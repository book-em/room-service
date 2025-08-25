package test

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/internal"
	"time"

	mock "github.com/stretchr/testify/mock"
)

func CreateTestRoomService() (
	internal.Service,
	*MockRoomRepo,
	*MockRoomAvailabilityRepo,
	*MockRoomPriceRepo,
	*MockUserClient,
) {
	mockRepo := new(MockRoomRepo)
	mockRoomAvailRepo := new(MockRoomAvailabilityRepo)
	mockRoomPriceRepo := new(MockRoomPriceRepo)
	mockUserClient := new(MockUserClient)

	svc := internal.NewService(mockRepo, mockRoomAvailRepo, mockRoomPriceRepo, mockUserClient)
	return svc, mockRepo, mockRoomAvailRepo, mockRoomPriceRepo, mockUserClient
}

// ----------------------------------------------- Mock Room repo

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

func (r *MockRoomRepo) FindByFilters(guestsNumber uint, location string) ([]internal.Room, error) {
	args := r.Called(uint(guestsNumber), string(location))
	rooms, _ := args.Get(0).([]internal.Room)
	return rooms, args.Error(1)
}

// ----------------------------------------------- Mock room availabilty repo

type MockRoomAvailabilityRepo struct {
	mock.Mock
}

func (m *MockRoomAvailabilityRepo) CreateList(list *internal.RoomAvailabilityList) error {
	args := m.Called(list)
	return args.Error(0)
}

func (m *MockRoomAvailabilityRepo) FindListById(id uint) (*internal.RoomAvailabilityList, error) {
	args := m.Called(id)
	list, _ := args.Get(0).(*internal.RoomAvailabilityList)
	return list, args.Error(1)
}

func (m *MockRoomAvailabilityRepo) FindListsByRoomId(roomId uint) ([]internal.RoomAvailabilityList, error) {
	args := m.Called(roomId)
	lists, _ := args.Get(0).([]internal.RoomAvailabilityList)
	return lists, args.Error(1)
}

func (m *MockRoomAvailabilityRepo) FindCurrentListOfRoom(roomId uint) (*internal.RoomAvailabilityList, error) {
	args := m.Called(roomId)
	list, _ := args.Get(0).(*internal.RoomAvailabilityList)
	return list, args.Error(1)
}

// ----------------------------------------------- Mock price repo

type MockRoomPriceRepo struct {
	mock.Mock
}

func (m *MockRoomPriceRepo) CreateList(list *internal.RoomPriceList) error {
	args := m.Called(list)
	return args.Error(0)
}

func (m *MockRoomPriceRepo) FindListById(id uint) (*internal.RoomPriceList, error) {
	args := m.Called(id)
	list, _ := args.Get(0).(*internal.RoomPriceList)
	return list, args.Error(1)
}

func (m *MockRoomPriceRepo) FindListsByRoomId(roomId uint) ([]internal.RoomPriceList, error) {
	args := m.Called(roomId)
	lists, _ := args.Get(0).([]internal.RoomPriceList)
	return lists, args.Error(1)
}

func (m *MockRoomPriceRepo) FindCurrentListOfRoom(roomId uint) (*internal.RoomPriceList, error) {
	args := m.Called(roomId)
	list, _ := args.Get(0).(*internal.RoomPriceList)
	return list, args.Error(1)
}

// ----------------------------------------------- Mock user client

type MockUserClient struct {
	mock.Mock
}

func (r *MockUserClient) FindById(id uint) (*userclient.UserDTO, error) {
	args := r.Called(id)
	user, _ := args.Get(0).(*userclient.UserDTO)
	return user, args.Error(1)
}

// ----------------------------------------------- Mock data

const (
	SMALL_IMG = "data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEASABIAAD/2wBDAAMCAgMCAgMDAwMEAwMEBQgFBQQEBQoHBwYIDAoMDAsKCwsNDhIQDQ4RDgsLEBYQERMUFRUVDA8XGBYUGBIUFRT/wAALCAABAAEBAREA/8QAFAABAAAAAAAAAAAAAAAAAAAACf/EABQQAQAAAAAAAAAAAAAAAAAAAAD/2gAIAQEAAD8AKp//2Q=="
)

