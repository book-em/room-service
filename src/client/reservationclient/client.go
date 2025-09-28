package reservationclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ReservationClient interface {
	// GetActiveHostReservations finds all reservations made by `host` that
	// haven't completed yet. The user must be a host.
	GetActiveHostReservations(jwt string, roomIDs []uint) ([]ReservationDTO, error)
}

type reservationClient struct {
	baseURL string
}

func NewReservationClient() ReservationClient {
	return &reservationClient{
		baseURL: "http://reservation-service:8080/api", // TODO: This should not be hardcoded
	}
}

func (c *reservationClient) GetActiveHostReservations(jwt string, roomIDs []uint) ([]ReservationDTO, error) {
	log.Printf("Get active host reservations")

	dto := RoomIDsDTO{IDs: roomIDs}
	jsonBytes, err := json.Marshal(dto)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/reservations/host/active", c.baseURL), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Error %v", err)
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Parsing response error: %v", err)
		return nil, err
	}

	var obj []ReservationDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		log.Printf("JSON Unmarshall error: %v", err)
		return nil, err
	}

	return obj, nil
}
