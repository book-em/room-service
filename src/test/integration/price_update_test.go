package integration

import (
	"bookem-room-service/internal"
	test "bookem-room-service/test/unit"
	"bookem-room-service/util"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegration_UpdatePriceList_Success(t *testing.T) {
	username := "host_pu_01"
	registerUser(username, "1234", util.Host)
	jwt := loginUser2(username, "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)

	dto := test.DefaultRoomCreateDTO
	dto.HostID = jwtObj.ID
	resp, _ := createRoom(jwt, dto)
	room := responseToRoom(resp)

	// [Phase 1] Create initial price list
	{
		createPriceDto := internal.CreateRoomPriceListDTO{
			RoomID: room.ID,
			Items: []internal.CreateRoomPriceItemDTO{
				{
					ExistingID: 0,
					DateFrom:   time.Date(2025, 8, 14, 0, 0, 0, 0, time.UTC),
					DateTo:     time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
					Price:      100,
				},
				{
					ExistingID: 0,
					DateFrom:   time.Date(2025, 11, 14, 0, 0, 0, 0, time.UTC),
					DateTo:     time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC),
					Price:      200,
				},
			},
		}

		resp, err := createRoomPrice(jwt, createPriceDto)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		priceList := responseToRoomPrice(resp)
		require.Equal(t, 2, len(priceList.Items))
		require.Equal(t, room.ID, priceList.RoomID)
	}

	// [Phase 2] Update price list: reuse one item, add a new one
	{
		resp, _ := findCurrentPriceListOfRoom(room.ID)
		currentPriceList := responseToRoomPrice(resp)
		require.Equal(t, 2, len(currentPriceList.Items))

		// Remove item [0], keep item [1], add new item
		itemToAdd := internal.CreateRoomPriceItemDTO{
			ExistingID: 0,
			DateFrom:   time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
			DateTo:     time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
			Price:      150,
		}
		itemToKeep := internal.CreateRoomPriceItemDTO{
			ExistingID: currentPriceList.Items[1].ID,
			DateFrom:   time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			DateTo:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			Price:      999, // ignored
		}

		updatePriceDto := internal.CreateRoomPriceListDTO{
			RoomID: room.ID,
			Items:  []internal.CreateRoomPriceItemDTO{itemToAdd, itemToKeep},
		}

		resp, err := createRoomPrice(jwt, updatePriceDto)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		priceList := responseToRoomPrice(resp)
		require.Equal(t, 2, len(priceList.Items))
		require.Equal(t, room.ID, priceList.RoomID)

		// Final check: confirm reuse and replacement
		{
			resp, _ := findCurrentPriceListOfRoom(room.ID)
			newPriceList := responseToRoomPrice(resp)
			require.Equal(t, 2, len(newPriceList.Items))

			foundIDs := []uint{
				newPriceList.Items[0].ID,
				newPriceList.Items[1].ID,
			}

			require.NotContains(t, foundIDs, currentPriceList.Items[0].ID) // Removed item
			require.Contains(t, foundIDs, currentPriceList.Items[1].ID)    // Reused item
		}
	}
}
