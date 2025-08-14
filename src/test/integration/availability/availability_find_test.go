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

func createUserAndRoom(username string) (string, *util.Jwt, internal.RoomDTO) {
	RegisterUser(username, "1234", userclient.Host)
	jwt := LoginUser2(username, "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)

	dto := test.DefaultRoomCreateDTO
	dto.HostID = jwtObj.ID
	resp, _ := CreateRoom(jwt, dto)
	room := ResponseToRoom(resp)

	return jwt, jwtObj, room
}

func createRoomAvailability(jwt string, room internal.RoomDTO) internal.RoomAvailabilityListDTO {
	dto := internal.CreateRoomAvailabilityListDTO{
		RoomID: room.ID,
		Items:  test.DefaultCreateAvailabilityListDTO.Items,
	}

	resp, err := CreateRoomAvailability(jwt, dto)
	if err != nil {
		panic(err)
	}

	return ResponseToRoomAvailability(resp)
}

func TestIntegration_FindCurrentAvailabilityListOfRoom_Success(t *testing.T) {
	username := "host_b_01"
	jwt, _, room := createUserAndRoom(username)

	// Create 2, so we can check if the second one overrides the first one.
	createRoomAvailability(jwt, room)
	availList2 := createRoomAvailability(jwt, room)

	resp, err := FindCurrentAvailabilityListOfRoom(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	availListGot := ResponseToRoomAvailability(resp)
	require.Equal(t, availList2.ID, availListGot.ID)
}

func TestIntegration_FindCurrentAvailabilityListOfRoom_NotFound(t *testing.T) {
	username := "host_b_02"
	_, _, room := createUserAndRoom(username)

	resp, err := FindCurrentAvailabilityListOfRoom(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_FindAvailabilityListsByRoomId_Success(t *testing.T) {
	username := "host_b_03"
	jwt, _, room := createUserAndRoom(username)

	createRoomAvailability(jwt, room)
	createRoomAvailability(jwt, room)

	resp, err := FindAvailabilityListsByRoomId(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listsGot := ResponseToRoomAvailabilityLists(resp)
	require.Equal(t, 2, len(listsGot))
}

func TestIntegration_FindAvailabilityListsByRoomId_NotFound(t *testing.T) {
	resp, err := FindAvailabilityListsByRoomId(999888777)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_FindAvailabilityListById_Success(t *testing.T) {
	username := "host_b_05"
	jwt, _, room := createUserAndRoom(username)

	li := createRoomAvailability(jwt, room)

	resp, err := FindAvailabilityListById(li.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listGot := ResponseToRoomAvailability(resp)
	require.Equal(t, li.ID, listGot.ID)
}

func TestIntegration_FindAvailabilityListById_NotFound(t *testing.T) {
	username := "host_b_06"
	_, _, room := createUserAndRoom(username)

	resp, err := FindAvailabilityListById(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
