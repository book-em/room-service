package test

import (
	internal "bookem-room-service/internal"

	mock "github.com/stretchr/testify/mock"
)

func createTestRoomService() (internal.Service, *MockRoomRepo) {
	mockRepo := new(MockRoomRepo)
	svc := internal.NewService(mockRepo)
	return svc, mockRepo
}

// ----------------------------------------------- Mock repo

type MockRoomRepo struct {
	mock.Mock
}

func (r *MockRoomRepo) Create(room *internal.Room) error {
	args := r.Called(room)
	return args.Error(0)
}

func (r *MockRoomRepo) Update(room *internal.Room) error {
	args := r.Called(room)
	return args.Error(0)
}

func (r *MockRoomRepo) Delete(room *internal.Room) error {
	args := r.Called(room)
	return args.Error(0)
}

func (r *MockRoomRepo) FindById(id uint) (*internal.Room, error) {
	args := r.Called(uint(id))
	user, _ := args.Get(0).(*internal.Room)
	return user, args.Error(1)
}

func (r *MockRoomRepo) FindByHost(hostId uint) ([]internal.Room, error) {
	args := r.Called(uint(hostId))
	user, _ := args.Get(0).([]internal.Room)
	return user, args.Error(1)
}

// ----------------------------------------------- Mock data

const (
	SMALL_IMG = "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAMCAgMCAgMDAwMEAwMEBQgFBQQEBQoHBwYIDAoMDAsKCwsNDhIQDQ4RDgsLEBYQERMUFRUVDA8XGBYUGBIUFRT/wAALCAABAAEBAREA/8QAFAABAAAAAAAAAAAAAAAAAAAACf/EABQQAQAAAAAAAAAAAAAAAAAAAAD/2gAIAQEAAD8AKp//2Q=="
)

var DefaultRoom = &internal.Room{
	ID:     0,
	HostID: 0,

	Name:        "Room Name",
	Description: "Room Desc",
	Address:     "Room Address",
	MinGuests:   1,
	MaxGuests:   5,
	Photos:      []string{"test.png"},
	Commodities: []string{"WiFi"},
}

var DefaultRoomDTO = internal.RoomDTO{
	HostID:      DefaultRoom.HostID,
	Name:        DefaultRoom.Name,
	Description: DefaultRoom.Description,
	Address:     DefaultRoom.Address,
	MinGuests:   DefaultRoom.MinGuests,
	MaxGuests:   DefaultRoom.MaxGuests,
	Photos:      DefaultRoom.Photos,
	Commodities: DefaultRoom.Commodities,
}

var DefaultRoomCreateDTO = internal.CreateRoomDTO{
	HostID:        DefaultRoom.HostID,
	Name:          DefaultRoom.Name,
	Description:   DefaultRoom.Description,
	Address:       DefaultRoom.Address,
	MinGuests:     DefaultRoom.MinGuests,
	MaxGuests:     DefaultRoom.MaxGuests,
	PhotosPayload: []string{SMALL_IMG},
	Commodities:   DefaultRoom.Commodities,
}
