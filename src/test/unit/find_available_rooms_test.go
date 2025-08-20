package test

import (
	"bookem-room-service/internal"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_FindAvailableRooms_Success(t *testing.T) {
	svc, mockRepo, _, _, _ := CreateTestRoomService()

	room1 := DefaultRoomResult
	room2 := DefaultRoomResult
	room3 := DefaultRoomResult
	info := DefaulPaginatedResultInfoDTO
	d := *DefaultRoomsQueryDTO

	rooms := []internal.RoomResultDTO{*room1, *room2, *room3}
	var totalHits int64 = 3

	mockRepo.
		On("FindAvailableRooms", d.Location, d.GuestsNumber, d.DateFrom, d.DateTo, d.PageNumber, d.PageSize).
		Return(rooms, totalHits, nil)

	roomsGot, infoGot, err := svc.FindAvailableRooms(d)

	assert.NoError(t, err)
	assert.Equal(t, rooms, roomsGot)
	assert.Equal(t, info, infoGot)
	mockRepo.AssertNumberOfCalls(t, "FindAvailableRooms", 1)
	mockRepo.AssertExpectations(t)
}

func Test_FindAvailableRooms_InvalidDate(t *testing.T) {
	svc, mockRepo, _, _, _ := CreateTestRoomService()

	d := *DefaultRoomsQueryDTO
	d.DateFrom = d.DateTo.Add(24 * time.Hour)

	roomsGot, infoGot, err := svc.FindAvailableRooms(d)

	assert.Nil(t, roomsGot)
	assert.Nil(t, infoGot)
	assert.Error(t, err)
	mockRepo.AssertNumberOfCalls(t, "FindAvailableRooms", 0)
	mockRepo.AssertExpectations(t)
}

func Test_FindAvailableRooms_DbError(t *testing.T) {
	svc, mockRepo, _, _, _ := CreateTestRoomService()

	d := *DefaultRoomsQueryDTO

	mockRepo.
		On("FindAvailableRooms", d.Location, d.GuestsNumber, d.DateFrom, d.DateTo, d.PageNumber, d.PageSize).
		Return(nil, nil, fmt.Errorf("db error"))

	roomsGot, infoGot, err := svc.FindAvailableRooms(d)

	assert.Nil(t, roomsGot)
	assert.Nil(t, infoGot)
	assert.Error(t, err)
	mockRepo.AssertNumberOfCalls(t, "FindAvailableRooms", 1)
	mockRepo.AssertExpectations(t)
}
