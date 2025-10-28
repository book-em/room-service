package test

import (
	"bookem-room-service/internal"
	"bookem-room-service/util"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ClearYear_Success(t *testing.T) {
	svc, _, _, _, _ := CreateTestRoomService()

	date1 := time.Now()
	date2 := time.Now().Add(24 * time.Hour)

	date1Res, date2Res := svc.ClearYear(context.Background(), date1, date2)

	assert.Equal(t, date1Res.Day(), date1.Day())
	assert.Equal(t, date1Res.Month(), date1.Month())
	assert.NotEqual(t, date1Res.Year(), date1.Year())
	assert.Equal(t, date2Res.Day(), date2.Day())
	assert.Equal(t, date2Res.Month(), date2.Month())
	assert.NotEqual(t, date2Res.Year(), date2.Year())
}

func Test_CalculatePriceForOneDay_BasePrice_Success(t *testing.T) {
	svc, _, _, _, _ := CreateTestRoomService()
	day := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	day = util.ClearYear(day)
	basePrice := uint(1000)
	rulePrice := uint(100)
	guests := uint(4)

	var rules internal.RoomPriceList = internal.RoomPriceList{
		ID:            1,
		RoomID:        1,
		EffectiveFrom: time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		BasePrice:     basePrice,
		PerGuest:      true,
		Items:         []internal.RoomPriceItem{},
	}
	rule := internal.RoomPriceItem{
		ID:       1,
		DateFrom: time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Price:    rulePrice,
	}
	rules.Items = append(rules.Items, rule)

	priceRes := svc.CalculatePriceForOneDay(context.Background(), day, guests, rules)

	assert.Equal(t, float32(guests*basePrice), priceRes)
}

func Test_CalculatePriceForOneDay_DefinedPriceByRule_Success(t *testing.T) {
	svc, _, _, _, _ := CreateTestRoomService()
	day := time.Date(2025, 8, 13, 0, 0, 0, 0, time.UTC)
	day = util.ClearYear(day)
	basePrice := uint(1000)
	rulePrice := uint(100)

	var rules internal.RoomPriceList = internal.RoomPriceList{
		ID:            1,
		RoomID:        1,
		EffectiveFrom: time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		BasePrice:     basePrice,
		PerGuest:      false,
		Items:         []internal.RoomPriceItem{},
	}
	rule := internal.RoomPriceItem{
		ID:       1,
		DateFrom: time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Price:    rulePrice,
	}
	rules.Items = append(rules.Items, rule)

	priceRes := svc.CalculatePriceForOneDay(context.Background(), day, uint(4), rules)

	assert.Equal(t, float32(rulePrice), priceRes)
}

func Test_CalculatePrice_UndefinedRules_Fail(t *testing.T) {
	svc, _, _, mockPriceRepo, _ := CreateTestRoomService()
	dateFrom := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 8, 25, 0, 0, 0, 0, time.UTC)
	guestsNumber := uint(2)
	roomId := uint(1)

	mockPriceRepo.On("FindCurrentListOfRoom", roomId).Return(nil, fmt.Errorf("not found"))
	totalPrice, perGuest, err := svc.CalculatePrice(context.Background(), dateFrom, dateTo, guestsNumber, roomId)

	assert.Error(t, err)
	assert.Equal(t, float32(0), totalPrice)
	assert.Equal(t, false, perGuest)
	mockPriceRepo.AssertNumberOfCalls(t, "FindCurrentListOfRoom", 1)
	mockPriceRepo.AssertExpectations(t)
}

