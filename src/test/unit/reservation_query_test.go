package test

import (
	"bookem-room-service/internal"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_QueryForReservation_Success(t *testing.T) {
	svc, mockRoomRepo, mockAvailRepo, mockRoomPriceRepo, mockUserClient := CreateTestRoomService()

	dto := internal.RoomReservationQueryDTO{
		RoomID:     DefaultRoom.ID,
		DateFrom:   time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		DateTo:     time.Date(2025, 8, 21, 0, 0, 0, 0, time.UTC),
		GuestCount: 2,
	}

	mockUserClient.On("FindById", DefaultUser_Guest.Id).Return(DefaultUser_Guest, nil)
	mockRoomRepo.On("FindById", DefaultRoom.ID).Return(DefaultRoom, nil)
	mockAvailRepo.On("FindCurrentListOfRoom", DefaultRoom.ID).Return(DefaultAvailabilityList, nil)
	mockRoomPriceRepo.On("FindCurrentListOfRoom", DefaultRoom.ID).Return(DefaultPriceList, nil)

	resp, err := svc.QueryForReservation(context.Background(), DefaultUser_Guest.Id, dto)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Available)
	assert.Equal(t, uint(400), resp.TotalCost) // 2 days × 100 x 2 guests

	mockUserClient.AssertExpectations(t)
	mockRoomRepo.AssertExpectations(t)
	mockAvailRepo.AssertExpectations(t)
	mockRoomPriceRepo.AssertExpectations(t)
}

func Test_QueryForReservation_RoomUnavailable(t *testing.T) {
	svc, mockRoomRepo, mockAvailRepo, _, mockUserClient := CreateTestRoomService()

	dto := internal.RoomReservationQueryDTO{
		RoomID:     DefaultRoom.ID,
		DateFrom:   time.Date(2025, 8, 26, 0, 0, 0, 0, time.UTC), // outside availability
		DateTo:     time.Date(2025, 8, 27, 0, 0, 0, 0, time.UTC),
		GuestCount: 2,
	}

	mockUserClient.On("FindById", DefaultUser_Guest.Id).Return(DefaultUser_Guest, nil)
	mockRoomRepo.On("FindById", DefaultRoom.ID).Return(DefaultRoom, nil)
	mockAvailRepo.On("FindCurrentListOfRoom", DefaultRoom.ID).Return(DefaultAvailabilityList, nil)

	resp, err := svc.QueryForReservation(context.Background(), DefaultUser_Guest.Id, dto)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Available)
	assert.Equal(t, uint(0), resp.TotalCost)

	mockUserClient.AssertExpectations(t)
	mockRoomRepo.AssertExpectations(t)
	mockAvailRepo.AssertExpectations(t)
}

func Test_QueryForReservation_PriceCalculationFails(t *testing.T) {
	svc, mockRoomRepo, mockAvailRepo, mockRoomPriceRepo, mockUserClient := CreateTestRoomService()

	dto := internal.RoomReservationQueryDTO{
		RoomID:     DefaultRoom.ID,
		DateFrom:   DefaultPriceItem.DateFrom,
		DateTo:     DefaultPriceItem.DateTo,
		GuestCount: 2,
	}

	mockUserClient.On("FindById", DefaultUser_Guest.Id).Return(DefaultUser_Guest, nil)
	mockRoomRepo.On("FindById", DefaultRoom.ID).Return(DefaultRoom, nil)
	mockAvailRepo.On("FindCurrentListOfRoom", DefaultRoom.ID).Return(DefaultAvailabilityList, nil)
	mockRoomPriceRepo.On("FindCurrentListOfRoom", DefaultRoom.ID).Return(nil, fmt.Errorf("pricing error"))

	resp, err := svc.QueryForReservation(context.Background(), DefaultUser_Guest.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockUserClient.AssertExpectations(t)
	mockRoomRepo.AssertExpectations(t)
	mockAvailRepo.AssertExpectations(t)
	mockRoomPriceRepo.AssertExpectations(t)
}

func Test_QueryForReservation_UnauthorizedUser(t *testing.T) {
	svc, _, _, _, mockUserClient := CreateTestRoomService()

	dto := internal.RoomReservationQueryDTO{
		RoomID:     DefaultRoom.ID,
		DateFrom:   DefaultPriceItem.DateFrom,
		DateTo:     DefaultPriceItem.DateTo,
		GuestCount: 2,
	}

	mockUserClient.On("FindById", DefaultUser_Host.Id).Return(DefaultUser_Host, nil)

	resp, err := svc.QueryForReservation(context.Background(), DefaultUser_Host.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, internal.ErrUnauthorized, err)

	mockUserClient.AssertExpectations(t)
}

func Test_QueryForReservation_UserNotFound(t *testing.T) {
	svc, _, _, _, mockUserClient := CreateTestRoomService()

	dto := internal.RoomReservationQueryDTO{
		RoomID:     DefaultRoom.ID,
		DateFrom:   DefaultPriceItem.DateFrom,
		DateTo:     DefaultPriceItem.DateTo,
		GuestCount: 2,
	}

	mockUserClient.On("FindById", uint(999)).Return(nil, fmt.Errorf("user not found"))

	resp, err := svc.QueryForReservation(context.Background(), 999, dto)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockUserClient.AssertExpectations(t)
}

func Test_QueryForReservation_RoomNotFound(t *testing.T) {
	svc, mockRoomRepo, _, _, mockUserClient := CreateTestRoomService()

	dto := internal.RoomReservationQueryDTO{
		RoomID:     999,
		DateFrom:   DefaultPriceItem.DateFrom,
		DateTo:     DefaultPriceItem.DateTo,
		GuestCount: 2,
	}

	mockUserClient.On("FindById", DefaultUser_Guest.Id).Return(DefaultUser_Guest, nil)
	mockRoomRepo.On("FindById", uint(999)).Return(nil, fmt.Errorf("room not found"))

	resp, err := svc.QueryForReservation(context.Background(), DefaultUser_Guest.Id, dto)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockUserClient.AssertExpectations(t)
	mockRoomRepo.AssertExpectations(t)
}
