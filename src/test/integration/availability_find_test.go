package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_FindCurrentAvailabilityListOfRoom_Success(t *testing.T) {
	username := "host_b_01"
	jwt, _, room := createUserAndRoom(username)

	// Create 2, so we can check if the second one overrides the first one.
	createRoomAvailabilityList(jwt, room)
	availList2 := createRoomAvailabilityList(jwt, room)

	resp, err := findCurrentAvailabilityListOfRoom(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	availListGot := responseToRoomAvailability(resp)
	require.Equal(t, availList2.ID, availListGot.ID)
}

func TestIntegration_FindCurrentAvailabilityListOfRoom_NotFound(t *testing.T) {
	username := "host_b_02"
	_, _, room := createUserAndRoom(username)

	resp, err := findCurrentAvailabilityListOfRoom(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_FindAvailabilityListsByRoomId_Success(t *testing.T) {
	username := "host_b_03"
	jwt, _, room := createUserAndRoom(username)

	createRoomAvailabilityList(jwt, room)
	createRoomAvailabilityList(jwt, room)

	resp, err := findAvailabilityListsByRoomId(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listsGot := responseToRoomAvailabilityLists(resp)
	require.Equal(t, 2, len(listsGot))
}

func TestIntegration_FindAvailabilityListsByRoomId_NotFound(t *testing.T) {
	resp, err := findAvailabilityListsByRoomId(999888777)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_FindAvailabilityListById_Success(t *testing.T) {
	username := "host_b_05"
	jwt, _, room := createUserAndRoom(username)

	li := createRoomAvailabilityList(jwt, room)

	resp, err := findAvailabilityListById(li.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listGot := responseToRoomAvailability(resp)
	require.Equal(t, li.ID, listGot.ID)
}

func TestIntegration_FindAvailabilityListById_NotFound(t *testing.T) {
	cleanup("room")
	cleanup("user")

	username := "host_b_06"
	_, _, room := createUserAndRoom(username)

	resp, err := findAvailabilityListById(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
