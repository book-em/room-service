package reservationclient

import (
	utils "bookem-room-service/util"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ReservationClient interface {
	// GetActiveHostReservations finds all reservations made by `host` that
	// haven't completed yet. The user must be a host.
	GetActiveHostReservations(ctx context.Context, jwt string, roomIDs []uint) ([]ReservationDTO, error)
}

type reservationClient struct {
	baseURL string
}

func NewReservationClient() ReservationClient {
	return &reservationClient{
		baseURL: "http://reservation-service:8080/api", // TODO: This should not be hardcoded
	}
}

func (c *reservationClient) GetActiveHostReservations(ctx context.Context, jwt string, roomIDs []uint) ([]ReservationDTO, error) {
	utils.TEL.Push(ctx, "get-active-reservations-for-host")
	defer utils.TEL.Pop()

	dto := RoomIDsDTO{IDs: roomIDs}
	jsonBytes, err := json.Marshal(dto)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/reservations/host/active", c.baseURL), bytes.NewBuffer(jsonBytes))
	if err != nil {
		utils.TEL.Error("preparing request error ", err)
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		utils.TEL.Error("request error ", err)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.TEL.Error("parsing response error", err)
		return nil, err
	}

	var obj []ReservationDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		utils.TEL.Error("JSON unmarshall error", err)
		return nil, err
	}

	return obj, nil
}
