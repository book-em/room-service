package test

import (
	"bookem-room-service/internal"
	. "bookem-room-service/test/unit"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindPriceListById_Success(t *testing.T) {
	svc, _, _, mockPriceRepo, _ := CreateTestRoomService()

	list := DefaultPriceList

	mockPriceRepo.On("FindListById", list.ID).Return(list, nil)

	listGot, err := svc.FindPriceListById(list.ID)

	assert.NoError(t, err)
	assert.Equal(t, list, listGot)
	mockPriceRepo.AssertNumberOfCalls(t, "FindListById", 1)
	mockPriceRepo.AssertExpectations(t)
}

func Test_FindPriceListById_NotFound(t *testing.T) {
	svc, _, _, mockPriceRepo, _ := CreateTestRoomService()

	mockPriceRepo.On("FindListById", uint(999)).Return(nil, fmt.Errorf("not found"))

	listGot, err := svc.FindPriceListById(999)

	assert.Error(t, err)
	assert.Nil(t, listGot)
	mockPriceRepo.AssertNumberOfCalls(t, "FindListById", 1)
	mockPriceRepo.AssertExpectations(t)
}

func Test_FindPriceListsByRoomId_Success(t *testing.T) {
	svc, mockRepo, _, mockPriceRepo, _ := CreateTestRoomService()

	room := DefaultRoom
	room.ID = 1

	list1 := *DefaultPriceList
	list2 := *DefaultPriceList
	list1.ID = 1
	list2.ID = 2
	list1.RoomID = room.ID
	list2.RoomID = room.ID

	lists := []internal.RoomPriceList{list1, list2}

	mockRepo.On("FindById", room.ID).Return(room, nil)
	mockPriceRepo.On("FindListsByRoomId", room.ID).Return(lists, nil)

	listsGot, err := svc.FindPriceListsByRoomId(room.ID)

	assert.NoError(t, err)
	assert.Equal(t, lists, listsGot)
	mockRepo.AssertExpectations(t)
	mockPriceRepo.AssertExpectations(t)
}

func Test_FindPriceListsByRoomId_Success_NoLists(t *testing.T) {
	svc, mockRepo, _, mockPriceRepo, _ := CreateTestRoomService()

	room := DefaultRoom
	room.ID = 999

	mockRepo.On("FindById", room.ID).Return(room, nil)
	mockPriceRepo.On("FindListsByRoomId", room.ID).Return([]internal.RoomPriceList{}, nil)

	listsGot, err := svc.FindPriceListsByRoomId(room.ID)

	assert.NoError(t, err)
	assert.Empty(t, listsGot)
	mockRepo.AssertExpectations(t)
	mockPriceRepo.AssertExpectations(t)
}

func Test_FindPriceListsByRoomId_NotFound(t *testing.T) {
	svc, mockRepo, _, mockPriceRepo, _ := CreateTestRoomService()

	room := DefaultRoom
	room.ID = 1

	mockRepo.On("FindById", room.ID).Return(room, nil)
	mockPriceRepo.On("FindListsByRoomId", room.ID).Return(nil, fmt.Errorf("not found"))

	listsGot, err := svc.FindPriceListsByRoomId(room.ID)

	assert.Error(t, err)
	assert.Nil(t, listsGot)
	mockRepo.AssertExpectations(t)
	mockPriceRepo.AssertExpectations(t)
}

func Test_FindPriceListsByRoomId_RoomNotFound(t *testing.T) {
	svc, mockRepo, _, _, _ := CreateTestRoomService()

	mockRepo.On("FindById", uint(999)).Return(nil, fmt.Errorf("not found"))

	listsGot, err := svc.FindPriceListsByRoomId(999)

	assert.Error(t, err)
	assert.Nil(t, listsGot)
	mockRepo.AssertExpectations(t)
}

func Test_FindCurrentPriceListOfRoom_Success(t *testing.T) {
	svc, _, _, mockPriceRepo, _ := CreateTestRoomService()

	list := DefaultPriceList

	mockPriceRepo.On("FindCurrentListOfRoom", list.RoomID).Return(list, nil)

	listGot, err := svc.FindCurrentPriceListOfRoom(list.RoomID)

	assert.NoError(t, err)
	assert.Equal(t, list, listGot)
	mockPriceRepo.AssertExpectations(t)
}

func Test_FindCurrentPriceListOfRoom_NotFound(t *testing.T) {
	svc, _, _, mockPriceRepo, _ := CreateTestRoomService()

	mockPriceRepo.On("FindCurrentListOfRoom", uint(999)).Return(nil, fmt.Errorf("not found"))

	listGot, err := svc.FindCurrentPriceListOfRoom(999)

	assert.Error(t, err)
	assert.Nil(t, listGot)
	mockPriceRepo.AssertExpectations(t)
}
