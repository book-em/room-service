package test

import (
	"bookem-room-service/client/userclient"
	test "bookem-room-service/test/unit"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_Create_Success(t *testing.T) {
	registerUser("user1", "1234", userclient.Host)
	jwt := loginUser2("user1", "1234")

	dto := test.DefaultRoomCreateDTO
	resp, err := createRoom(jwt, dto)
	require.NoError(t, err)

	result := responseToRoom(resp)
	require.Equal(t, fmt.Sprintf("room-%d-%d.jpg", result.ID, 0), result.Photos[0])
	require.Equal(t, result.Name, dto.Name)
	require.Equal(t, result.Address, dto.Address)
	require.Equal(t, result.Commodities, dto.Commodities)
}
