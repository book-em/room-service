package internal

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
