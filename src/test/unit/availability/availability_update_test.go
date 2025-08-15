package test

import (
	"bookem-room-service/internal"
	. "bookem-room-service/test/unit"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_UpdateAvailability_Success(t *testing.T) {
	svc, mockRepo, mockAvailRepo, mockUserClient := CreateTestRoomService()

	dto := DefaultCreateAvailabilityListDTO
	user := DefaultUser_Host
	room := DefaultRoom
	room.HostID = user.Id

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(room, nil)
	mockAvailRepo.On("CreateList", mock.Anything).Return(nil)

	got, err := svc.UpdateAvailability(user.Id, dto)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockAvailRepo.AssertExpectations(t)
}

func Test_UpdateAvailability_UserNotFound(t *testing.T) {
	svc, _, _, mockUserClient := CreateTestRoomService()

	dto := DefaultCreateAvailabilityListDTO
	hostID := uint(1234)

	mockUserClient.On("FindById", hostID).Return(nil, fmt.Errorf("not found"))

	got, err := svc.UpdateAvailability(hostID, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
}

func Test_UpdateAvailability_UserNotHost(t *testing.T) {
	svc, _, _, mockUserClient := CreateTestRoomService()

	dto := DefaultCreateAvailabilityListDTO
	user := DefaultUser_Guest

	mockUserClient.On("FindById", user.Id).Return(user, nil)

	got, err := svc.UpdateAvailability(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
}

func Test_UpdateAvailability_RoomNotFound(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	dto := DefaultCreateAvailabilityListDTO
	user := DefaultUser_Host

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(nil, fmt.Errorf("not found"))

	got, err := svc.UpdateAvailability(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_UpdateAvailability_HostNotOwnRoom(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	dto := DefaultCreateAvailabilityListDTO
	user := DefaultUser_Host
	room := DefaultRoom
	room.HostID = user.Id + 1 // mismatch

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(room, nil)

	got, err := svc.UpdateAvailability(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_UpdateAvailability_BadDateRange(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	// We're creating a new one here because this:
	//
	// dto := DefaultCreateAvailabilityListDTO
	//
	// creates a shallow copy, so modifying dto.Items will affect the global
	// object and make later tests fail.
	dto := internal.CreateRoomAvailabilityListDTO{
		RoomID: DefaultRoom.ID,
		Items: []internal.CreateRoomAvailabilityItemDTO{internal.CreateRoomAvailabilityItemDTO{
			DateFrom:  time.Date(3025, 8, 20, 0, 0, 0, 0, time.UTC),
			DateTo:    time.Date(2020, 8, 25, 0, 0, 0, 0, time.UTC),
			Available: DefaultAvailabilityItem.Available,
		}},
	}

	dto.Items[0].DateFrom = dto.Items[0].DateTo.Add(time.Hour) // invalid range

	user := DefaultUser_Host
	room := DefaultRoom
	room.HostID = user.Id

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(room, nil)

	got, err := svc.UpdateAvailability(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_UpdateAvailability_DuplicateDateRange(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	dto := internal.CreateRoomAvailabilityListDTO{
		RoomID: DefaultRoom.ID,
		Items:  []internal.CreateRoomAvailabilityItemDTO{DefaultCreateAvailabilityItemDTO, DefaultCreateAvailabilityItemDTO},
	}

	user := DefaultUser_Host
	room := DefaultRoom
	room.HostID = user.Id

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(room, nil)

	got, err := svc.UpdateAvailability(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_UpdateAvailability_DBError(t *testing.T) {
	svc, mockRepo, mockAvailRepo, mockUserClient := CreateTestRoomService()

	dto := DefaultCreateAvailabilityListDTO
	user := DefaultUser_Host
	room := DefaultRoom
	room.HostID = user.Id

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(room, nil)
	mockAvailRepo.On("CreateList", mock.Anything).Return(fmt.Errorf("db error"))

	got, err := svc.UpdateAvailability(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockAvailRepo.AssertExpectations(t)
}
