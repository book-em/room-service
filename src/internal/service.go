package internal

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/util"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
)

type Service interface {
	Create(callerID uint, dto CreateRoomDTO) (*Room, error)
	FindById(id uint) (*Room, error)
	FindByHost(hostId uint) ([]Room, error)
	FindAvailableRooms(dto RoomsQueryDTO) ([]RoomResultDTO, *PaginatedResultInfoDTO, error)

	FindAvailabilityListById(id uint) (*RoomAvailabilityList, error)
	FindAvailabilityListsByRoomId(roomId uint) ([]RoomAvailabilityList, error)
	FindCurrentAvailabilityListOfRoom(roomId uint) (*RoomAvailabilityList, error)
	UpdateAvailability(callerID uint, dto CreateRoomAvailabilityListDTO) (*RoomAvailabilityList, error)

	FindPriceListById(id uint) (*RoomPriceList, error)
	FindPriceListsByRoomId(roomId uint) ([]RoomPriceList, error)
	FindCurrentPriceListOfRoom(roomId uint) (*RoomPriceList, error)
	UpdatePriceList(callerID uint, dto CreateRoomPriceListDTO) (*RoomPriceList, error)

	ClearYear(dateFrom time.Time, dateTo time.Time) (time.Time, time.Time)
	// CalculatePriceForOneDay computes the price for the room for a single night.
	// If the room is priced by guest, then the resulting price is multiplied by the number of guests.
	//
	// In other words, this is the total price for a single night. If you want the price for a single
	// guest, you need to determine if the room is priced per guest and if so, divide by the number of
	// guests.
	//
	// TODO: This should NOT return float32.
	CalculatePriceForOneDay(day time.Time, guests uint, rules RoomPriceList) float32
	// CalculatePrice calculates the price of the room between dateFrom and dateTo.
	//
	// It's assumed that the room can be booked in this date range.
	// Returns the total price, whether the price is flat or per guest and any error.
	// If the room is priced per guest, the returned price is the total price for all guests.
	// So if you want the price for a single guest, divide by the number of guests.
	//
	// TODO: This should NOT return float32.
	CalculatePrice(dateFrom time.Time, dateTo time.Time, guestsNumber uint, roomId uint) (float32, bool, error)
	IsRoomAvailableForOneDay(day time.Time, rules []RoomAvailabilityItem) bool
	IsRoomAvailable(dateFrom time.Time, dateTo time.Time, roomId uint) bool
	CalculateUnitPrice(perGuest bool, guestsNumber uint, dateFrom time.Time, dateTo time.Time, totalPrice float32) float32
	PreparePaginatedResult(hits []RoomResultDTO, pageNumber uint, pageSize uint) ([]RoomResultDTO, PaginatedResultInfoDTO)

	QueryForReservation(callerID uint, dto RoomReservationQueryDTO) (*RoomReservationQueryResponseDTO, error)
}

type service struct {
	repo            Repository
	availabiltyRepo RoomAvailabilityRepo
	priceRepo       RoomPriceRepo
	userClient      userclient.UserClient
}

func NewService(
	roomRepo Repository,
	availabiltyRepo RoomAvailabilityRepo,
	priceRepo RoomPriceRepo,
	userClient userclient.UserClient) Service {
	return &service{roomRepo, availabiltyRepo, priceRepo, userClient}
}

