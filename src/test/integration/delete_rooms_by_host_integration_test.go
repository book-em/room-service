package integration

import (
	"bookem-room-service/internal"
	"bookem-room-service/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_DeleteRoomsByHostId_NoRoomsSuccess(t *testing.T) {
	registerUser("host_del1", "1234", util.Host)
	jwt := loginUser2("host_del1", "1234")
	roomsExpect := []internal.RoomDTO{}

	resp, err := deleteRoomsByHostId(jwt)

	require.NoError(t, err)
	roomsGot := responseToRooms(resp)
	require.Equal(t, roomsExpect, roomsGot)
}

func TestIntegration_DeleteRoomsByHostId_InvalidRole(t *testing.T) {
	registerUser("guest_del1", "1234", util.Guest)
	jwt := loginUser2("guest_del1", "1234")

	resp, err := deleteRoomsByHostId(jwt)

	require.NoError(t, err)
	require.Equal(t, 401, resp.StatusCode)
}

func TestIntegration_DeleteRoomsByHostId_Success(t *testing.T) {
	username := "host_del2"
	jwt, _, room := createUserAndRoom(username)
	roomsExpect := []internal.RoomDTO{room}

	resp, err := deleteRoomsByHostId(jwt)

	require.NoError(t, err)
	roomsGot := responseToRooms(resp)
	require.Equal(t, roomsExpect[0].ID, roomsGot[0].ID)
	require.Equal(t, true, roomsGot[0].Deleted)
}
