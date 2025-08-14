package test

import (
	"bookem-room-service/client/userclient"
	test "bookem-room-service/test/unit"
	"bookem-room-service/util"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_Create_Success(t *testing.T) {
	RegisterUser("user1", "1234", userclient.Host)
	jwt := LoginUser2("user1", "1234")
	jwtObj, err := util.GetJwtFromString(jwt)
	if err != nil {
		panic(err)
	}

	dto := test.DefaultRoomCreateDTO
	dto.HostID = jwtObj.ID
	resp, err := CreateRoom(jwt, dto)
	require.NoError(t, err)

	result := ResponseToRoom(resp)
	require.Equal(t, fmt.Sprintf("room-%d-%d.jpg", result.ID, 0), result.Photos[0])
	require.Equal(t, result.Name, dto.Name)
	require.Equal(t, result.Address, dto.Address)
	require.Equal(t, result.Commodities, dto.Commodities)
}