func (s *service) Create(callerID uint, dto CreateRoomDTO) (*Room, error) {
	// Check if user exists.

	caller, err := s.userClient.FindById(callerID)
	if err != nil {
		return nil, err
	}

	// Check if user is host.

	if caller.Role != string(userclient.Host) {
		log.Printf("Unauthorized (bad role %s)", caller.Role)
		return nil, ErrUnauthorized
	}

	// User must be creating a room for himself.

	if caller.Id != dto.HostID {
		log.Printf("Unauthorized (wrong user %d but caller is %d)", dto.HostID, caller.Id)
		return nil, ErrUnauthorized
	}

	// First create the room without photos.

	room := &Room{
		HostID:      dto.HostID,
		Name:        dto.Name,
		Description: dto.Description,
		Address:     dto.Address,
		MinGuests:   dto.MinGuests,
		MaxGuests:   dto.MaxGuests,
		Photos:      []string{},
		Commodities: dto.Commodities,
	}

	err = s.repo.Create(room)
	if err != nil {
		return nil, err
	}

	// Then add the photos (because we want deterministic filenames, so we need the ID).

	var photos = make([]string, 0)
	for _, imageBase64 := range dto.PhotosPayload {
		imgFname := fmt.Sprintf("room-%d-%d", room.ID, len(photos))
		_, path, err := util.SaveImageB64(imageBase64, imgFname)
		if err != nil {
			log.Printf("Could not save image %s: %v", imgFname, err)
			s.repo.Delete(room)
			return nil, err
		}
		photos = append(photos, path)
	}

	// Then update the model with the photos.

	room.Photos = photos
	err = s.repo.Update(room)
	if err != nil {
		log.Printf("Could not update room with images: %v", err)
		s.repo.Delete(room)
		return nil, err
	}

	return room, nil
}

func (s *service) FindById(id uint) (*Room, error) {
	room, err := s.repo.FindById(id)
	if err != nil {
		return nil, ErrNotFound("room", id)
	}
	return room, nil
}

func (s *service) FindByHost(hostId uint) ([]Room, error) {
	log.Printf("Find rooms by host %d", hostId)

	// Check if user exists.

	host, err := s.userClient.FindById(hostId)
	if err != nil {
		return nil, ErrNotFound("host", hostId)
	}

	// Check if user is host.

	if host.Role != string(userclient.Host) {
		log.Printf("Unauthorized (user %d is not host)", hostId)
		return nil, ErrUnauthorized
	}

	// Fetch rooms.

	rooms, err := s.repo.FindByHost(hostId)

	if err != nil {
		log.Printf("%s", err.Error())
		return nil, ErrNotFound("rooms of host", hostId)
	}
	return rooms, nil
}

func (s *service) FindAvailabilityListById(id uint) (*RoomAvailabilityList, error) {
	li, err := s.availabiltyRepo.FindListById(id)
	if err != nil {
		return nil, ErrNotFound("room availability list", id)
	}
	return li, err
}

func (s *service) FindAvailabilityListsByRoomId(roomId uint) ([]RoomAvailabilityList, error) {
	_, err := s.FindById(roomId)
	if err != nil {
		return nil, ErrNotFound("room", roomId)
	}

	lists, err := s.availabiltyRepo.FindListsByRoomId(roomId)
	if err != nil {
		return nil, ErrNotFound("room availability lists", roomId)
	}
	return lists, err
}

func (s *service) FindCurrentAvailabilityListOfRoom(roomId uint) (*RoomAvailabilityList, error) {
	li, err := s.availabiltyRepo.FindCurrentListOfRoom(roomId)
	if err != nil {
		return nil, ErrNotFound("room availability list", roomId)
	}
	return li, err
}

