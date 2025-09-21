package integration

import (
	"bookem-room-service/client/userclient"
	"bookem-room-service/internal"
	test "bookem-room-service/test/unit"
	"bookem-room-service/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const url_user = "http://user-service:8080/api/"
const url_room = "http://room-service:8080/api/"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genName(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func registerUser(username_or_email string, password string, role userclient.UserRole) (*http.Response, error) {
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
		Name:     genName(6),
		Surname:  genName(6),
		Address:  genName(10),
	}

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url_user+"register", "application/json", bytes.NewBuffer(jsonBytes))
	return resp, err
}

func loginUser(username_or_email string, password string) (*http.Response, error) {
	dto := userclient.LoginDTO{
		UsernameOrEmail: username_or_email,
		Password:        password,
	}

	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url_user+"login", "application/json", bytes.NewBuffer(jsonBytes))
	return resp, err
}

func loginUser2(username_or_email string, password string) string {
	resp, _ := loginUser(username_or_email, password)

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

func createRoom(jwt string, dto internal.CreateRoomDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url_room+"new", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func findRoomById(id uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%s%d", url_room, id)) // No forward slash between them, it's in `URL`
	return resp, err
}

func findAvailableRooms(dto internal.RoomsQueryDTO) (*http.Response, error) {
	params := url.Values{}
	val, _ := dto.DateFrom.UTC().MarshalText()
	params.Add("dateFrom", string(val))
	val, _ = dto.DateTo.UTC().MarshalText()
	params.Add("dateTo", string(val))
	params.Add("address", dto.Address)
	params.Add("guestsNumber", fmt.Sprintf("%d", dto.GuestsNumber))
	params.Add("pageNumber", fmt.Sprintf("%d", dto.PageNumber))
	params.Add("pageSize", fmt.Sprintf("%d", dto.PageSize))

	req, err := http.NewRequest(http.MethodGet, url_room+"all?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func findRoomsByHostId(hostId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%shost/%d", url_room, hostId)) // No forward slash between them, it's in `URL`
	return resp, err
}

func responseToRoom(resp *http.Response) internal.RoomDTO {
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

func responseToRooms(resp *http.Response) []internal.RoomDTO {
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

func createRoomAvailability(jwt string, dto internal.CreateRoomAvailabilityListDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url_room+"available", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func responseToRoomAvailability(resp *http.Response) internal.RoomAvailabilityListDTO {
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

func responseToFindAvailableRooms(resp *http.Response) internal.RoomsResultDTO {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("failed to read response body: %v", err))
	}

	var obj internal.RoomsResultDTO
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		fmt.Print(string(bodyBytes))
		panic(fmt.Sprintf("failed to unmarshal: %v", err))
	}

	return obj
}

func findCurrentAvailabilityListOfRoom(roomId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%savailable/room/%d", url_room, roomId))
	return resp, err
}

func findAvailabilityListsByRoomId(roomId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%savailable/room/all/%d", url_room, roomId))
	return resp, err
}

func findAvailabilityListById(id uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%savailable/%d", url_room, id))
	return resp, err
}

func responseToRoomAvailabilityLists(resp *http.Response) []internal.RoomAvailabilityListDTO {
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

func createRoomPrice(jwt string, dto internal.CreateRoomPriceListDTO) (*http.Response, error) {
	jsonBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url_room+"price", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+jwt)
	return http.DefaultClient.Do(req)
}

func responseToRoomPrice(resp *http.Response) internal.RoomPriceListDTO {
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
func responseToRoomPriceLists(resp *http.Response) []internal.RoomPriceListDTO {
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
func findCurrentPriceListOfRoom(roomId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%sprice/room/%d", url_room, roomId))
	return resp, err
}
func findPriceListsByRoomId(roomId uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%sprice/room/all/%d", url_room, roomId))
	return resp, err
}
func findPriceListById(id uint) (*http.Response, error) {
	resp, err := http.Get(fmt.Sprintf("%sprice/%d", url_room, id))
	return resp, err
}

func createUserAndRoom(username string) (string, *util.Jwt, internal.RoomDTO) {
	registerUser(username, "1234", userclient.Host)
	jwt := loginUser2(username, "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)

	dto := test.DefaultRoomCreateDTO
	dto.HostID = jwtObj.ID
	resp, _ := createRoom(jwt, dto)
	room := responseToRoom(resp)

	return jwt, jwtObj, room
}

func createRoomAvailabilityList(jwt string, room internal.RoomDTO) internal.RoomAvailabilityListDTO {
	dto := internal.CreateRoomAvailabilityListDTO{
		RoomID: room.ID,
		Items:  test.DefaultCreateAvailabilityListDTO.Items,
	}

	resp, err := createRoomAvailability(jwt, dto)
	if err != nil {
		panic(err)
	}

	return responseToRoomAvailability(resp)
}

func createUserAndRoomForPrice(username string) (string, *util.Jwt, internal.RoomDTO) {
	registerUser(username, "1234", userclient.Host)
	jwt := loginUser2(username, "1234")
	jwtObj, _ := util.GetJwtFromString(jwt)

	dto := test.DefaultRoomCreateDTO
	dto.HostID = jwtObj.ID
	resp, _ := createRoom(jwt, dto)
	room := responseToRoom(resp)

	return jwt, jwtObj, room
}

func createRoomPriceList(jwt string, room internal.RoomDTO) internal.RoomPriceListDTO {
	dto := internal.CreateRoomPriceListDTO{
		RoomID: room.ID,
		Items:  test.DefaultCreatePriceListDTO.Items,
	}

	resp, err := createRoomPrice(jwt, dto)
	if err != nil {
		panic(err)
	}

	return responseToRoomPrice(resp)
}

func setupRooms(quantity int) {
	for i := range quantity {
		username := fmt.Sprintf("host%d", i)
		jwt, _, room := createUserAndRoom(username)
		createRoomAvailabilityList(jwt, room)
		createRoomPriceList(jwt, room)
	}
}
