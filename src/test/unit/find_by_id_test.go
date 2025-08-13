package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindById_Success(t *testing.T) {
	svc, mockRepo, _ := createTestRoomService()

	room := DefaultRoom

	mockRepo.On("FindById", room.ID).Return(room, nil)

	roomGot, err := svc.FindById(room.ID)

	assert.NoError(t, err)
	assert.Equal(t, room, roomGot)
	mockRepo.AssertNumberOfCalls(t, "FindById", 1)
	mockRepo.AssertExpectations(t)
}

func Test_FindById_NotFound(t *testing.T) {
	svc, mockRepo, _ := createTestRoomService()

	room := DefaultRoom

	mockRepo.On("FindById", room.ID).Return(nil, fmt.Errorf("not found"))

	roomGot, err := svc.FindById(room.ID)

	assert.Error(t, err)
	assert.Nil(t, roomGot)
	mockRepo.AssertNumberOfCalls(t, "FindById", 1)
	mockRepo.AssertExpectations(t)
}