func (s *service) UpdateAvailability(callerID uint, dto CreateRoomAvailabilityListDTO) (*RoomAvailabilityList, error) {
	// Idea:
	//
	// Each list is read-only, when you change it, you're actually creating a new one.
	// Our API allows modifying a list by giving it the entire array of items.
	// So this method does both updating and deleting.

	log.Printf("[1] User exists")

	caller, err := s.userClient.FindById(callerID)
	if err != nil {
		return nil, err
	}

	log.Printf("[2] User is host")

	if caller.Role != string(userclient.Host) {
		log.Printf("Unauthorized (bad role %s)", caller.Role)
		return nil, ErrUnauthorized
	}

	log.Printf("[3] Room exists")

	room, err := s.FindById(dto.RoomID)
	if err != nil {
		return nil, err
	}

	log.Printf("[4] Host owns the room")

	if room.HostID != callerID {
		return nil, ErrUnauthorized
	}

	log.Printf("[5] Create availability list")

	newList := RoomAvailabilityList{
		RoomID:        dto.RoomID,
		EffectiveFrom: time.Now(),
		Items:         make([]RoomAvailabilityItem, 0, len(dto.Items)),
	}

	log.Printf("[6] Validate and create availability list items")

	for i, item := range dto.Items {
		from := util.ClearYear(item.DateFrom)
		to := util.ClearYear(item.DateTo)

		if from.After(to) {
			log.Printf("invalid date range: %v > %v", from, to)

			return nil, ErrBadRequestCustom(fmt.Sprintf("invalid date range: %v > %v", from, to))
		}

		// This loop could be optimized.
		for j, item2 := range dto.Items {
			if i == j {
				continue
			}

			from2 := util.ClearYear(item2.DateFrom)
			to2 := util.ClearYear(item2.DateTo)

			if from == from2 && to == to2 {
				return nil, ErrBadRequestCustom(fmt.Sprintf("duplicate availability rule at index %d and %d", i, j))
			}
		}

		newList.Items = append(newList.Items, RoomAvailabilityItem{
			ID:        item.ExistingID,
			DateFrom:  item.DateFrom,
			DateTo:    item.DateTo,
			Available: item.Available,
		})
	}

	log.Printf("[7] Save availability list to DB")

	err = s.availabiltyRepo.CreateList(&newList)
	if err != nil {
		return nil, err
	}

	return &newList, nil
}

func (s *service) FindPriceListById(id uint) (*RoomPriceList, error) {
	list, err := s.priceRepo.FindListById(id)
	if err != nil {
		return nil, ErrNotFound("room price list", id)
	}
	return list, nil
}

func (s *service) FindPriceListsByRoomId(roomId uint) ([]RoomPriceList, error) {
	_, err := s.FindById(roomId)
	if err != nil {
		return nil, ErrNotFound("room", roomId)
	}

	lists, err := s.priceRepo.FindListsByRoomId(roomId)
	if err != nil {
		return nil, ErrNotFound("room price lists", roomId)
	}
	return lists, nil
}

func (s *service) FindCurrentPriceListOfRoom(roomId uint) (*RoomPriceList, error) {
	list, err := s.priceRepo.FindCurrentListOfRoom(roomId)
	if err != nil {
		return nil, ErrNotFound("room price list", roomId)
	}
	return list, nil
}

func (s *service) UpdatePriceList(callerID uint, dto CreateRoomPriceListDTO) (*RoomPriceList, error) {
	log.Printf("[1] User exists")

	caller, err := s.userClient.FindById(callerID)
	if err != nil {
		return nil, err
	}

	log.Printf("[2] User is host")

	if caller.Role != string(userclient.Host) {
		log.Printf("Unauthorized (bad role %s)", caller.Role)
		return nil, ErrUnauthorized
	}

	log.Printf("[3] Room exists")

	room, err := s.FindById(dto.RoomID)
	if err != nil {
		return nil, err
	}

	log.Printf("[4] Host owns the room")

	if room.HostID != callerID {
		return nil, ErrUnauthorized
	}

	log.Printf("[5] Create price list")

	newList := RoomPriceList{
		RoomID:        dto.RoomID,
		EffectiveFrom: time.Now(),
		BasePrice:     dto.BasePrice,
		PerGuest:      dto.PerGuest,
		Items:         make([]RoomPriceItem, 0, len(dto.Items)),
	}

	log.Printf("[6] Validate and create price list items")

	for i, item := range dto.Items {
		from := util.ClearYear(item.DateFrom)
		to := util.ClearYear(item.DateTo)

		if from.After(to) {
			return nil, ErrBadRequestCustom(fmt.Sprintf("invalid date range: %v > %v", from, to))
		}

		for j, item2 := range dto.Items {
			if i == j {
				continue
			}

			from2 := util.ClearYear(item2.DateFrom)
			to2 := util.ClearYear(item2.DateTo)

			if !from.After(to2) && !from2.After(to) {
				return nil, ErrBadRequestCustom(fmt.Sprintf("price rules at index %d and %d conflict (no intersections allowed)", i, j))
			}
		}

		newList.Items = append(newList.Items, RoomPriceItem{
			ID:       item.ExistingID,
			DateFrom: item.DateFrom,
			DateTo:   item.DateTo,
			Price:    item.Price,
		})
	}

	log.Printf("[7] Save price list to DB")

	err = s.priceRepo.CreateList(&newList)
	if err != nil {
		return nil, err
	}

	return &newList, nil
}

