package internal

import "fmt"

type Service interface {
	Create(dto CreateRoomDTO) (*Room, error)
	FindById(id uint) (*Room, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) Create(dto CreateRoomDTO) (*Room, error) {
	var photos []string // TODO: Actually save photos and fetch the filenames

	room := &Room{
		HostID:      dto.HostID,
		Name:        dto.Name,
		Description: dto.Description,
		Address:     dto.Address,
		MinGuests:   dto.MinGuests,
		MaxGuests:   dto.MaxGuests,
		Photos:      photos,
		Commodities: dto.Commodities,
	}

	err := s.repo.Create(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (s *service) FindById(id uint) (*Room, error) {
	room, err := s.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("Not found")
	}
	return room, nil
}
