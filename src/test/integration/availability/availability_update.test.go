package test

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/internal"
	. "bookem-room-service/test/integration"
	test "bookem-room-service/test/unit"
	"bookem-room-service/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_UpdateAvailability_Success(t *testing.T) {
	username := "host_a_01"
	RegisterUser(username, "1234", userclient.Host)
	jwt := LoginUser2(username, "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)

	dto := test.DefaultRoomCreateDTO
	dto.HostID = jwtObj.ID
	resp, _ := CreateRoom(jwt, dto)
	room := ResponseToRoom(resp)

	createRoomDto := internal.CreateRoomAvailabilityListDTO{
		RoomID: room.ID,
		Items:  test.DefaultCreateAvailabilityListDTO.Items,
	}

	resp, err := CreateRoomAvailability(jwt, createRoomDto)
	if err != nil {
		panic(err)
	}
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	availList := ResponseToRoomAvailability(resp)

	require.Equal(t, 1, len(availList.Items))
	require.Equal(t, room.ID, availList.RoomID)
}
