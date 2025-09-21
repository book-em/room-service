package internal

import "time"

type RoomDTO struct {
	ID          uint     `json:"id"`
	HostID      uint     `json:"hostID"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	MinGuests   uint     `json:"minGuests"`
	MaxGuests   uint     `json:"maxGuests"`
	Photos      []string `json:"photos"`
	Commodities []string `json:"commodities"`
}

type CreateRoomDTO struct {
	HostID        uint     `json:"hostID"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Address       string   `json:"address"`
	MinGuests     uint     `json:"minGuests"`
	MaxGuests     uint     `json:"maxGuests"`
	PhotosPayload []string `json:"photosPayload"`
	Commodities   []string `json:"commodities"`
}

func NewRoomDTO(r *Room) RoomDTO {
	return RoomDTO{
		ID:          r.ID,
		HostID:      r.HostID,
		Name:        r.Name,
		Description: r.Description,
		Address:     r.Address,
		MinGuests:   r.MinGuests,
		MaxGuests:   r.MaxGuests,
		Photos:      r.Photos,
		Commodities: r.Commodities,
	}
}

// ---------------------------------------------------------------

type CreateRoomAvailabilityListDTO struct {
	RoomID uint                            `json:"roomId"`
	Items  []CreateRoomAvailabilityItemDTO `json:"items"`
}

type RoomAvailabilityListDTO struct {
	ID            uint                      `json:"id"`
	RoomID        uint                      `json:"roomId"`
	EffectiveFrom time.Time                 `json:"effectiveFrom"`
	Items         []RoomAvailabilityItemDTO `json:"items"`
}

func NewRoomAvailabilityListDTO(list *RoomAvailabilityList) RoomAvailabilityListDTO {
	items := make([]RoomAvailabilityItemDTO, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, NewRoomAvailabilityItemDTO(item))
	}

	return RoomAvailabilityListDTO{
		ID:            list.ID,
		RoomID:        list.RoomID,
		EffectiveFrom: list.EffectiveFrom,
		Items:         items,
	}
}

type RoomAvailabilityItemDTO struct {
	ID        uint      `json:"id"`
	DateFrom  time.Time `json:"dateFrom"`
	DateTo    time.Time `json:"dateTo"`
	Available bool      `json:"available"`
}

type CreateRoomAvailabilityItemDTO struct {
	// ExistingID is either the ID of an RoomAvailabilityItem that already
	// exists, or 0 if this is a new item. When 0, a new one will be created in
	// the DB. When not 0, it will reuse the existing object.
	ExistingID uint      `json:"existingId"`
	DateFrom   time.Time `json:"dateFrom"`
	DateTo     time.Time `json:"dateTo"`
	Available  bool      `json:"available"`
}

func NewRoomAvailabilityItemDTO(item RoomAvailabilityItem) RoomAvailabilityItemDTO {
	return RoomAvailabilityItemDTO{
		ID:        item.ID,
		DateFrom:  item.DateFrom,
		DateTo:    item.DateTo,
		Available: item.Available,
	}
}

// ---------------------------------------------------------------

type CreateRoomPriceListDTO struct {
	RoomID    uint                     `json:"roomId"`
	Items     []CreateRoomPriceItemDTO `json:"items"`
	BasePrice uint                     `json:"basePrice"`
	PerGuest  bool                     `json:"perGuest"`
}

type RoomPriceListDTO struct {
	ID            uint               `json:"id"`
	RoomID        uint               `json:"roomId"`
	EffectiveFrom time.Time          `json:"effectiveFrom"`
	BasePrice     uint               `json:"basePrice"`
	Items         []RoomPriceItemDTO `json:"items"`
	PerGuest      bool               `json:"perGuest"`
}

func NewRoomPriceListDTO(list *RoomPriceList) RoomPriceListDTO {
	items := make([]RoomPriceItemDTO, 0, len(list.Items))
	for _, item := range list.Items {
		items = append(items, NewRoomPriceItemDTO(item))
	}

	return RoomPriceListDTO{
		ID:            list.ID,
		RoomID:        list.RoomID,
		EffectiveFrom: list.EffectiveFrom,
		BasePrice:     list.BasePrice,
		Items:         items,
		PerGuest:      list.PerGuest,
	}
}

type RoomPriceItemDTO struct {
	ID       uint      `json:"id"`
	DateFrom time.Time `json:"dateFrom"`
	DateTo   time.Time `json:"dateTo"`
	Price    uint      `json:"price"`
}

type CreateRoomPriceItemDTO struct {
	// ExistingID is either the ID of an RoomPriceItem that already
	// exists, or 0 if this is a new item. When 0, a new one will be created in
	// the DB. When not 0, it will reuse the existing object.
	ExistingID uint      `json:"existingId"`
	DateFrom   time.Time `json:"dateFrom"`
	DateTo     time.Time `json:"dateTo"`
	Price      uint      `json:"price"`
}

func NewRoomPriceItemDTO(item RoomPriceItem) RoomPriceItemDTO {
	return RoomPriceItemDTO{
		ID:       item.ID,
		DateFrom: item.DateFrom,
		DateTo:   item.DateTo,
		Price:    item.Price,
	}
}

type RoomsQueryDTO struct {
	Address      string    `form:"address"`
	GuestsNumber uint      `form:"guestsNumber" binding:"required,min=1"`
	DateFrom     time.Time `form:"dateFrom" binding:"required"`
	DateTo       time.Time `form:"dateTo" binding:"required"`
	PageNumber   uint      `form:"pageNumber" binding:"required,min=1"`
	PageSize     uint      `form:"pageSize" binding:"required,min=1"`
}

type PaginatedResultInfoDTO struct {
	Page       uint `json:"page"`
	PageSize   uint `json:"pageSize"`
	TotalPages uint `json:"totalPages"`
	TotalHits  uint `json:"totalHits"`
}

type RoomResultDTO struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Address     string   `json:"address"`
	Photos      []string `json:"photos" gorm:"type:text;serializer:json"`
	PerGuest    bool     `json:"perGuest"`
	UnitPrice   float32  `json:"unitPrice"`
	TotalPrice  float32  `json:"totalPrice"`
}

func NewRoomResultDTO(room Room, perGuest bool, unitPrice float32, totalPrice float32) RoomResultDTO {
	return RoomResultDTO{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Address:     room.Address,
		Photos:      room.Photos,
		PerGuest:    perGuest,
		UnitPrice:   unitPrice,
		TotalPrice:  totalPrice,
	}
}

type RoomsResultDTO struct {
	Hits []RoomResultDTO        `json:"hits"`
	Info PaginatedResultInfoDTO `json:"info"`
}

func NewRoomsResultDTO(hits []RoomResultDTO, info PaginatedResultInfoDTO) RoomsResultDTO {
	return RoomsResultDTO{
		Hits: hits,
		Info: info,
	}
}

// --------------------------------------------------------

type RoomReservationQueryDTO struct {
	RoomID     uint      `json:"roomId"`
	DateFrom   time.Time `json:"dateFrom"`
	DateTo     time.Time `json:"dateTo"`
	GuestCount uint      `json:"guestCount"`
}

type RoomReservationQueryResponseDTO struct {
	Available bool `json:"available"`
	TotalCost uint `json:"totalCost"`
}
