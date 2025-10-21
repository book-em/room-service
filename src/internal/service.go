package internal

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/util"
	"context"
	"fmt"
	"math"
	"strings"
	"time"
)

type Service interface {
	Create(context context.Context, callerID uint, dto CreateRoomDTO) (*Room, error)
	FindById(context context.Context, id uint) (*Room, error)
	FindByHost(context context.Context, hostId uint) ([]Room, error)
	FindAvailableRooms(context context.Context, dto RoomsQueryDTO) ([]RoomResultDTO, *PaginatedResultInfoDTO, error)

	FindAvailabilityListById(context context.Context, id uint) (*RoomAvailabilityList, error)
	FindAvailabilityListsByRoomId(context context.Context, roomId uint) ([]RoomAvailabilityList, error)
	FindCurrentAvailabilityListOfRoom(context context.Context, roomId uint) (*RoomAvailabilityList, error)
	UpdateAvailability(context context.Context, callerID uint, dto CreateRoomAvailabilityListDTO) (*RoomAvailabilityList, error)

	FindPriceListById(context context.Context, id uint) (*RoomPriceList, error)
	FindPriceListsByRoomId(context context.Context, roomId uint) ([]RoomPriceList, error)
	FindCurrentPriceListOfRoom(context context.Context, roomId uint) (*RoomPriceList, error)
	UpdatePriceList(context context.Context, callerID uint, dto CreateRoomPriceListDTO) (*RoomPriceList, error)

	ClearYear(context context.Context, dateFrom time.Time, dateTo time.Time) (time.Time, time.Time)
	// CalculatePriceForOneDay computes the price for the room for a single night.
	// If the room is priced by guest, then the resulting price is multiplied by the number of guests.
	//
	// In other words, this is the total price for a single night. If you want the price for a single
	// guest, you need to determine if the room is priced per guest and if so, divide by the number of
	// guests.
	//
	// TODO: This should NOT return float32.
	CalculatePriceForOneDay(context context.Context, day time.Time, guests uint, rules RoomPriceList) float32
	// CalculatePrice calculates the price of the room between dateFrom and dateTo.
	//
	// It's assumed that the room can be booked in this date range.
	// Returns the total price, whether the price is flat or per guest and any error.
	// If the room is priced per guest, the returned price is the total price for all guests.
	// So if you want the price for a single guest, divide by the number of guests.
	//
	// TODO: This should NOT return float32.
	CalculatePrice(context context.Context, dateFrom time.Time, dateTo time.Time, guestsNumber uint, roomId uint) (float32, bool, error)
	IsRoomAvailableForOneDay(context context.Context, day time.Time, rules []RoomAvailabilityItem) bool
	IsRoomAvailable(context context.Context, dateFrom time.Time, dateTo time.Time, roomId uint) bool
	CalculateUnitPrice(context context.Context, perGuest bool, guestsNumber uint, dateFrom time.Time, dateTo time.Time, totalPrice float32) float32
	PreparePaginatedResult(context context.Context, hits []RoomResultDTO, pageNumber uint, pageSize uint) ([]RoomResultDTO, PaginatedResultInfoDTO)

	QueryForReservation(context context.Context, callerID uint, dto RoomReservationQueryDTO) (*RoomReservationQueryResponseDTO, error)
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

func (s *service) Create(context context.Context, callerID uint, dto CreateRoomDTO) (*Room, error) {
	util.TEL.Info("user wants to create a room", "caller_id", callerID)

	util.TEL.Push(context, "validate-user")
	defer util.TEL.Pop()

	// Check if user exists.

	util.TEL.Debug("check if user exists", "id", callerID)
	caller, err := s.userClient.FindById(util.TEL.Ctx(), callerID)
	if err != nil {
		util.TEL.Error("user does not exist", err, "id", callerID)
		return nil, err
	}

	// Check if user is host.

	util.TEL.Debug("check if user is a host", "id", callerID)
	if caller.Role != string(util.Host) {
		util.TEL.Error("user has a bad role", nil, "role", caller.Role)
		return nil, ErrUnauthorized
	}

	// User must be creating a room for himself.

	util.TEL.Debug("user must be creating a room for himself")
	if caller.Id != dto.HostID {
		util.TEL.Error("wrong user", nil, "caller_id", caller.Id, "host_id", dto.HostID)
		return nil, ErrUnauthorized
	}

	// First create the room without photos.

	util.TEL.Push(context, "create-room-in-db")
	defer util.TEL.Pop()

	room := &Room{
		HostID:      dto.HostID,
		Name:        dto.Name,
		Description: dto.Description,
		Address:     dto.Address,
		MinGuests:   dto.MinGuests,
		MaxGuests:   dto.MaxGuests,
		Photos:      []string{},
		Commodities: dto.Commodities,
		AutoApprove: dto.AutoApprove,
	}

	err = s.repo.Create(room)
	if err != nil {
		return nil, err
	}

	// Then add the photos (because we want deterministic filenames, so we need the ID).

	util.TEL.Push(context, "add-all-photos-on-disk")
	defer util.TEL.Pop()

	var photos = make([]string, 0)
	for _, imageBase64 := range dto.PhotosPayload {
		imgFname := fmt.Sprintf("room-%d-%d", room.ID, len(photos))
		_, path, err := util.SaveImageB64(imageBase64, imgFname)
		if err != nil {
			util.TEL.Error("could not save image", err, "fname", imgFname)
			s.repo.Delete(room)
			return nil, err
		}
		photos = append(photos, path)
	}

	// Then update the model with the photos.

	util.TEL.Push(context, "add-refs-to-photos-in-db")
	defer util.TEL.Pop()

	room.Photos = photos
	err = s.repo.Update(room)
	if err != nil {
		util.TEL.Error("could not create room (updating with photos failed)", err)
		util.TEL.Warn("deleting room...")
		s.repo.Delete(room)
		return nil, err
	}

	return room, nil
}

func (s *service) FindById(context context.Context, id uint) (*Room, error) {
	util.TEL.Info("find room", "id", id)

	util.TEL.Push(context, "find-room-in-db")
	defer util.TEL.Pop()

	room, err := s.repo.FindById(id)
	if err != nil {
		util.TEL.Error("room not found", err, "id", id)
		return nil, ErrNotFound("room", id)
	}
	return room, nil
}

func (s *service) FindByHost(context context.Context, hostId uint) ([]Room, error) {
	util.TEL.Info("find rooms owned by host", "host_id", hostId)

	util.TEL.Push(context, "validate-user")
	defer util.TEL.Pop()

	// Check if user exists.

	util.TEL.Debug("check if user exists", "id", hostId)
	host, err := s.userClient.FindById(util.TEL.Ctx(), hostId)
	if err != nil {
		util.TEL.Error("user does not exist", err, "id", hostId)
		return nil, ErrNotFound("host", hostId)
	}

	// Check if user is host.

	util.TEL.Debug("check if user is a host", "id", hostId)
	if host.Role != string(util.Host) {
		util.TEL.Error("user has a bad role", nil, "role", host.Role)
		return nil, ErrUnauthorized
	}

	// Fetch rooms.

	util.TEL.Push(context, "find-rooms-in-db")
	defer util.TEL.Pop()

	rooms, err := s.repo.FindByHost(hostId)
	if err != nil {
		util.TEL.Error("could not find rooms by host", err)
		return nil, ErrNotFound("rooms of host", hostId)
	}
	return rooms, nil
}

func (s *service) FindAvailabilityListById(context context.Context, id uint) (*RoomAvailabilityList, error) {
	util.TEL.Info("find room availability list", "list_id", id)

	util.TEL.Push(context, "find-availability-list-in-db")
	defer util.TEL.Pop()

	li, err := s.availabiltyRepo.FindListById(id)
	if err != nil {
		util.TEL.Error("availability list not found", err, "list_id", id)
		return nil, ErrNotFound("room availability list", id)
	}
	return li, err
}

func (s *service) FindAvailabilityListsByRoomId(context context.Context, roomId uint) ([]RoomAvailabilityList, error) {
	util.TEL.Info("find availability lists by room", "id", roomId)

	// TODO: Should I push and pop here?

	_, err := s.FindById(util.TEL.Ctx(), roomId)
	if err != nil {
		util.TEL.Error("room not found", err, "id", roomId)
		return nil, ErrNotFound("room", roomId)
	}

	util.TEL.Push(context, "find-availability-lists-in-db")
	defer util.TEL.Pop()

	lists, err := s.availabiltyRepo.FindListsByRoomId(roomId)
	if err != nil {
		util.TEL.Error("availability lists not found", err)
		return nil, ErrNotFound("room availability lists", roomId)
	}
	return lists, err
}

func (s *service) FindCurrentAvailabilityListOfRoom(context context.Context, roomId uint) (*RoomAvailabilityList, error) {
	util.TEL.Info("find current availability list of room", "room_id", roomId)

	util.TEL.Push(context, "find-current-availability-list-in-db")
	defer util.TEL.Pop()

	li, err := s.availabiltyRepo.FindCurrentListOfRoom(roomId)
	if err != nil {
		util.TEL.Error("availability list not found", err)
		return nil, ErrNotFound("room availability list", roomId)
	}
	return li, err
}

func (s *service) UpdateAvailability(context context.Context, callerID uint, dto CreateRoomAvailabilityListDTO) (*RoomAvailabilityList, error) {
	// Idea:
	//
	// Each list is read-only, when you change it, you're actually creating a new one.
	// Our API allows modifying a list by giving it the entire array of items.
	// So this method does both updating and deleting.

	util.TEL.Info("update availability of room", "room_id", dto.RoomID)

	util.TEL.Push(context, "validate-room-and-user")
	defer util.TEL.Pop()

	util.TEL.Debug("check if user exists", "id", callerID)
	caller, err := s.userClient.FindById(util.TEL.Ctx(), callerID)
	if err != nil {
		util.TEL.Error("user does not exist", err, "id", callerID)
		return nil, err
	}

	util.TEL.Debug("check if user is a host", "id", callerID)
	if caller.Role != string(util.Host) {
		util.TEL.Error("user has a bad role", nil, "role", caller.Role)
		return nil, ErrUnauthorized
	}

	util.TEL.Debug("find room", "id", dto.RoomID)
	// TODO: Should I push and pop here?
	room, err := s.FindById(util.TEL.Ctx(), dto.RoomID)
	if err != nil {
		util.TEL.Error("room not found", err, "id", dto.RoomID)
		return nil, err
	}

	util.TEL.Debug("caller must own the room")
	if room.HostID != callerID {
		util.TEL.Error("user is not owner of this room", nil, "user_id", callerID, "owner_id", room.HostID)
		return nil, ErrUnauthorized
	}

	util.TEL.Push(context, "validate-availability-list")
	defer util.TEL.Pop()

	util.TEL.Debug("create availability list")
	newList := RoomAvailabilityList{
		RoomID:        dto.RoomID,
		EffectiveFrom: time.Now(),
		Items:         make([]RoomAvailabilityItem, 0, len(dto.Items)),
	}

	util.TEL.Debug("validate and create items for the availability list")
	for i, item := range dto.Items {
		from := util.ClearYear(item.DateFrom)
		to := util.ClearYear(item.DateTo)

		if from.After(to) {
			util.TEL.Error("invalid date range", nil, "from", from, "to", to)
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
				util.TEL.Error("duplicate availability rule", nil, "index1", i, "index2", j)
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

	util.TEL.Push(context, "save-availability-list-to-db")
	defer util.TEL.Pop()

	err = s.availabiltyRepo.CreateList(&newList)
	if err != nil {
		util.TEL.Error("could not create availability list in db", err)
		return nil, err
	}

	return &newList, nil
}

func (s *service) FindPriceListById(context context.Context, id uint) (*RoomPriceList, error) {
	util.TEL.Info("find room price list", "list_id", id)

	util.TEL.Push(context, "find-availability-list-in-db")
	defer util.TEL.Pop()

	list, err := s.priceRepo.FindListById(id)
	if err != nil {
		util.TEL.Error("room price list not found", err, "list_id", id)
		return nil, ErrNotFound("room price list", id)
	}
	return list, nil
}

func (s *service) FindPriceListsByRoomId(context context.Context, roomId uint) ([]RoomPriceList, error) {
	util.TEL.Info("find room price lists by room", "room_id", roomId)

	util.TEL.Push(context, "find-availability-lists-in-db")
	defer util.TEL.Pop()

	// TODO: Should I push and pop here?

	_, err := s.FindById(util.TEL.Ctx(), roomId)
	if err != nil {
		util.TEL.Error("room not found", err, "id", roomId)
		return nil, ErrNotFound("room", roomId)
	}

	lists, err := s.priceRepo.FindListsByRoomId(roomId)
	if err != nil {
		util.TEL.Error("room price lists of room not found", err, "room_id", roomId)
		return nil, ErrNotFound("room price lists", roomId)
	}
	return lists, nil
}

func (s *service) FindCurrentPriceListOfRoom(context context.Context, roomId uint) (*RoomPriceList, error) {
	util.TEL.Info("find current price list of room", nil, "room_id", roomId)

	util.TEL.Push(context, "find-current-price-list-in-db")
	defer util.TEL.Pop()

	list, err := s.priceRepo.FindCurrentListOfRoom(roomId)
	if err != nil {
		util.TEL.Error("room current price lists of room not found", err, "room_id", roomId)
		return nil, ErrNotFound("room price list", roomId)
	}
	return list, nil
}

func (s *service) UpdatePriceList(context context.Context, callerID uint, dto CreateRoomPriceListDTO) (*RoomPriceList, error) {
	util.TEL.Info("update price list of room %d", nil, "room_id", dto.RoomID)

	util.TEL.Push(context, "validate-room-and-user")
	defer util.TEL.Pop()

	util.TEL.Debug("check if user exists", "id", callerID)
	caller, err := s.userClient.FindById(util.TEL.Ctx(), callerID)
	if err != nil {
		util.TEL.Error("user does not exist", err, "id", callerID)
		return nil, err
	}

	util.TEL.Debug("check if user is a host", "id", callerID)
	if caller.Role != string(util.Host) {
		util.TEL.Error("user has a bad role", nil, "role", caller.Role)
		return nil, ErrUnauthorized
	}

	util.TEL.Debug("find room", "id", dto.RoomID)
	// TODO: Should I push and pop here?
	room, err := s.FindById(util.TEL.Ctx(), dto.RoomID)
	if err != nil {
		util.TEL.Error("room not found", err, "id", dto.RoomID)
		return nil, err
	}

	util.TEL.Debug("caller must own the room")
	if room.HostID != callerID {
		util.TEL.Error("user is not owner of this room", nil, "user_id", callerID, "owner_id", room.HostID)
		return nil, ErrUnauthorized
	}

	util.TEL.Push(context, "validate-availability-list")
	defer util.TEL.Pop()

	util.TEL.Debug("create price list")
	newList := RoomPriceList{
		RoomID:        dto.RoomID,
		EffectiveFrom: time.Now(),
		BasePrice:     dto.BasePrice,
		PerGuest:      dto.PerGuest,
		Items:         make([]RoomPriceItem, 0, len(dto.Items)),
	}

	util.TEL.Debug("validate and create items for the price list")
	for i, item := range dto.Items {
		from := util.ClearYear(item.DateFrom)
		to := util.ClearYear(item.DateTo)

		if from.After(to) {
			util.TEL.Error("invalid date range", nil, "from", from, "to", to)
			return nil, ErrBadRequestCustom(fmt.Sprintf("invalid date range: %v > %v", from, to))
		}

		for j, item2 := range dto.Items {
			if i == j {
				continue
			}

			from2 := util.ClearYear(item2.DateFrom)
			to2 := util.ClearYear(item2.DateTo)

			if !from.After(to2) && !from2.After(to) {
				util.TEL.Error("price rules conflict (no intersections allowed)", nil, "item1", i, "item2", j)
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

	util.TEL.Push(context, "save-price-list-to-db")
	defer util.TEL.Pop()

	err = s.priceRepo.CreateList(&newList)
	if err != nil {
		util.TEL.Error("could not create price list in db", err)
		return nil, err
	}

	return &newList, nil
}

func (s *service) ClearYear(context context.Context, dateFrom time.Time, dateTo time.Time) (time.Time, time.Time) {
	util.TEL.Debug("Clearing year from date range", "from", dateFrom, "to", dateTo)
	dateFrom = util.ClearYear(dateFrom)
	dateTo = util.ClearYear(dateTo)
	util.TEL.Debug("Resulting date range", "from", dateFrom, "to", dateTo)
	return dateFrom, dateTo
}

func (s *service) CalculatePriceForOneDay(context context.Context, day time.Time, guests uint, rules RoomPriceList) float32 {
	util.TEL.Info("calculating price for one day", "day", day, "guests", guests, "room_id", rules.RoomID, "pricelist_id", rules.ID)

	normalizedDay := util.ClearYear(day)
	price := rules.BasePrice

	for _, rule := range rules.Items {
		rule.DateFrom, rule.DateTo = s.ClearYear(util.TEL.Ctx(), rule.DateFrom, rule.DateTo)

		if !normalizedDay.Before(rule.DateFrom) && !normalizedDay.After(rule.DateTo) {
			price = rule.Price
		}
	}
	util.TEL.Debug("unit price for this day is", "price", price)

	if rules.PerGuest {
		util.TEL.Debug("price is per guest")
		return float32(price * guests)
	}

	util.TEL.Debug("price is flat rate")
	return float32(price)
}

func (s *service) CalculatePrice(context context.Context, dateFrom time.Time, dateTo time.Time, guests uint, roomId uint) (float32, bool, error) {
	util.TEL.Info("calculating price for a date range", "from", dateFrom, "to", dateTo, "guests", guests, "room_id", roomId)

	rules, err := s.FindCurrentPriceListOfRoom(util.TEL.Ctx(), roomId)
	if err != nil {
		return float32(0), false, err
	}

	dateFrom, dateTo = s.ClearYear(util.TEL.Ctx(), dateFrom, dateTo)
	var totalPrice float32

	for day := dateFrom; !day.After(dateTo); day = day.Add(24 * time.Hour) {
		totalPrice += s.CalculatePriceForOneDay(util.TEL.Ctx(), day, guests, *rules)
	}
	util.TEL.Debug("result", "total_price", totalPrice, "price_is_per_guest", rules.PerGuest)

	return totalPrice, rules.PerGuest, nil
}

func (s *service) IsRoomAvailableForOneDay(context context.Context, day time.Time, rules []RoomAvailabilityItem) bool {
	util.TEL.Info("is the room available on a specific day", "day", day)

	leastRule := RoomAvailabilityItem{
		DateFrom:  time.Time{}, // The zero value represents the earliest possible time
		DateTo:    time.Now(),
		Available: false,
	}

	dayNormalized := util.ClearYear(day)

	for _, rule := range rules {
		rule.DateFrom, rule.DateTo = s.ClearYear(util.TEL.Ctx(), rule.DateFrom, rule.DateTo)

		if !dayNormalized.Before(rule.DateFrom) && !dayNormalized.After(rule.DateTo) {
			if rule.DateTo.Sub(rule.DateFrom) < leastRule.DateTo.Sub(leastRule.DateFrom) {
				leastRule = rule
			}
		}
	}

	return leastRule.Available
}

func (s *service) IsRoomAvailable(context context.Context, dateFrom time.Time, dateTo time.Time, roomId uint) bool {
	util.TEL.Info("is the room available between multiple days", "from", dateFrom, "to", dateTo, "room_id", roomId)

	dateFrom, dateTo = s.ClearYear(util.TEL.Ctx(), dateFrom, dateTo)

	rules, err := s.FindCurrentAvailabilityListOfRoom(util.TEL.Ctx(), roomId)
	if err != nil {
		util.TEL.Debug("no availability list => room is unavailable")
		return false
	}

	for day := dateFrom; !day.After(dateTo); day = day.Add(24 * time.Hour) {
		if s.IsRoomAvailableForOneDay(util.TEL.Ctx(), day, rules.Items) == false {
			util.TEL.Debug("room is unavailable on this day", "day", day)
			return false
		}
	}
	return true
}

func (s *service) CalculateUnitPrice(context context.Context, perGuest bool, guestsNumber uint, dateFrom time.Time, dateTo time.Time, totalPrice float32) float32 {
	util.TEL.Info("calculating unit price", "guests", guestsNumber, "per_guest", perGuest, "from", dateFrom, "to", dateTo, "total_price", totalPrice)

	var unitPrice float32
	interval := float32(dateTo.Sub(dateFrom).Hours()/24) + 1

	if perGuest {
		unitPrice = totalPrice / interval / float32(guestsNumber)
	} else {
		unitPrice = totalPrice / interval
	}

	util.TEL.Info("unit price is", "price", unitPrice)

	return unitPrice
}

func (s *service) PreparePaginatedResult(context context.Context, hits []RoomResultDTO, pageNumber uint, pageSize uint) ([]RoomResultDTO, PaginatedResultInfoDTO) {
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

func (s *service) FindAvailableRooms(context context.Context, dto RoomsQueryDTO) ([]RoomResultDTO, *PaginatedResultInfoDTO, error) {
	util.TEL.Info("find available rooms from query", "query", fmt.Sprintf("%+v", dto))

	from := util.ClearYear(dto.DateFrom)
	to := util.ClearYear(dto.DateTo)

	util.TEL.Push(context, "find by filters")
	defer util.TEL.Pop()

	if from.After(to) {
		util.TEL.Error("invalid date range", nil, "from", from, "to", to)
		return nil, nil, ErrBadRequestCustom(fmt.Sprintf("invalid date range: %v > %v", from, to))
	}

	rooms, err := s.repo.FindByFilters(dto.GuestsNumber, strings.TrimSpace(dto.Address))
	if err != nil {
		util.TEL.Error("could not perform query", err)
		return nil, nil, err
	}

	util.TEL.Push(context, "get price for each hit")
	defer util.TEL.Pop()

	var hits []RoomResultDTO
	for _, room := range rooms {
		canBook := s.IsRoomAvailable(util.TEL.Ctx(), from, to, room.ID)

		if canBook {
			totalPrice, perGuest, err := s.CalculatePrice(util.TEL.Ctx(), from, to, dto.GuestsNumber, room.ID)
			if err != nil {
				util.TEL.Error("could not calculate price", err)
				continue
			}
			unitPrice := s.CalculateUnitPrice(util.TEL.Ctx(), perGuest, dto.GuestsNumber, from, to, totalPrice)
			hits = append(hits, NewRoomResultDTO(room, perGuest, unitPrice, totalPrice))
		}
	}

	util.TEL.Push(context, "build result")
	defer util.TEL.Pop()

	hits, resultInfo := s.PreparePaginatedResult(util.TEL.Ctx(), hits, dto.PageNumber, dto.PageSize)

	return hits, &resultInfo, nil
}

func (s *service) QueryForReservation(context context.Context, callerID uint, dto RoomReservationQueryDTO) (*RoomReservationQueryResponseDTO, error) {
	util.TEL.Info("query room for reservation", "id", dto.RoomID)

	util.TEL.Push(context, "validate-room-and-user")
	defer util.TEL.Pop()

	util.TEL.Debug("check if user exists", "id", callerID)
	caller, err := s.userClient.FindById(util.TEL.Ctx(), callerID)
	if err != nil {
		util.TEL.Error("user does not exist", err, "id", callerID)
		return nil, err
	}

	util.TEL.Debug("check if user is a guest", "id", callerID)
	if caller.Role != string(util.Guest) {
		util.TEL.Error("user has a bad role", nil, "role", caller.Role)
		return nil, ErrUnauthorized
	}

	util.TEL.Debug("find room", "id", dto.RoomID)
	// TODO: Should I push and pop here? and elsewhere where i call subfunc
	room, err := s.FindById(util.TEL.Ctx(), dto.RoomID)
	if err != nil {
		util.TEL.Error("room not found", err, "id", dto.RoomID)
		return nil, err
	}

	util.TEL.Push(context, "query")
	defer util.TEL.Pop()

	util.TEL.Debug("find room availability", "id", room.ID)
	isAvailable := s.IsRoomAvailable(util.TEL.Ctx(), dto.DateFrom, dto.DateTo, room.ID)

	if !isAvailable {
		util.TEL.Error("room cannot be booked at this date range - returning early", nil)

		return &RoomReservationQueryResponseDTO{
			Available: isAvailable,
			TotalCost: 0,
		}, nil
	}

	util.TEL.Debug("calculate price for this potential reservation")
	fullPrice, _, err := s.CalculatePrice(util.TEL.Ctx(), dto.DateFrom, dto.DateTo, dto.GuestCount, room.ID)

	if err != nil {
		return nil, err
	}

	return &RoomReservationQueryResponseDTO{
		Available: isAvailable,
		TotalCost: uint(fullPrice), // TODO: Remove this cast once CalculatePrice returns uint.
	}, nil
}
