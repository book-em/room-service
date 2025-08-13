package test

import (
	"bookem-room-service/client/userclient"
	test "bookem-room-service/test/unit"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_FindById_Success(t *testing.T) {
	registerUser("user2", "1234", userclient.Host)
	jwt := loginUser2("user2", "1234")
	resp, _ := createRoom(jwt, test.DefaultRoomCreateDTO)
	room := responseToRoom(resp)

	roomId := room.ID

	resp, err := findRoomById(roomId)

	require.NoError(t, err)
	roomGot := responseToRoom(resp)

	require.Equal(t, roomId, roomGot.ID)
	require.Equal(t, room, roomGot)
}

func TestIntegration_FindById_MissingId(t *testing.T) {
	resp, err := http.Get(URL_room)

	require.NoError(t, err)
	require.Equal(t, 400, resp.StatusCode)
}
