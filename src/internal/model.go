package internal

type Room struct {
	ID          uint     `json:"id"       gorm:"primaryKey"`
	HostID      uint     `json:"hostID"`
	Name        string   `json:"name"     gorm:"type:varchar(50);not null"`
	Description string   `json:"description"`
	Address     string   `json:"address"  gorm:"type:varchar(150);not null"`
	MinGuests   uint     `json:"minGuests"`
	MaxGuests   uint     `json:"maxGuests"`
	Photos      []string `json:"photos"   gorm:"type:text;serializer:json"`
	Commodities []string `json:"commodities" gorm:"type:text;serializer:json"`
}