func (s *service) ClearYear(dateFrom time.Time, dateTo time.Time) (time.Time, time.Time) {
	dateFrom = util.ClearYear(dateFrom)
	dateTo = util.ClearYear(dateTo)
	return dateFrom, dateTo
}

func (s *service) CalculatePriceForOneDay(day time.Time, guests uint, rules RoomPriceList) float32 {
	price := rules.BasePrice

	for _, rule := range rules.Items {
		rule.DateFrom, rule.DateTo = s.ClearYear(rule.DateFrom, rule.DateTo)

		if !day.Before(rule.DateFrom) && !day.After(rule.DateTo) {
			price = rule.Price
		}
	}

	if rules.PerGuest {
		return float32(price * guests)
	}

	return float32(price)
}

func (s *service) CalculatePrice(dateFrom time.Time, dateTo time.Time, guests uint, roomId uint) (float32, bool, error) {
	rules, err := s.FindCurrentPriceListOfRoom(roomId)
	if err != nil {
		return float32(0), false, err
	}

	dateFrom, dateTo = s.ClearYear(dateFrom, dateTo)
	var totalPrice float32

	for day := dateFrom; !day.After(dateTo); day = day.Add(24 * time.Hour) {
		totalPrice += s.CalculatePriceForOneDay(day, guests, *rules)
	}

	return totalPrice, rules.PerGuest, nil
}

func (s *service) IsRoomAvailableForOneDay(day time.Time, rules []RoomAvailabilityItem) bool {
	leastRule := RoomAvailabilityItem{
		DateFrom:  time.Time{}, // The zero value represents the earliest possible time
		DateTo:    time.Now(),
		Available: false,
	}

	for _, rule := range rules {
		rule.DateFrom, rule.DateTo = s.ClearYear(rule.DateFrom, rule.DateTo)

		if !day.Before(rule.DateFrom) && !day.After(rule.DateTo) {
			if rule.DateTo.Sub(rule.DateFrom) < leastRule.DateTo.Sub(leastRule.DateFrom) {
				leastRule = rule
			}
		}
	}

	return leastRule.Available
}

func (s *service) IsRoomAvailable(dateFrom time.Time, dateTo time.Time, roomId uint) bool {
	dateFrom, dateTo = s.ClearYear(dateFrom, dateTo)

	rules, err := s.FindCurrentAvailabilityListOfRoom(roomId)
	if err != nil {
		return false
	}

	for day := dateFrom; !day.After(dateTo); day = day.Add(24 * time.Hour) {
		if s.IsRoomAvailableForOneDay(day, rules.Items) == false {
			return false
		}
	}
	return true
}

func (s *service) CalculateUnitPrice(perGuest bool, guestsNumber uint, dateFrom time.Time, dateTo time.Time, totalPrice float32) float32 {
	var unitPrice float32
	interval := float32(dateTo.Sub(dateFrom).Hours()/24) + 1

	if perGuest {
		unitPrice = totalPrice / interval / float32(guestsNumber)
	} else {
		unitPrice = totalPrice / interval
	}

	return unitPrice
}