func Test_CalculatePrice_PerGuest_Success(t *testing.T) {
	svc, _, _, mockPriceRepo, _ := CreateTestRoomService()
	dateFrom := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 8, 24, 0, 0, 0, 0, time.UTC)
	guestsNumber := uint(2)
	roomId := uint(1)

	var rules internal.RoomPriceList = internal.RoomPriceList{
		ID:            1,
		RoomID:        roomId,
		EffectiveFrom: dateFrom,
		BasePrice:     300,
		PerGuest:      true,
		Items:         []internal.RoomPriceItem{},
	}
	rule1 := internal.RoomPriceItem{
		ID:       1,
		DateFrom: time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 14, 0, 0, 0, 0, time.UTC),
		Price:    100,
	}
	rule2 := internal.RoomPriceItem{
		ID:       2,
		DateFrom: time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 24, 0, 0, 0, 0, time.UTC),
		Price:    200,
	}
	rules.Items = append(rules.Items, rule1, rule2)

	mockPriceRepo.On("FindCurrentListOfRoom", roomId).Return(&rules, nil)
	totalPrice, perGuest, err := svc.CalculatePrice(context.Background(), dateFrom, dateTo, guestsNumber, roomId)

	assert.NoError(t, err)
	// 5 x 100 x 2  +  5 x 300 x 2  +  5 x 200 x 2  =  6000
	assert.Equal(t, float32(6000), totalPrice)
	assert.Equal(t, true, perGuest)
	mockPriceRepo.AssertNumberOfCalls(t, "FindCurrentListOfRoom", 1)
	mockPriceRepo.AssertExpectations(t)
}

func Test_CalculatePrice_FlatPrice_Success(t *testing.T) {
	svc, _, _, mockPriceRepo, _ := CreateTestRoomService()
	dateFrom := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 8, 24, 0, 0, 0, 0, time.UTC)
	guestsNumber := uint(2)
	roomId := uint(1)

	var rules internal.RoomPriceList = internal.RoomPriceList{
		ID:            1,
		RoomID:        roomId,
		EffectiveFrom: dateFrom,
		BasePrice:     300,
		PerGuest:      false,
		Items:         []internal.RoomPriceItem{},
	}
	rule1 := internal.RoomPriceItem{
		ID:       1,
		DateFrom: time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 14, 0, 0, 0, 0, time.UTC),
		Price:    100,
	}
	rule2 := internal.RoomPriceItem{
		ID:       2,
		DateFrom: time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 24, 0, 0, 0, 0, time.UTC),
		Price:    200,
	}
	rules.Items = append(rules.Items, rule1, rule2)

	mockPriceRepo.On("FindCurrentListOfRoom", roomId).Return(&rules, nil)
	totalPrice, perGuest, err := svc.CalculatePrice(context.Background(), dateFrom, dateTo, guestsNumber, roomId)

	assert.NoError(t, err)
	// 5 x 100  +  5 x 300  +  5 x 200  =  3000
	assert.Equal(t, float32(3000), totalPrice)
	assert.Equal(t, false, perGuest)
	mockPriceRepo.AssertNumberOfCalls(t, "FindCurrentListOfRoom", 1)
	mockPriceRepo.AssertExpectations(t)
}

func Test_IsRoomAvailableForOneDay_OverlappingAvailable_Success(t *testing.T) {
	// One interval overlaps with another, choose the least one
	svc, _, _, _, _ := CreateTestRoomService()
	day := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	day = util.ClearYear(day)

	rule1 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		Available: false,
	}
	rule2 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		Available: true,
	}

	rules := []internal.RoomAvailabilityItem{}
	rules = append(rules, rule1, rule2)

	isAvailable := svc.IsRoomAvailableForOneDay(context.Background(), day, rules)

	assert.Equal(t, true, isAvailable)
}

func Test_IsRoomAvailableForOneDay_OverlappingUnavailable_Success(t *testing.T) {
	// One interval overlaps with another, choose the least one
	svc, _, _, _, _ := CreateTestRoomService()
	day := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	day = util.ClearYear(day)

	rule1 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		Available: true,
	}
	rule2 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		Available: false,
	}

	rules := []internal.RoomAvailabilityItem{}
	rules = append(rules, rule1, rule2)

	isAvailable := svc.IsRoomAvailableForOneDay(context.Background(), day, rules)

	assert.Equal(t, false, isAvailable)
}

func Test_IsRoomAvailableForOneDay_NoOverlappingUnavailable_Success(t *testing.T) {
	// For undefined day, the room is unavailable by default
	svc, _, _, _, _ := CreateTestRoomService()
	day := time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC)
	day = util.ClearYear(day)

	rule1 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		Available: true,
	}
	rule2 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		Available: false,
	}

	rules := []internal.RoomAvailabilityItem{}
	rules = append(rules, rule1, rule2)

	isAvailable := svc.IsRoomAvailableForOneDay(context.Background(), day, rules)

	assert.Equal(t, false, isAvailable)
}

