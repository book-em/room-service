package test

import (
	"bookem-room-service/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func Test_Create_Success(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	room := DefaultRoom
	dto := DefaultRoomCreateDTO

	mockRepo.On("Create", mock.AnythingOfType("*internal.Room")).Return(nil)
	mockRepo.On("Update", mock.AnythingOfType("*internal.Room")).Return(nil)
	mockUserClient.On("FindById", mock.AnythingOfType("uint")).Return(DefaultUser_Host, nil)
	util.SaveImageB64 = func(base64Image string, filename string) (string, string, error) {
		return "foo/" + room.Photos[0], room.Photos[0], nil
	}

	roomGot, err := svc.Create(DefaultUser_Host.Id, dto)

	assert.NoError(t, err)
	assert.Equal(t, room, roomGot)
	mockRepo.AssertNumberOfCalls(t, "Create", 1)
	mockRepo.AssertNumberOfCalls(t, "Update", 1)
	mockRepo.AssertNumberOfCalls(t, "Delete", 0)
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}

func Test_Create_InsertFailed(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	dto := DefaultRoomCreateDTO

	mockRepo.On("Create", mock.AnythingOfType("*internal.Room")).Return(fmt.Errorf("db error"))
	mockUserClient.On("FindById", mock.AnythingOfType("uint")).Return(DefaultUser_Host, nil)

	roomGot, err := svc.Create(DefaultUser_Host.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, roomGot)
	mockRepo.AssertNumberOfCalls(t, "Create", 1)
	mockRepo.AssertNumberOfCalls(t, "Update", 0)
	mockRepo.AssertNumberOfCalls(t, "Delete", 0)
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}

func Test_Create_ImageSaveFailed(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	dto := DefaultRoomCreateDTO

	mockRepo.On("Create", mock.AnythingOfType("*internal.Room")).Return(nil)
	mockRepo.On("Delete", mock.AnythingOfType("*internal.Room")).Return(nil)
	mockUserClient.On("FindById", mock.AnythingOfType("uint")).Return(DefaultUser_Host, nil)

	util.SaveImageB64 = func(base64Image string, filename string) (string, string, error) {
		return "", "", fmt.Errorf("some error")
	}

	roomGot, err := svc.Create(DefaultUser_Host.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, roomGot)
	mockRepo.AssertNumberOfCalls(t, "Create", 1)
	mockRepo.AssertNumberOfCalls(t, "Update", 0)
	mockRepo.AssertNumberOfCalls(t, "Delete", 1)
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}

func Test_Create_UpdateFailed(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	room := DefaultRoom
	dto := DefaultRoomCreateDTO

	mockRepo.On("Create", mock.AnythingOfType("*internal.Room")).Return(nil)
	mockRepo.On("Update", mock.AnythingOfType("*internal.Room")).Return(fmt.Errorf("error"))
	mockRepo.On("Delete", mock.AnythingOfType("*internal.Room")).Return(nil)
	mockUserClient.On("FindById", mock.AnythingOfType("uint")).Return(DefaultUser_Host, nil)
	util.SaveImageB64 = func(base64Image string, filename string) (string, string, error) {
		return "foo/" + room.Photos[0], room.Photos[0], nil
	}

	roomGot, err := svc.Create(DefaultUser_Host.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, roomGot)
	mockRepo.AssertNumberOfCalls(t, "Create", 1)
	mockRepo.AssertNumberOfCalls(t, "Update", 1)
	mockRepo.AssertNumberOfCalls(t, "Delete", 1)
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}

func Test_Create_HostNotFound(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	dto := DefaultRoomCreateDTO

	mockUserClient.On("FindById", mock.AnythingOfType("uint")).Return(nil, fmt.Errorf("user not found"))
	roomGot, err := svc.Create(DefaultUser_Host.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, roomGot)
	mockRepo.AssertNumberOfCalls(t, "Create", 0)
	mockRepo.AssertNumberOfCalls(t, "Update", 0)
	mockRepo.AssertNumberOfCalls(t, "Delete", 0)
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}

func Test_Create_HostHasBadRole(t *testing.T) {
	svc, mockRepo, _, mockUserClient := CreateTestRoomService()

	dto := DefaultRoomCreateDTO

	mockUserClient.On("FindById", mock.AnythingOfType("uint")).Return(DefaultUser_Guest, nil)
	roomGot, err := svc.Create(DefaultUser_Guest.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, roomGot)
	mockRepo.AssertNumberOfCalls(t, "Create", 0)
	mockRepo.AssertNumberOfCalls(t, "Update", 0)
	mockRepo.AssertNumberOfCalls(t, "Delete", 0)
	mockRepo.AssertExpectations(t)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
}