func (s *service) PreparePaginatedResult(hits []RoomResultDTO, pageNumber uint, pageSize uint) ([]RoomResultDTO, PaginatedResultInfoDTO) {

	totalHits := len(hits)
	totalPages := uint(math.Ceil(float64(totalHits) / float64(pageSize)))
	startIdx := uint((pageNumber - 1) * pageSize)
	endIdx := startIdx + pageSize
	lastPage := totalPages - 1
	lastPageStartIdx := lastPage * pageSize
	lastPageEndIdx := uint(totalHits)

	resultInfo := PaginatedResultInfoDTO{
		Page:       pageNumber,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalHits:  uint(totalHits),
	}

	if totalHits != 0 {
		// Show last page result if page number exceeds total
		if startIdx > lastPageEndIdx {
			startIdx = lastPageStartIdx
			endIdx = lastPageEndIdx
		}

		// Case when the last page is selected
		if endIdx > lastPageEndIdx {
			endIdx = lastPageEndIdx
		}

		hits = hits[startIdx:endIdx]
	}

	return hits, resultInfo
}

func (s *service) FindAvailableRooms(dto RoomsQueryDTO) ([]RoomResultDTO, *PaginatedResultInfoDTO, error) {

	from := util.ClearYear(dto.DateFrom)
	to := util.ClearYear(dto.DateTo)

	if from.After(to) {
		return nil, nil, ErrBadRequestCustom(fmt.Sprintf("invalid date range: %v > %v", from, to))
	}

	rooms, err := s.repo.FindByFilters(dto.GuestsNumber, strings.TrimSpace(dto.Address))
	if err != nil {
		return nil, nil, err
	}

	var hits []RoomResultDTO
	for _, room := range rooms {
		canBook := s.IsRoomAvailable(from, to, room.ID)

		if canBook {
			totalPrice, perGuest, err := s.CalculatePrice(from, to, dto.GuestsNumber, room.ID)
			if err != nil {
				continue
			}
			unitPrice := s.CalculateUnitPrice(perGuest, dto.GuestsNumber, from, to, totalPrice)
			hits = append(hits, NewRoomResultDTO(room, perGuest, unitPrice, totalPrice))
		}

	}

	hits, resultInfo := s.PreparePaginatedResult(hits, dto.PageNumber, dto.PageSize)

	return hits, &resultInfo, nil
}

func (s *service) QueryForReservation(callerID uint, dto RoomReservationQueryDTO) (*RoomReservationQueryResponseDTO, error) {
	log.Printf("[1] User exists")
	log.Printf("%d", callerID)

	caller, err := s.userClient.FindById(callerID)
	if err != nil {
		return nil, err
	}

	log.Printf("[2] User is guest")

	if caller.Role != string(userclient.Guest) {
		log.Printf("Unauthorized (bad role %s)", caller.Role)
		return nil, ErrUnauthorized
	}

	log.Printf("[3] Room exists")

	room, err := s.FindById(dto.RoomID)
	if err != nil {
		return nil, err
	}

	log.Printf("[4] Find room availability")

	isAvailable := s.IsRoomAvailable(dto.DateFrom, dto.DateTo, room.ID)

	if !isAvailable {
		log.Printf("[4.1] Room cannot be booked at this date range - returning early")

		return &RoomReservationQueryResponseDTO{
			Available: isAvailable,
			TotalCost: 0,
		}, nil
	}

	log.Printf("[5] Find price for this reservation")

	fullPrice, _, err := s.CalculatePrice(dto.DateFrom, dto.DateTo, dto.GuestCount, room.ID)

	if err != nil {
		return nil, err
	}

	return &RoomReservationQueryResponseDTO{
		Available: isAvailable,
		TotalCost: uint(fullPrice), // TODO: Remove this cast once CalculatePrice returns uint.
	}, nil
}
