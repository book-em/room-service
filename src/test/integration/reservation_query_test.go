package integration

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/internal"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegration_QueryForReservation_Success(t *testing.T) {
	// [Step 1] Create a host and a room
	hostUsername := "host_query_test"
	hostJwt, _, room := createUserAndRoom(hostUsername)

	// [Step 2] Set availability and pricing for the room
	createRoomAvailabilityList(hostJwt, room)

	createRoomPriceList(hostJwt, room)

	// [Step 3] Register a guest user
	guestUsername := "guest_query_test"
	registerUser(guestUsername, "1234", userclient.Guest)
	guestJwt := loginUser2(guestUsername, "1234")

	// [Step 4] Prepare reservation query DTO
	dto := internal.RoomReservationQueryDTO{
		RoomID:     room.ID,
		DateFrom:   time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		DateTo:     time.Date(2025, 8, 21, 0, 0, 0, 0, time.UTC),
		GuestCount: 2,
	}

	jsonBytes, err := json.Marshal(dto)
	require.NoError(t, err)

	// [Step 5] Make the request
	req, err := http.NewRequest(http.MethodPost, url_room+"reservation/query", bytes.NewBuffer(jsonBytes))
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+guestJwt)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// [Step 6] Validate response
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result internal.RoomReservationQueryResponseDTO
	err = json.Unmarshal(bodyBytes, &result)
	require.NoError(t, err)

	require.True(t, result.Available)
	require.Equal(t, uint(200), result.TotalCost) // 2 days x 100 (flat rate)
}
