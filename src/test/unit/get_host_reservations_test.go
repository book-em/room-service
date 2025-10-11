package test

import (
	reservationClient "bookem-room-service/client/reservationclient"
	"bookem-room-service/internal"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetActiveHostReservations_UserNotFound(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient, _ := CreateTestRoomService()

	host := DefaultUser_Host
	jwt := "token"

	mockUserClient.On("FindById", context.Background(), host.Id).Return(nil, errors.New("User is not found"))

	reservations, err := svc.GetActiveHostReservations(context.Background(), host.Id, jwt)

	assert.Error(t, err)
	assert.Nil(t, reservations)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 0)
	mockRepo.AssertExpectations(t)
}

func Test_GetActiveHostReservations_FindRoomsError(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient, mockReservationClient := CreateTestRoomService()

	host := DefaultUser_Host
	jwt := "token"

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(nil, errors.New("Rooms are not found"))

	reservations, err := svc.GetActiveHostReservations(context.Background(), host.Id, jwt)

	assert.Error(t, err)
	assert.Nil(t, reservations)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
	mockReservationClient.AssertNumberOfCalls(t, "GetActiveHostReservations", 0)
	mockReservationClient.AssertExpectations(t)
}

func Test_GetActiveHostReservations_NoRoomsSuccess(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient, mockReservationClient := CreateTestRoomService()

	host := DefaultUser_Host
	jwt := "token"
	rooms := []internal.Room{}

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(rooms, nil)
	mockReservationClient.On("GetActiveHostReservations", context.Background(), jwt, mock.Anything).Return(nil, nil)

	reservations, err := svc.GetActiveHostReservations(context.Background(), host.Id, jwt)

	assert.NoError(t, err)
	assert.Nil(t, reservations)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
	mockReservationClient.AssertNumberOfCalls(t, "GetActiveHostReservations", 1)
	mockReservationClient.AssertExpectations(t)
}

func Test_GetActiveHostReservations_FoundRoomsNotFoundReservations(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient, mockReservationClient := CreateTestRoomService()

	host := DefaultUser_Host
	jwt := "token"
	room1 := DefaultRoom
	room2 := DefaultRoom
	room3 := DefaultRoom
	room1.HostID = host.Id
	room2.HostID = host.Id
	room3.HostID = host.Id
	rooms := []internal.Room{*room1, *room2, *room3}
	roomIds := []uint{room1.ID, room2.ID, room3.ID}

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(rooms, nil)
	mockReservationClient.On("GetActiveHostReservations", context.Background(), jwt, roomIds).Return(nil, nil)

	reservations, err := svc.GetActiveHostReservations(context.Background(), host.Id, jwt)

	assert.NoError(t, err)
	assert.Nil(t, reservations)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
	mockReservationClient.AssertNumberOfCalls(t, "GetActiveHostReservations", 1)
	mockReservationClient.AssertExpectations(t)
}

func Test_GetActiveHostReservations_ReservationsErr(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient, mockReservationClient := CreateTestRoomService()

	host := DefaultUser_Host
	jwt := "token"
	rooms := []internal.Room{}

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(rooms, nil)
	mockReservationClient.On("GetActiveHostReservations", context.Background(), jwt, mock.Anything).Return(nil, errors.New("Reservations db error"))

	reservations, err := svc.GetActiveHostReservations(context.Background(), host.Id, jwt)

	assert.Error(t, err)
	assert.Nil(t, reservations)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
	mockReservationClient.AssertNumberOfCalls(t, "GetActiveHostReservations", 1)
	mockReservationClient.AssertExpectations(t)
}

func Test_GetActiveHostReservations_Success(t *testing.T) {
	svc, mockRepo, _, _, mockUserClient, mockReservationClient := CreateTestRoomService()

	host := DefaultUser_Host
	jwt := "token"
	room1 := DefaultRoom
	room2 := DefaultRoom
	room3 := DefaultRoom
	room1.HostID = host.Id
	room2.HostID = host.Id
	room3.HostID = host.Id
	rooms := []internal.Room{*room1, *room2, *room3}
	roomIds := []uint{room1.ID, room2.ID, room3.ID}
	reservartion1 := DefaultReservationDTO
	reservartion2 := DefaultReservationDTO
	reservartion3 := DefaultReservationDTO
	reservations := []reservationClient.ReservationDTO{*reservartion1, *reservartion2, *reservartion3}

	mockUserClient.On("FindById", context.Background(), host.Id).Return(host, nil)
	mockRepo.On("FindByHost", host.Id).Return(rooms, nil)
	mockReservationClient.On("GetActiveHostReservations", context.Background(), jwt, roomIds).Return(reservations, nil)

	reservationsGot, err := svc.GetActiveHostReservations(context.Background(), host.Id, jwt)

	assert.NoError(t, err)
	assert.Equal(t, reservations, reservationsGot)
	mockUserClient.AssertNumberOfCalls(t, "FindById", 1)
	mockUserClient.AssertExpectations(t)
	mockRepo.AssertNumberOfCalls(t, "FindByHost", 1)
	mockRepo.AssertExpectations(t)
	mockReservationClient.AssertNumberOfCalls(t, "GetActiveHostReservations", 1)
	mockReservationClient.AssertExpectations(t)
}