func Test_IsRoomAvailable_UndefinedRulesUnavailable(t *testing.T) {
	svc, _, mockRepo, _, _ := CreateTestRoomService()
	dateFrom := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)
	roomId := uint(1)

	mockRepo.On("FindCurrentListOfRoom", roomId).Return(nil, fmt.Errorf("room availability list not found"))

	canBook := svc.IsRoomAvailable(context.Background(), dateFrom, dateTo, roomId)

	assert.Equal(t, false, canBook)
	mockRepo.AssertNumberOfCalls(t, "FindCurrentListOfRoom", 1)
	mockRepo.AssertExpectations(t)
}

func Test_IsRoomAvailable_OverlappingAvailable(t *testing.T) {
	svc, _, mockRepo, _, _ := CreateTestRoomService()
	dateFrom := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC)
	roomId := uint(1)

	var rules internal.RoomAvailabilityList = internal.RoomAvailabilityList{
		ID:     1,
		RoomID: roomId,
		Items:  []internal.RoomAvailabilityItem{},
	}
	rule1 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC),
		Available: false,
	}
	rule2 := internal.RoomAvailabilityItem{
		ID:        2,
		DateFrom:  time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		Available: true,
	}
	rules.Items = append(rules.Items, rule1, rule2)

	mockRepo.On("FindCurrentListOfRoom", roomId).Return(&rules, nil)

	canBook := svc.IsRoomAvailable(context.Background(), dateFrom, dateTo, roomId)

	assert.Equal(t, true, canBook)
	mockRepo.AssertNumberOfCalls(t, "FindCurrentListOfRoom", 1)
	mockRepo.AssertExpectations(t)
}

func Test_IsRoomAvailable_OverlappingAvailableAdvanced(t *testing.T) {
	svc, _, mockRepo, _, _ := CreateTestRoomService()
	dateFrom := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)
	roomId := uint(1)

	var rules internal.RoomAvailabilityList = internal.RoomAvailabilityList{
		ID:     1,
		RoomID: roomId,
		Items:  []internal.RoomAvailabilityItem{},
	}
	rule1 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC),
		Available: false,
	}
	rule2 := internal.RoomAvailabilityItem{
		ID:        2,
		DateFrom:  time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Available: true,
	}
	rule3 := internal.RoomAvailabilityItem{
		ID:        3,
		DateFrom:  time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		Available: true,
	}
	rules.Items = append(rules.Items, rule1, rule2, rule3)

	mockRepo.On("FindCurrentListOfRoom", roomId).Return(&rules, nil)

	canBook := svc.IsRoomAvailable(context.Background(), dateFrom, dateTo, roomId)

	assert.Equal(t, true, canBook)
	mockRepo.AssertNumberOfCalls(t, "FindCurrentListOfRoom", 1)
	mockRepo.AssertExpectations(t)
}

func Test_IsRoomAvailable_OverlappingUnavailable(t *testing.T) {
	svc, _, mockRepo, _, _ := CreateTestRoomService()
	dateFrom := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 8, 18, 0, 0, 0, 0, time.UTC)
	roomId := uint(1)

	var rules internal.RoomAvailabilityList = internal.RoomAvailabilityList{
		ID:     1,
		RoomID: roomId,
		Items:  []internal.RoomAvailabilityItem{},
	}
	rule1 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC),
		Available: true,
	}
	rule2 := internal.RoomAvailabilityItem{
		ID:        2,
		DateFrom:  time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		Available: false,
	}
	rules.Items = append(rules.Items, rule1, rule2)

	mockRepo.On("FindCurrentListOfRoom", roomId).Return(&rules, nil)

	canBook := svc.IsRoomAvailable(context.Background(), dateFrom, dateTo, roomId)

	assert.Equal(t, false, canBook)
	mockRepo.AssertNumberOfCalls(t, "FindCurrentListOfRoom", 1)
	mockRepo.AssertExpectations(t)
}

func Test_CalculateUnitPricePerGuest(t *testing.T) {
	svc, _, _, _, _ := CreateTestRoomService()
	perGuest := true
	guestsNumber := uint(2)
	dateFrom := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC)
	totalPrice := float32(10000)

	unitPrice := svc.CalculateUnitPrice(context.Background(), perGuest, guestsNumber, dateFrom, dateTo, totalPrice)

	// 10000 / 10 / 2  =  500
	assert.Equal(t, float32(500), unitPrice)
}

