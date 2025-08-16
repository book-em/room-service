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

func Test_UpdatePriceList_Success(t *testing.T) {
	svc, mockRepo, _, mockPriceRepo, mockUserClient := CreateTestRoomService()

	dto := DefaultCreatePriceListDTO
	user := DefaultUser_Host
	room := DefaultRoom
	room.HostID = user.Id

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(room, nil)
	mockPriceRepo.On("CreateList", mock.Anything).Return(nil)

	got, err := svc.UpdatePriceList(user.Id, dto)

	assert.NoError(t, err)
	assert.NotNil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockPriceRepo.AssertExpectations(t)
}

func Test_UpdatePriceList_UserNotFound(t *testing.T) {
	svc, _, _, _, mockUserClient := CreateTestRoomService()

	dto := DefaultCreatePriceListDTO
	hostID := uint(1234)

	mockUserClient.On("FindById", hostID).Return(nil, fmt.Errorf("not found"))

	got, err := svc.UpdatePriceList(hostID, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
}

func Test_UpdatePriceList_UserNotHost(t *testing.T) {
	svc, _, _, _, mockUserClient := CreateTestRoomService()

	dto := DefaultCreatePriceListDTO
	user := DefaultUser_Guest

	mockUserClient.On("FindById", user.Id).Return(user, nil)

	got, err := svc.UpdatePriceList(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
}

func Test_UpdatePriceList_RoomNotFound(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	dto := DefaultCreatePriceListDTO
	user := DefaultUser_Host

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(nil, fmt.Errorf("not found"))

	got, err := svc.UpdatePriceList(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_UpdatePriceList_HostNotOwnRoom(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	dto := DefaultCreatePriceListDTO
	user := DefaultUser_Host
	room := DefaultRoom
	room.HostID = user.Id + 1 // mismatch

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(room, nil)

	got, err := svc.UpdatePriceList(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_UpdatePriceList_BadDateRange(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	dto := internal.CreateRoomPriceListDTO{
		RoomID: DefaultRoom.ID,
		Items: []internal.CreateRoomPriceItemDTO{
			{
				DateFrom: time.Date(3025, 8, 20, 0, 0, 0, 0, time.UTC),
				DateTo:   time.Date(2020, 8, 25, 0, 0, 0, 0, time.UTC),
				Price:    100,
			},
		},
	}

	dto.Items[0].DateFrom = dto.Items[0].DateTo.Add(time.Hour) // invalid range

	user := DefaultUser_Host
	room := DefaultRoom
	room.HostID = user.Id

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(room, nil)

	got, err := svc.UpdatePriceList(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_UpdatePriceList_DuplicateDateRange(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	dto := internal.CreateRoomPriceListDTO{
		RoomID: DefaultRoom.ID,
		Items:  []internal.CreateRoomPriceItemDTO{DefaultCreatePriceItemDTO, DefaultCreatePriceItemDTO},
	}

	user := DefaultUser_Host
	room := DefaultRoom
	room.HostID = user.Id

	mockUserClient.On("FindById", user.Id).Return(user, nil)
	mockRepo.On("FindById", dto.RoomID).Return(room, nil)

	got, err := svc.UpdatePriceList(user.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, got)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
