package internal

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/util"
	"fmt"
	"log"
)

type Service interface {
	Create(callerID uint, dto CreateRoomDTO) (*Room, error)
	FindById(id uint) (*Room, error)
	FindByHost(hostId uint) ([]Room, error)
}

type service struct {
	repo       Repository
	userClient userclient.UserClient
}

func NewService(r Repository, userClient userclient.UserClient) Service {
	return &service{r, userClient}
}

func (s *service) Create(callerID uint, dto CreateRoomDTO) (*Room, error) {
	// Check if user exists.

	caller, err := s.userClient.FindById(callerID)
	if err != nil {
		return nil, err
	}

	// Check if user is host.

	if caller.Role != string(userclient.Host) {
		return nil, fmt.Errorf("Unauthorized (bad role)")
	}

	// User must be creating a room for himself.

	if caller.Id != dto.HostID {
		return nil, fmt.Errorf("Unauthorized (wrong user %d but caller is %d)", dto.HostID, caller.Id)
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
		return nil, fmt.Errorf("room %d not found: %v", id, err)
	}
	return room, nil
}

func (s *service) FindByHost(hostId uint) ([]Room, error) {
	// Check if user exists.

	host, err := s.userClient.FindById(hostId)
	if err != nil {
		return nil, err
	}

	// Check if user is host.

	if host.Role != string(userclient.Host) {
		return nil, fmt.Errorf("user %d is not host", hostId)
	}

	// Fetch rooms.

	rooms, err := s.repo.FindByHost(hostId)

	if err != nil {
		return nil, fmt.Errorf("rooms of host %d not found: %v", hostId, err)
	}
	return rooms, nil
}