var DefaultRoom = &internal.Room{
	ID:     0,
	HostID: 2,

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

var DefaultUser_Guest = &userclient.UserDTO{
	Id:       1,
	Username: "guser",
	Email:    "gemail@mail.com",
	Name:     "gname",
	Surname:  "gsurname",
	Role:     "guest",
	Address:  "gAddress 123",
}

var DefaultUser_Host = &userclient.UserDTO{
	Id:       2,
	Username: "huser",
	Email:    "hemail@mail.com",
	Name:     "hname",
	Surname:  "hsurname",
	Role:     "host",
	Address:  "hAddress 123",
}
var DefaultAvailabilityItem = internal.RoomAvailabilityItem{
	ID:        1,
	DateFrom:  time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
	DateTo:    time.Date(2025, 8, 25, 0, 0, 0, 0, time.UTC),
	Available: true,
}

var DefaultAvailabilityList = &internal.RoomAvailabilityList{
	ID:            1,
	RoomID:        DefaultRoom.ID,
	EffectiveFrom: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
	Items:         []internal.RoomAvailabilityItem{DefaultAvailabilityItem},
}

var DefaultAvailabilityItemDTO = internal.RoomAvailabilityItemDTO{
	ID:        DefaultAvailabilityItem.ID,
	DateFrom:  DefaultAvailabilityItem.DateFrom,
	DateTo:    DefaultAvailabilityItem.DateTo,
	Available: DefaultAvailabilityItem.Available,
}

var DefaultAvailabilityListDTO = internal.RoomAvailabilityListDTO{
	ID:            DefaultAvailabilityList.ID,
	RoomID:        DefaultAvailabilityList.RoomID,
	EffectiveFrom: DefaultAvailabilityList.EffectiveFrom,
	Items:         []internal.RoomAvailabilityItemDTO{DefaultAvailabilityItemDTO},
}

var DefaultCreateAvailabilityItemDTO = internal.CreateRoomAvailabilityItemDTO{
	DateFrom:  DefaultAvailabilityItem.DateFrom,
	DateTo:    DefaultAvailabilityItem.DateTo,
	Available: DefaultAvailabilityItem.Available,
}

var DefaultCreateAvailabilityListDTO = internal.CreateRoomAvailabilityListDTO{
	RoomID: DefaultRoom.ID,
	Items:  []internal.CreateRoomAvailabilityItemDTO{DefaultCreateAvailabilityItemDTO},
}

var DefaultPriceItem = internal.RoomPriceItem{
	ID:       1,
	DateFrom: time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
	DateTo:   time.Date(2025, 8, 25, 0, 0, 0, 0, time.UTC),
	Price:    100,
}

var DefaultPriceList = &internal.RoomPriceList{
	ID:            1,
	RoomID:        DefaultRoom.ID,
	EffectiveFrom: time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
	BasePrice:     100,
	PerGuest:      true,
	Items:         []internal.RoomPriceItem{DefaultPriceItem},
}

var DefaultPriceItemDTO = internal.RoomPriceItemDTO{
	ID:       DefaultPriceItem.ID,
	DateFrom: DefaultPriceItem.DateFrom,
	DateTo:   DefaultPriceItem.DateTo,
	Price:    DefaultPriceItem.Price,
}

var DefaultPriceListDTO = internal.RoomPriceListDTO{
	ID:            DefaultPriceList.ID,
	RoomID:        DefaultPriceList.RoomID,
	EffectiveFrom: DefaultPriceList.EffectiveFrom,
	BasePrice:     DefaultPriceList.BasePrice,
	PerGuest:      DefaultPriceList.PerGuest,
	Items:         []internal.RoomPriceItemDTO{DefaultPriceItemDTO},
}

var DefaultCreatePriceItemDTO = internal.CreateRoomPriceItemDTO{
	DateFrom: DefaultPriceItem.DateFrom,
	DateTo:   DefaultPriceItem.DateTo,
	Price:    DefaultPriceItem.Price,
}

var DefaultCreatePriceListDTO = internal.CreateRoomPriceListDTO{
	RoomID: DefaultRoom.ID,
	Items:  []internal.CreateRoomPriceItemDTO{DefaultCreatePriceItemDTO},
}

var DefaultRoomsQueryDTO = &internal.RoomsQueryDTO{
	Address:      "address",
	GuestsNumber: 4,
	DateFrom:     time.Date(2025, 8, 6, 0, 0, 0, 0, time.UTC),
	DateTo:       time.Date(2025, 8, 7, 0, 0, 0, 0, time.UTC),
	PageNumber:   1,
	PageSize:     10,
}

var DefaultRoomResult = &internal.RoomResultDTO{
	ID:          1,
	Name:        "Room Name",
	Description: "Room Desc",
	Address:     "Room Address",
	Photos:      []string{"test.png"},
	UnitPrice:   100.0,
	TotalPrice:  200.0,
}

var DefaulPaginatedResultInfoDTO = &internal.PaginatedResultInfoDTO{
	Page:       1,
	PageSize:   10,
	TotalPages: 1,
	TotalHits:  3,
}
