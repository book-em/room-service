package test

import (
	"bookem-room-service/internal"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindByHost_Success(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	host := DefaultUser_Host

	room1 := DefaultRoom
	room2 := DefaultRoom
	room3 := DefaultRoom
	room1.HostID = host.Id
	room2.HostID = host.Id
	room3.HostID = host.Id

	rooms := []internal.Room{*room1, *room2, *room3}

	mockUserClient.On("FindById", host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(rooms, nil)

	roomsGot, err := svc.FindByHost(context.Background(), host.Id)

	assert.NoError(t, err)
	assert.Equal(t, rooms, roomsGot)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
}

func Test_FindByHost_UserNotFound(t *testing.T) {
	svc, _, _, _, mockUserClient := CreateTestRoomService()

	hostId := uint(123)
	mockUserClient.On("FindById", hostId).Return(nil, fmt.Errorf("user not found"))

	roomsGot, err := svc.FindByHost(context.Background(), hostId)

	assert.Error(t, err)
	assert.Nil(t, roomsGot)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}

func Test_FindByHost_UserNotHost(t *testing.T) {
	svc, _, _, _, mockUserClient := CreateTestRoomService()

	notHost := DefaultUser_Guest
	mockUserClient.On("FindById", notHost.Id).Return(notHost, nil)

	roomsGot, err := svc.FindByHost(context.Background(), notHost.Id)

	assert.Error(t, err)
	assert.Nil(t, roomsGot)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}

func Test_FindByHost_DbError(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	host := DefaultUser_Host

	mockUserClient.On("FindById", host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(nil, fmt.Errorf("db error"))

	roomsGot, err := svc.FindByHost(context.Background(), host.Id)

	assert.Error(t, err)
	assert.Nil(t, roomsGot)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
}
