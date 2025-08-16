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
	// ExistingID is either the ID of an RoomAvailabiltyItem that already
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
