package test

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/internal"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
)

const URL_user = "http://user-service:8080/api/"
const URL_room = "http://room-service:8080/api/"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenName(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RegisterUser(username_or_email string, password string, role userclient.UserRole) (*http.Response, error) {
	username := username_or_email
	email := username + "@gmail.com"

	if strings.HasSuffix(username_or_email, "@gmail.com") {
		username = strings.Split(username_or_email, "@")[0]
		email = username_or_email
	}

	dto := userclient.UserCreateDTO{
		Username: username,
		Password: password,
		Email:    email,
		Role:     string(role),
		Name:     GenName(6),
		Surname:  GenName(6),
		Address:  GenName(10),
	}

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(URL_user+"register", "application/json", bytes.NewBuffer(jsonBytes))
	return resp, err
}

func LoginUser(username_or_email string, password string) (*http.Response, error) {
	dto := userclient.LoginDTO{
		UsernameOrEmail: username_or_email,
		Password:        password,
	}

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(URL_user+"login", "application/json", bytes.NewBuffer(jsonBytes))
	return resp, err
}

func LoginUser2(username_or_email string, password string) string {
	resp, _ := LoginUser(username_or_email, password)

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var token userclient.JWTDTO
	if err := json.Unmarshal(bodyBytes, &token); err != nil {
		panic(fmt.Sprintf("failed to unmarshal jwt: %v", err))
	}

	return token.Jwt
}

func CreateRoom(jwt string, dto internal.CreateRoomDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL_room+"new", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func FindRoomById(id uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%s%d", URL_room, id)) // No forward slash between them, it's in `URL`
	return resp, err
}

func FindRoomsByHostId(hostId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%shost/%d", URL_room, hostId)) // No forward slash between them, it's in `URL`
	return resp, err
}

func ResponseToRoom(resp *http.Response) internal.RoomDTO {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var obj internal.RoomDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		fmt.Print(string(bodyBytes))
		panic(fmt.Sprintf("failed to unmarshal: %v", err))
	}

	return obj
}

func ResponseToRooms(resp *http.Response) []internal.RoomDTO {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var obj []internal.RoomDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		panic(fmt.Sprintf("failed to unmarshal: %v", err))
	}

	return obj
}

func CreateRoomAvailability(jwt string, dto internal.CreateRoomAvailabilityListDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL_room+"available", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func ResponseToRoomAvailability(resp *http.Response) internal.RoomAvailabilityListDTO {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var obj internal.RoomAvailabilityListDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		panic(fmt.Sprintf("failed to unmarshal: %v", err))
	}

	return obj
}

func FindCurrentAvailabilityListOfRoom(roomId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%savailable/room/%d", URL_room, roomId))
	return resp, err
}

func FindAvailabilityListsByRoomId(roomId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%savailable/room/all/%d", URL_room, roomId))
	return resp, err
}

func FindAvailabilityListById(id uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%savailable/%d", URL_room, id))
	return resp, err
}

func ResponseToRoomAvailabilityLists(resp *http.Response) []internal.RoomAvailabilityListDTO {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var obj []internal.RoomAvailabilityListDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		panic(fmt.Sprintf("failed to unmarshal: %v", err))
	}

	return obj
}

func CreateRoomPrice(jwt string, dto internal.CreateRoomPriceListDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, URL_room+"price", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func ResponseToRoomPrice(resp *http.Response) internal.RoomPriceListDTO {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var obj internal.RoomPriceListDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		panic(fmt.Sprintf("failed to unmarshal: %v", err))
	}

	return obj
}
func ResponseToRoomPriceLists(resp *http.Response) []internal.RoomPriceListDTO {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var obj []internal.RoomPriceListDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		panic(fmt.Sprintf("failed to unmarshal: %v", err))
	}

	return obj
}
func FindCurrentPriceListOfRoom(roomId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%sprice/room/%d", URL_room, roomId))
	return resp, err
}
func FindPriceListsByRoomId(roomId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%sprice/room/all/%d", URL_room, roomId))
	return resp, err
}
func FindPriceListById(id uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%sprice/%d", URL_room, id))
	return resp, err
}
