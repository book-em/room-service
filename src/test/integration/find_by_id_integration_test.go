package integration

import (
	test "bookem-room-service/test/unit"
	"bookem-room-service/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_FindById_Success(t *testing.T) {
	cleanup("room")
	cleanup("user")

	registerUser("user2", "1234", util.Host)
	jwt := loginUser2("user2", "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)

	roomCreateDTO := test.DefaultRoomCreateDTO
	roomCreateDTO.HostID = jwtObj.ID

	resp, _ := createRoom(jwt, roomCreateDTO)
	room := responseToRoom(resp)

	roomId := room.ID

	resp, err := findRoomById(roomId)

	require.NoError(t, err)
	roomGot := responseToRoom(resp)

	require.Equal(t, roomId, roomGot.ID)
	require.Equal(t, room, roomGot)
}

func TestIntegration_FindById_MissingId(t *testing.T) {
	resp, err := http.Get(url_room)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_FindById_Deleted(t *testing.T) {
	cleanup("room")
	cleanup("user")

	registerUser("user2", "1234", util.Host)
	jwt := loginUser2("user2", "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)

	roomCreateDTO := test.DefaultRoomCreateDTO
	roomCreateDTO.HostID = jwtObj.ID
	roomCreateDTO.Deleted = true

	resp, _ := createRoom(jwt, roomCreateDTO)
	room := responseToRoom(resp)

	roomId := room.ID

	resp, err := findRoomById(roomId)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
