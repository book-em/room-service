package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindById_Success(t *testing.T) {
	svc, mockRepo, _, _, _ := CreateTestRoomService()

	room := DefaultRoom

	mockRepo.On("FindById", room.ID).Return(room, nil)

	roomGot, err := svc.FindById(context.Background(), room.ID)

	assert.NoError(t, err)
	assert.Equal(t, room, roomGot)
	mockRepo.AssertNumberOfCalls(t, "FindById", 1)
	mockRepo.AssertExpectations(t)
}

func Test_FindById_NotFound(t *testing.T) {
	svc, mockRepo, _, _, _ := CreateTestRoomService()

	room := DefaultRoom

	mockRepo.On("FindById", room.ID).Return(nil, fmt.Errorf("not found"))

	roomGot, err := svc.FindById(context.Background(), room.ID)

	assert.Error(t, err)
	assert.Nil(t, roomGot)
	mockRepo.AssertNumberOfCalls(t, "FindById", 1)
	mockRepo.AssertExpectations(t)
}
