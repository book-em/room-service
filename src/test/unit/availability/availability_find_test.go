package test

import (
	"bookem-room-service/internal"
	. "bookem-room-service/test/unit"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindAvailabilityListById_Success(t *testing.T) {
	svc, _, mockAvailRepo, _, _ := CreateTestRoomService()

	li := DefaultAvailabilityList

	mockAvailRepo.On("FindListById", li.ID).Return(li, nil)

	liGot, err := svc.FindAvailabilityListById(li.ID)

	assert.NoError(t, err)
	assert.Equal(t, li, liGot)
	mockAvailRepo.AssertNumberOfCalls(t, "FindListById", 1)
	mockAvailRepo.AssertExpectations(t)
}

func Test_FindAvailabilityListById_NotFound(t *testing.T) {
	svc, _, mockAvailRepo, _, _ := CreateTestRoomService()

	li := DefaultAvailabilityList

	mockAvailRepo.On("FindListById", li.ID).Return(nil, fmt.Errorf("not found"))

	liGot, err := svc.FindAvailabilityListById(li.ID)

	assert.Error(t, err)
	assert.Nil(t, liGot)
	mockAvailRepo.AssertNumberOfCalls(t, "FindListById", 1)
	mockAvailRepo.AssertExpectations(t)
}

func Test_FindAvailabilityListsByRoomId_Success(t *testing.T) {
	svc, mockRepo, mockAvailRepo, _, _ := CreateTestRoomService()

	room := DefaultRoom
	room.ID = uint(1)
	li1 := *DefaultAvailabilityList
	li2 := *DefaultAvailabilityList
	li1.ID = 1
	li2.ID = 2
	li1.RoomID = room.ID
	li2.RoomID = 2

	lists := []internal.RoomAvailabilityList{li1, li2}

	mockAvailRepo.On("FindListsByRoomId", li1.RoomID).Return(lists, nil)
	mockRepo.On("FindById", room.ID).Return(room, nil)

	listsGot, err := svc.FindAvailabilityListsByRoomId(room.ID)

	assert.NoError(t, err)
	assert.Equal(t, lists, listsGot)
	mockAvailRepo.AssertNumberOfCalls(t, "FindListsByRoomId", 1)
	mockAvailRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_FindAvailabilityListsByRoomId_Success_NoLists(t *testing.T) {
	svc, mockRepo, mockAvailRepo, _, _ := CreateTestRoomService()

	room := DefaultRoom
	room.ID = uint(999)

	lists := []internal.RoomAvailabilityList{}

	mockAvailRepo.On("FindListsByRoomId", uint(999)).Return(lists, nil)
	mockRepo.On("FindById", room.ID).Return(room, nil)

	listsGot, err := svc.FindAvailabilityListsByRoomId(uint(999))

	assert.NoError(t, err)
	assert.Equal(t, lists, listsGot)
	mockAvailRepo.AssertNumberOfCalls(t, "FindListsByRoomId", 1)
	mockAvailRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_FindAvailabilityListsByRoomId_NotFound(t *testing.T) {
	svc, mockRepo, mockAvailRepo, _, _ := CreateTestRoomService()

	room := DefaultRoom
	room.ID = uint(1)

	li := DefaultAvailabilityList
	li.RoomID = room.ID

	mockAvailRepo.On("FindListsByRoomId", li.RoomID).Return(nil, fmt.Errorf("not found"))
	mockRepo.On("FindById", room.ID).Return(room, nil)

	liGot, err := svc.FindAvailabilityListsByRoomId(li.RoomID)

	assert.Error(t, err)
	assert.Nil(t, liGot)
	mockAvailRepo.AssertNumberOfCalls(t, "FindListsByRoomId", 1)
	mockAvailRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func Test_FindAvailabilityListsByRoomId_RoomNotFound(t *testing.T) {
	svc, mockRepo, _, _, _ := CreateTestRoomService()

	li := DefaultAvailabilityList

	mockRepo.On("FindById", li.RoomID).Return(nil, fmt.Errorf("not found"))

	liGot, err := svc.FindAvailabilityListsByRoomId(li.RoomID)

	assert.Error(t, err)
	assert.Nil(t, liGot)
	mockRepo.AssertExpectations(t)
}

func Test_FindCurrentAvailabilityListOfRoom_Success(t *testing.T) {
	svc, _, mockAvailRepo, _, _ := CreateTestRoomService()

	li := DefaultAvailabilityList

	mockAvailRepo.On("FindCurrentListOfRoom", li.RoomID).Return(li, nil)

	liGot, err := svc.FindCurrentAvailabilityListOfRoom(li.RoomID)

	assert.NoError(t, err)
	assert.Equal(t, li, liGot)
	mockAvailRepo.AssertNumberOfCalls(t, "FindCurrentListOfRoom", 1)
	mockAvailRepo.AssertExpectations(t)
}

func Test_FindCurrentAvailabilityListOfRoom_NotFound(t *testing.T) {
	svc, _, mockAvailRepo, _, _ := CreateTestRoomService()

	li := DefaultAvailabilityList

	mockAvailRepo.On("FindCurrentListOfRoom", li.RoomID).Return(nil, fmt.Errorf("not found"))

	liGot, err := svc.FindCurrentAvailabilityListOfRoom(li.RoomID)

	assert.Error(t, err)
	assert.Nil(t, liGot)
	mockAvailRepo.AssertNumberOfCalls(t, "FindCurrentListOfRoom", 1)
	mockAvailRepo.AssertExpectations(t)
}
