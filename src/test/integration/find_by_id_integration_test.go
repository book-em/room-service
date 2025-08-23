package integration

import (
	"bookem-room-service/client/userclient"
	test "bookem-room-service/test/unit"
	"bookem-room-service/util"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_FindById_Success(t *testing.T) {
	cleanup("room")
	cleanup("user")

	RegisterUser("user2", "1234", userclient.Host)
	jwt := LoginUser2("user2", "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)

	roomCreateDTO := test.DefaultRoomCreateDTO
	roomCreateDTO.HostID = jwtObj.ID

	resp, _ := CreateRoom(jwt, roomCreateDTO)
	room := ResponseToRoom(resp)

	roomId := room.ID

	resp, err := FindRoomById(roomId)

	require.NoError(t, err)
	roomGot := ResponseToRoom(resp)

	require.Equal(t, roomId, roomGot.ID)
	require.Equal(t, room, roomGot)
}

func TestIntegration_FindById_MissingId(t *testing.T) {
	resp, err := http.Get(URL_room)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
