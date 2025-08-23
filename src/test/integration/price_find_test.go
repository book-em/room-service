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
	createRoomPriceList(jwt, room)
	priceList2 := createRoomPriceList(jwt, room)

	resp, err := findCurrentPriceListOfRoom(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	priceListGot := responseToRoomPrice(resp)
	require.Equal(t, priceList2.ID, priceListGot.ID)
}

func TestIntegration_FindCurrentPriceListOfRoom_NotFound(t *testing.T) {
	username := "host_p_02"
	_, _, room := createUserAndRoomForPrice(username)

	resp, err := findCurrentPriceListOfRoom(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_FindPriceListsByRoomId_Success(t *testing.T) {
	username := "host_p_03"
	jwt, _, room := createUserAndRoomForPrice(username)

	createRoomPriceList(jwt, room)
	createRoomPriceList(jwt, room)

	resp, err := findPriceListsByRoomId(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listsGot := responseToRoomPriceLists(resp)
	require.Equal(t, 2, len(listsGot))
}

func TestIntegration_FindPriceListsByRoomId_NotFound(t *testing.T) {
	resp, err := findPriceListsByRoomId(999888777)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_FindPriceListById_Success(t *testing.T) {
	username := "host_p_05"
	jwt, _, room := createUserAndRoomForPrice(username)

	li := createRoomPriceList(jwt, room)

	resp, err := findPriceListById(li.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	listGot := responseToRoomPrice(resp)
	require.Equal(t, li.ID, listGot.ID)
}

func TestIntegration_FindPriceListById_NotFound(t *testing.T) {
	username := "host_p_06"
	_, _, room := createUserAndRoomForPrice(username)

	resp, err := findPriceListById(room.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}
