package internal

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/util"
	"fmt"
	"log"
	"time"
)

type Service interface {
	Create(callerID uint, dto CreateRoomDTO) (*Room, error)
	FindById(id uint) (*Room, error)
	FindByHost(hostId uint) ([]Room, error)

	FindAvailabilityListById(id uint) (*RoomAvailabilityList, error)
	FindAvailabilityListsByRoomId(roomId uint) ([]RoomAvailabilityList, error)
	FindCurrentAvailabilityListOfRoom(roomId uint) (*RoomAvailabilityList, error)
	UpdateAvailability(callerID uint, dto CreateRoomAvailabilityListDTO) (*RoomAvailabilityList, error)

	FindPriceListById(id uint) (*RoomPriceList, error)
	FindPriceListsByRoomId(roomId uint) ([]RoomPriceList, error)
	FindCurrentPriceListOfRoom(roomId uint) (*RoomPriceList, error)
	UpdatePriceList(callerID uint, dto CreateRoomPriceListDTO) (*RoomPriceList, error)
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
