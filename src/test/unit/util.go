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

var defaultRoom = &internal.Room{
	ID:     0,
	HostID: 0,

	Name:        "Room Name",
	Description: "Room Desc",
	Address:     "Room Address",
	MinGuests:   1,
	MaxGuests:   5,
	Photos:      []string{"img-1.png"},
	Commodities: []string{"test.png"},
}

var defaultRoomDTO = internal.RoomDTO{
	HostID:      defaultRoom.HostID,
	Name:        defaultRoom.Name,
	Description: defaultRoom.Description,
	Address:     defaultRoom.Address,
	MinGuests:   defaultRoom.MinGuests,
	MaxGuests:   defaultRoom.MaxGuests,
	Photos:      defaultRoom.Photos,
	Commodities: defaultRoom.Commodities,
}

var defaultRoomCreate = internal.CreateRoomDTO{
	HostID:        defaultRoom.HostID,
	Name:          defaultRoom.Name,
	Description:   defaultRoom.Description,
	Address:       defaultRoom.Address,
	MinGuests:     defaultRoom.MinGuests,
	MaxGuests:     defaultRoom.MaxGuests,
	PhotosPayload: []string{"some_base64_data_with_mime_type"},
	Commodities:   defaultRoom.Commodities,
}
