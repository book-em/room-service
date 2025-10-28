package test

import (
	"bookem-room-service/internal"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DeleteRooomsByHostId_UserNotFound(t *testing.T) {
	svc, _, _, _, mockUserClient := CreateTestRoomService()

	hostId := uint(123)
	mockUserClient.On("FindById", context.Background(), hostId).Return(nil, fmt.Errorf("user not found"))

	rooms, err := svc.DeleteRoomsByHostId(context.Background(), hostId)

	assert.Error(t, err)
	assert.Nil(t, rooms)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}

func Test_DeleteRooomsByHostId_UserNotHost(t *testing.T) {
	svc, _, _, _, mockUserClient := CreateTestRoomService()

	notHost := DefaultUser_Guest
	mockUserClient.On("FindById", context.Background(), notHost.Id).Return(notHost, nil)

	rooms, err := svc.DeleteRoomsByHostId(context.Background(), notHost.Id)

	assert.Error(t, err)
	assert.Nil(t, rooms)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}

func Test_DeleteRooomsByHostId_FetchingRoomsDBError(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	host := DefaultUser_Host

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(nil, fmt.Errorf("db error"))

	rooms, err := svc.DeleteRoomsByHostId(context.Background(), host.Id)

	assert.Error(t, err)
	assert.Nil(t, rooms)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
}

func Test_DeleteRooomsByHostId_DeletingRoomsDbError(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	host := DefaultUser_Host
	room1 := DefaultRoom
	room1.HostID = host.Id
	rooms := []internal.Room{*room1}

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(rooms, nil)
	mockRepo.On("DeleteRoomsByHostId", host.Id).Return(fmt.Errorf("db error"))

	roomsGot, err := svc.DeleteRoomsByHostId(context.Background(), host.Id)

	assert.Error(t, err)
	assert.Nil(t, roomsGot)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "DeleteRoomsByHostId", 1)
	mockRepo.AssertExpectations(t)
}

func Test_DeleteRooomsByHostId_RefetchRoomsDbError(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	host := DefaultUser_Host
	room1 := DefaultRoom
	room1.HostID = host.Id
	rooms := []internal.Room{*room1}

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(rooms, nil).Once()
	mockRepo.On("DeleteRoomsByHostId", host.Id).Return(nil)
	mockRepo.On("FindByHost", host.Id).Return(nil, fmt.Errorf("db error"))

	roomsGot, err := svc.DeleteRoomsByHostId(context.Background(), host.Id)

	assert.Error(t, err)
	assert.Nil(t, roomsGot)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 2)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "DeleteRoomsByHostId", 1)
	mockRepo.AssertExpectations(t)
}

func Test_DeleteRooomsByHostId_NoRoomsSuccess(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	host := DefaultUser_Host
	rooms := []internal.Room{}

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(rooms, nil)

	roomsGot, err := svc.DeleteRoomsByHostId(context.Background(), host.Id)

	assert.NoError(t, err)
	assert.Equal(t, rooms, roomsGot)
	assert.Equal(t, 0, len(roomsGot))
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
}

func Test_DeleteRooomsByHostId_Success(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient := CreateTestRoomService()

	host := DefaultUser_Host
	room1 := DefaultRoom
	room2 := DefaultRoom
	room3 := DefaultRoom
	room1.HostID = host.Id
	room2.HostID = host.Id
	room3.HostID = host.Id
	rooms := []internal.Room{*room1, *room2, *room3}

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(rooms, nil)
	mockRepo.On("DeleteRoomsByHostId", host.Id).Return(nil)

	roomsGot, err := svc.DeleteRoomsByHostId(context.Background(), host.Id)

	assert.NoError(t, err)
	assert.Equal(t, rooms, roomsGot)
	assert.Equal(t, 3, len(roomsGot))
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 2)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "DeleteRoomsByHostId", 1)
	mockRepo.AssertExpectations(t)
}
