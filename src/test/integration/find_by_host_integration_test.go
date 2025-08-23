package integration

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/internal"
	test "bookem-room-service/test/unit"
	"bookem-room-service/util"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_FindByHost_Success(t *testing.T) {
	registerUser("user3", "1234", userclient.Host)
	jwt := loginUser2("user3", "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)
	hostId := jwtObj.ID

	roomCreateDTO := test.DefaultRoomCreateDTO
	roomCreateDTO.HostID = hostId

	resp, _ := createRoom(jwt, roomCreateDTO)
	room := responseToRoom(resp)
	resp, _ = createRoom(jwt, roomCreateDTO)
	room2 := responseToRoom(resp)
	roomsExpect := []internal.RoomDTO{room, room2}

	resp, err := findRoomsByHostId(hostId)

	require.NoError(t, err)
	roomsGot := responseToRooms(resp)
	require.Equal(t, roomsExpect, roomsGot)
}

func TestIntegration_FindByHost_MissingId(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%shost/", url_room))

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIntegration_FindByHost_HostNotFound(t *testing.T) {
	resp, err := findRoomsByHostId(888888)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