func Test_CalculateUnitPriceFlat(t *testing.T) {
	svc, _, _, _, _ := CreateTestRoomService()
	perGuest := false
	guestsNumber := uint(99)
	dateFrom := time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 8, 19, 0, 0, 0, 0, time.UTC)
	totalPrice := float32(10000)

	unitPrice := svc.CalculateUnitPrice(context.Background(), perGuest, guestsNumber, dateFrom, dateTo, totalPrice)

	// 10000 / 10  =  1000
	assert.Equal(t, float32(1000), unitPrice)
}

func Test_PreparePaginatedResult_Success(t *testing.T) {
	svc, _, _, _, _ := CreateTestRoomService()

	pageNumber := uint(1)
	pageSize := uint(10)
	hits := []internal.RoomResultDTO{}
	for i := 0; i < 50; i++ {
		hits = append(hits, *DefaultRoomResult)
	}

	hitsResult, resultInfo := svc.PreparePaginatedResult(context.Background(), hits, pageNumber, pageSize)

	assert.Equal(t, pageNumber, resultInfo.Page)
	assert.Equal(t, pageSize, resultInfo.PageSize)
	assert.Equal(t, uint(5), resultInfo.TotalPages)
	assert.Equal(t, uint(50), resultInfo.TotalHits)
	assert.Equal(t, int(10), len(hitsResult))
}

func Test_PreparePaginatedResult_OutOfMargin(t *testing.T) {
	// Show last page result if page number exceeds total
	svc, _, _, _, _ := CreateTestRoomService()

	pageNumber := uint(99999)
	pageSize := uint(10)
	hits := []internal.RoomResultDTO{}
	for i := 0; i < 15; i++ {
		hits = append(hits, *DefaultRoomResult)
	}

	hitsResult, resultInfo := svc.PreparePaginatedResult(context.Background(), hits, pageNumber, pageSize)

	assert.Equal(t, pageNumber, resultInfo.Page)
	assert.Equal(t, pageSize, resultInfo.PageSize)
	assert.Equal(t, pageSize, resultInfo.PageSize)
	assert.Equal(t, uint(2), resultInfo.TotalPages)
	assert.Equal(t, uint(15), resultInfo.TotalHits)
	assert.Equal(t, int(5), len(hitsResult))
}

func Test_PreparePaginatedResult_LastPage(t *testing.T) {
	// Case when the last page is selected
	svc, _, _, _, _ := CreateTestRoomService()

	pageNumber := uint(2)
	pageSize := uint(10)
	hits := []internal.RoomResultDTO{}
	for i := 0; i < 15; i++ {
		hits = append(hits, *DefaultRoomResult)
	}

	hitsResult, resultInfo := svc.PreparePaginatedResult(context.Background(), hits, pageNumber, pageSize)

	assert.Equal(t, pageNumber, resultInfo.Page)
	assert.Equal(t, pageSize, resultInfo.PageSize)
	assert.Equal(t, pageSize, resultInfo.PageSize)
	assert.Equal(t, uint(2), resultInfo.TotalPages)
	assert.Equal(t, uint(15), resultInfo.TotalHits)
	assert.Equal(t, int(5), len(hitsResult))
}

// ---------------------------------------------------- OLD

func Test_FindAvailableRooms_InvalidDate(t *testing.T) {
	svc, mockRepo, _, _, _ := CreateTestRoomService()

	d := *DefaultRoomsQueryDTO
	d.DateFrom = d.DateTo.Add(24 * time.Hour)

	roomsGot, infoGot, err := svc.FindAvailableRooms(context.Background(), d)

	assert.Nil(t, roomsGot)
	assert.Nil(t, infoGot)
	assert.Error(t, err)
	mockRepo.AssertNumberOfCalls(t, "FindByFilters", 0)
	mockRepo.AssertExpectations(t)
}

func Test_FindAvailableRooms_DbError(t *testing.T) {
	svc, mockRepo, _, _, _ := CreateTestRoomService()

	d := *DefaultRoomsQueryDTO

	mockRepo.On("FindByFilters", d.GuestsNumber, d.Address).Return(nil, fmt.Errorf("db error"))

	roomsGot, infoGot, err := svc.FindAvailableRooms(context.Background(), d)

	assert.Nil(t, roomsGot)
	assert.Nil(t, infoGot)
	assert.Error(t, err)
	mockRepo.AssertNumberOfCalls(t, "FindByFilters", 1)
	mockRepo.AssertExpectations(t)
}

