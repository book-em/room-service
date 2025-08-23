package integration

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/internal"
	test "bookem-room-service/test/unit"
	"bookem-room-service/util"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIntegration_UpdateAvailability_Success(t *testing.T) {
	username := "host_a_01"
	RegisterUser(username, "1234", userclient.Host)
	jwt := LoginUser2(username, "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)

	dto := test.DefaultRoomCreateDTO
	dto.HostID = jwtObj.ID
	resp, _ := CreateRoom(jwt, dto)
	room := ResponseToRoom(resp)

	// This test has two phases:
	//
	// [Phase 1] create the initial availability list.
	//
	{
		createRoomDto := internal.CreateRoomAvailabilityListDTO{
			RoomID: room.ID,
			Items: []internal.CreateRoomAvailabilityItemDTO{
				internal.CreateRoomAvailabilityItemDTO{
					ExistingID: 0,
					DateFrom:   time.Date(2025, 8, 14, 0, 0, 0, 0, time.UTC),
					DateTo:     time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
					Available:  true,
				}, internal.CreateRoomAvailabilityItemDTO{
					ExistingID: 0,
					DateFrom:   time.Date(2025, 11, 14, 0, 0, 0, 0, time.UTC),
					DateTo:     time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC),
					Available:  true,
				},
			},
		}

		resp, err := CreateRoomAvailability(jwt, createRoomDto)
		if err != nil {
			panic(err)
		}
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		availList := ResponseToRoomAvailability(resp)

		require.Equal(t, 2, len(availList.Items))
		require.Equal(t, room.ID, availList.RoomID)
	}
	//
	// [Phase 2] do it again, but we add a new item and remove an existing item.
	// The goal of this phase is to confirm that the one item that wasn't removed
	// is reused in the DB.
	//
	{
		resp, _ := FindCurrentAvailabilityListOfRoom(room.ID)
		currentAvailabilityList := ResponseToRoomAvailability(resp)
		require.Equal(t, 2, len(currentAvailabilityList.Items))

		// We will remove item [0] and keep item [1].

		itemToAdd := internal.CreateRoomAvailabilityItemDTO{
			ExistingID: 0,
			DateFrom:   time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC),
			DateTo:     time.Date(2025, 9, 15, 0, 0, 0, 0, time.UTC),
			Available:  true,
		}
		itemToKeep := internal.CreateRoomAvailabilityItemDTO{
			ExistingID: currentAvailabilityList.Items[1].ID,
			//
			// These fields don't matter if we're reusing existing items.
			// So you can change them and it won't matter.
			// But validation still happens.
			//
			DateFrom:  time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			DateTo:    time.Date(2025, 3, 2, 0, 0, 0, 0, time.UTC),
			Available: true,
		}

		createRoomDto := internal.CreateRoomAvailabilityListDTO{
			RoomID: room.ID,
			Items: []internal.CreateRoomAvailabilityItemDTO{
				itemToAdd,
				itemToKeep,
			},
		}

		resp, err := CreateRoomAvailability(jwt, createRoomDto)
		if err != nil {
			panic(err)
		}
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		availList := ResponseToRoomAvailability(resp)

		require.Equal(t, 2, len(availList.Items))
		require.Equal(t, room.ID, availList.RoomID)

		// After it passes, we check if everything is OK

		{
			resp, _ := FindCurrentAvailabilityListOfRoom(room.ID)
			newAvailabilityList := ResponseToRoomAvailability(resp)
			require.Equal(t, 2, len(currentAvailabilityList.Items))

			foundIDs := []uint{
				newAvailabilityList.Items[0].ID,
				newAvailabilityList.Items[1].ID,
			}

			require.NotContains(t, foundIDs, currentAvailabilityList.Items[0].ID) // Removed item
			require.Contains(t, foundIDs, currentAvailabilityList.Items[1].ID)    // New item
		}
	}
}
