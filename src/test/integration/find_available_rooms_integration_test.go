package integration

import (
	test "bookem-room-service/test/unit"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegration_FindAvailableRooms_Success(t *testing.T) {
	cleanup("room")
	cleanup("user")
	setupRooms(5)

	query := test.DefaultRoomsQueryDTO
	query.Address = "Room Address"
	query.DateFrom = time.Date(2025, 8, 22, 0, 0, 0, 0, time.UTC)
	query.DateTo = time.Date(2025, 8, 23, 0, 0, 0, 0, time.UTC)

	resp, err := findAvailableRooms(*query)
	result := responseToFindAvailableRooms(resp)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 5, int(result.Info.TotalHits))
	require.Equal(t, 5, len(result.Hits))
	require.Equal(t, query.PageNumber, result.Info.Page)
	require.Equal(t, query.PageSize, result.Info.PageSize)
	require.Equal(t, 1, int(result.Info.TotalPages))
}

func TestIntegration_FindAvailableRooms_Pagination_Success(t *testing.T) {
	cleanup("room")
	cleanup("user")
	setupRooms(5)

	query := test.DefaultRoomsQueryDTO
	query.Address = "  Room ADDRESS  "
	query.DateFrom = time.Date(2025, 8, 22, 0, 0, 0, 0, time.UTC)
	query.DateTo = time.Date(2025, 8, 23, 0, 0, 0, 0, time.UTC)
	query.PageNumber = 3
	query.PageSize = 2

	resp, err := findAvailableRooms(*query)
	result := responseToFindAvailableRooms(resp)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 5, int(result.Info.TotalHits))
	require.Equal(t, 1, len(result.Hits))
	require.Equal(t, query.PageNumber, result.Info.Page)
	require.Equal(t, query.PageSize, result.Info.PageSize)
	require.Equal(t, 3, int(result.Info.TotalPages))

	query.PageNumber = 1
	query.PageSize = 2

	resp, err = findAvailableRooms(*query)
	result = responseToFindAvailableRooms(resp)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 5, int(result.Info.TotalHits))
	require.Equal(t, 2, len(result.Hits))
	require.Equal(t, query.PageNumber, result.Info.Page)
	require.Equal(t, query.PageSize, result.Info.PageSize)
	require.Equal(t, 3, int(result.Info.TotalPages))
}

func TestIntegration_FindAvailableRooms_NoHits_Success(t *testing.T) {
	cleanup("room")
	cleanup("user")
	setupRooms(5)

	// [1] No available rooms for a certain date range
	query := test.DefaultRoomsQueryDTO
	query.Address = "Room Address"
	query.DateFrom = time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC)
	query.DateTo = time.Date(2025, 10, 23, 0, 0, 0, 0, time.UTC)

	resp, err := findAvailableRooms(*query)
	result := responseToFindAvailableRooms(resp)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 0, int(result.Info.TotalHits))
	require.Equal(t, 0, len(result.Hits))
	require.Equal(t, query.PageNumber, result.Info.Page)
	require.Equal(t, query.PageSize, result.Info.PageSize)
	require.Equal(t, 0, int(result.Info.TotalPages))

	// [2] No available rooms for a specific address
	query = test.DefaultRoomsQueryDTO
	query.Address = "unknown address"

	resp, err = findAvailableRooms(*query)
	result = responseToFindAvailableRooms(resp)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 0, int(result.Info.TotalHits))
	require.Equal(t, 0, len(result.Hits))
	require.Equal(t, query.PageNumber, result.Info.Page)
	require.Equal(t, query.PageSize, result.Info.PageSize)
	require.Equal(t, 0, int(result.Info.TotalPages))

	// [3] No available rooms for the given number of guests
	query = test.DefaultRoomsQueryDTO
	query.Address = "Room Address"
	query.GuestsNumber = 9999

	resp, err = findAvailableRooms(*query)
	result = responseToFindAvailableRooms(resp)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 0, int(result.Info.TotalHits))
	require.Equal(t, 0, len(result.Hits))
	require.Equal(t, query.PageNumber, result.Info.Page)
	require.Equal(t, query.PageSize, result.Info.PageSize)
	require.Equal(t, 0, int(result.Info.TotalPages))
}

func TestIntegration_FindAvailableRooms_InvalidDate(t *testing.T) {
	cleanup("room")
	cleanup("user")
	setupRooms(5)

	query := test.DefaultRoomsQueryDTO
	query.Address = "Room Address"
	query.DateFrom = time.Date(2025, 10, 23, 0, 0, 0, 0, time.UTC)
	query.DateTo = time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC)

	resp, err := findAvailableRooms(*query)

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