func Test_FindAvailableRooms_Success(t *testing.T) {
	svc, mockRepo, mockAvailRepo, mockPriceRepo, _ := CreateTestRoomService()

	rooms := []internal.Room{}
	room1 := internal.Room{
		ID:        1,
		HostID:    1,
		Name:      "room1",
		Address:   "address1",
		MinGuests: 2,
		MaxGuests: 5,
	}
	room2 := internal.Room{
		ID:        2,
		HostID:    2,
		Name:      "room2",
		Address:   "address2",
		MinGuests: 1,
		MaxGuests: 4,
	}
	rooms = append(rooms, room1, room2)

	// room availability rules
	var availRules1 internal.RoomAvailabilityList = internal.RoomAvailabilityList{
		ID:     1,
		RoomID: room1.ID,
		Items:  []internal.RoomAvailabilityItem{},
	}
	rule1 := internal.RoomAvailabilityItem{
		ID:        1,
		DateFrom:  time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC),
		Available: false,
	}
	rule2 := internal.RoomAvailabilityItem{
		ID:        2,
		DateFrom:  time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		Available: true,
	}
	availRules1.Items = append(availRules1.Items, rule1, rule2)

	var availRules2 internal.RoomAvailabilityList = internal.RoomAvailabilityList{
		ID:     2,
		RoomID: room2.ID,
		Items:  []internal.RoomAvailabilityItem{},
	}
	rule1 = internal.RoomAvailabilityItem{
		ID:        3,
		DateFrom:  time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 5, 0, 0, 0, 0, time.UTC),
		Available: false,
	}
	rule2 = internal.RoomAvailabilityItem{
		ID:        4,
		DateFrom:  time.Date(2025, 8, 6, 0, 0, 0, 0, time.UTC),
		DateTo:    time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC),
		Available: true,
	}
	availRules2.Items = append(availRules2.Items, rule1, rule2)

	// price rules
	var priceRules1 internal.RoomPriceList = internal.RoomPriceList{
		ID:            1,
		RoomID:        room1.ID,
		EffectiveFrom: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		BasePrice:     300,
		PerGuest:      true,
		Items:         []internal.RoomPriceItem{},
	}
	priceRule1 := internal.RoomPriceItem{
		ID:       1,
		DateFrom: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Price:    100,
	}
	priceRule2 := internal.RoomPriceItem{
		ID:       2,
		DateFrom: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC),
		Price:    200,
	}
	priceRules1.Items = append(priceRules1.Items, priceRule1, priceRule2)

	var priceRules2 internal.RoomPriceList = internal.RoomPriceList{
		ID:            2,
		RoomID:        room2.ID,
		EffectiveFrom: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		BasePrice:     500,
		PerGuest:      false,
		Items:         []internal.RoomPriceItem{},
	}
	priceRule1 = internal.RoomPriceItem{
		ID:       3,
		DateFrom: time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC),
		Price:    400,
	}
	priceRule2 = internal.RoomPriceItem{
		ID:       4,
		DateFrom: time.Date(2025, 8, 16, 0, 0, 0, 0, time.UTC),
		DateTo:   time.Date(2025, 8, 30, 0, 0, 0, 0, time.UTC),
		Price:    600,
	}
	priceRules2.Items = append(priceRules2.Items, priceRule1, priceRule2)

	query := internal.RoomsQueryDTO{
		Address:      "address",
		GuestsNumber: 3,
		DateFrom:     time.Date(2025, 8, 10, 0, 0, 0, 0, time.UTC),
		DateTo:       time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC),
		PageNumber:   1,
		PageSize:     1,
	}

	// [1] none address
	query.Address = "none"
	mockRepo.On("FindByFilters", query.GuestsNumber, query.Address).Return(nil, nil)

	roomsGot, infoGot, err := svc.FindAvailableRooms(context.Background(), query)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(roomsGot))
	assert.Equal(t, query.PageNumber, infoGot.Page)
	assert.Equal(t, query.PageSize, infoGot.PageSize)
	assert.Equal(t, uint(0), infoGot.TotalHits)
	assert.Equal(t, uint(0), infoGot.TotalPages)
	mockRepo.AssertNumberOfCalls(t, "FindByFilters", 1)
	mockRepo.AssertExpectations(t)

	// [2] testing pagination
	query.PageNumber = uint(5)
	query.PageSize = uint(2)
	query.Address = "none"
	mockRepo.On("FindByFilters", query.GuestsNumber, query.Address).Return(nil, nil)

	roomsGot, infoGot, err = svc.FindAvailableRooms(context.Background(), query)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(roomsGot))
	assert.Equal(t, query.PageNumber, infoGot.Page)
	assert.Equal(t, query.PageSize, infoGot.PageSize)
	assert.Equal(t, uint(0), infoGot.TotalHits)
	assert.Equal(t, uint(0), infoGot.TotalPages)
	mockRepo.AssertNumberOfCalls(t, "FindByFilters", 2)
	mockRepo.AssertExpectations(t)

	// [3] Both rooms are available
	query.PageSize = 1
	query.PageNumber = 1
	query.Address = "address"
	mockRepo.On("FindByFilters", query.GuestsNumber, query.Address).Return(rooms, nil)
	mockAvailRepo.On("FindCurrentListOfRoom", room1.ID).Return(&availRules1, nil)
	mockAvailRepo.On("FindCurrentListOfRoom", room2.ID).Return(&availRules2, nil)
	mockPriceRepo.On("FindCurrentListOfRoom", room1.ID).Return(&priceRules1, nil)
	mockPriceRepo.On("FindCurrentListOfRoom", room2.ID).Return(&priceRules2, nil)

	roomsGot, infoGot, err = svc.FindAvailableRooms(context.Background(), query)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(roomsGot))
	assert.Equal(t, uint(2), infoGot.TotalHits)
	assert.Equal(t, query.PageSize, infoGot.PageSize)
	assert.Equal(t, uint(2), infoGot.TotalPages)
	mockRepo.AssertNumberOfCalls(t, "FindByFilters", 3)
	mockRepo.AssertExpectations(t)

	// [4] Single room is available
	query.PageSize = 4
	query.PageNumber = 1
	query.DateFrom = time.Date(2025, 8, 6, 0, 0, 0, 0, time.UTC)
	query.DateTo = time.Date(2025, 8, 7, 0, 0, 0, 0, time.UTC)
	mockRepo.On("FindByFilters", query.GuestsNumber, query.Address).Return(rooms, nil)
	mockAvailRepo.On("FindCurrentListOfRoom", room1.ID).Return(&availRules1, nil)
	mockAvailRepo.On("FindCurrentListOfRoom", room2.ID).Return(&availRules2, nil)
	mockPriceRepo.On("FindCurrentListOfRoom", room1.ID).Return(&priceRules1, nil)
	mockPriceRepo.On("FindCurrentListOfRoom", room2.ID).Return(&priceRules2, nil)

	roomsGot, infoGot, err = svc.FindAvailableRooms(context.Background(), query)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(roomsGot))
	assert.Equal(t, uint(1), infoGot.TotalHits)
	assert.Equal(t, query.PageSize, infoGot.PageSize)
	assert.Equal(t, uint(1), infoGot.TotalPages)
	mockRepo.AssertNumberOfCalls(t, "FindByFilters", 4)
	mockRepo.AssertExpectations(t)
}

func Test_ExcludeDeletedRooms_Success(t *testing.T) {
	svc, _, _, _, _ := CreateTestRoomService()

	rooms := []internal.Room{}
	room1 := internal.Room{
		ID:        1,
		HostID:    1,
		Name:      "room1",
		Address:   "address1",
		MinGuests: 2,
		MaxGuests: 5,
		Deleted:   false,
	}
	room2 := internal.Room{
		ID:        2,
		HostID:    2,
		Name:      "room2",
		Address:   "address2",
		MinGuests: 1,
		MaxGuests: 4,
		Deleted:   true,
	}
	rooms = append(rooms, room1, room2)

	roomsGot := svc.ExcludeDeletedRooms(context.Background(), rooms)

	assert.Equal(t, 1, len(roomsGot))
	assert.Equal(t, room1, roomsGot[0])
}
