package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration_FindCurrentPriceListOfRoom_Success(t *testing.T) {
	username := "host_p_01"
	jwt, _, room := createUserAndRoomForPrice(username)

	// Create 2, so we can check if the second one overrides the first one.
	createRoomPrice(jwt, room)
	priceList2 := createRoomPrice(jwt, room)

	resp, err := FindCurrentPriceListOfRoom(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	priceListGot := ResponseToRoomPrice(resp)
	require.Equal(t, priceList2.ID, priceListGot.ID)
}

func TestIntegration_FindCurrentPriceListOfRoom_NotFound(t *testing.T) {
	username := "host_p_02"
	_, _, room := createUserAndRoomForPrice(username)

	resp, err := FindCurrentPriceListOfRoom(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_FindPriceListsByRoomId_Success(t *testing.T) {
	username := "host_p_03"
	jwt, _, room := createUserAndRoomForPrice(username)

	createRoomPrice(jwt, room)
	createRoomPrice(jwt, room)

	resp, err := FindPriceListsByRoomId(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listsGot := ResponseToRoomPriceLists(resp)
	require.Equal(t, 2, len(listsGot))
}

func TestIntegration_FindPriceListsByRoomId_NotFound(t *testing.T) {
	resp, err := FindPriceListsByRoomId(999888777)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_FindPriceListById_Success(t *testing.T) {
	username := "host_p_05"
	jwt, _, room := createUserAndRoomForPrice(username)

	li := createRoomPrice(jwt, room)

	resp, err := FindPriceListById(li.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listGot := ResponseToRoomPrice(resp)
	require.Equal(t, li.ID, listGot.ID)
}

func TestIntegration_FindPriceListById_NotFound(t *testing.T) {
	username := "host_p_06"
	_, _, room := createUserAndRoomForPrice(username)

	resp, err := FindPriceListById(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
