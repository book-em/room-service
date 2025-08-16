package internal

import (
	"time"
)

type Room struct {
	ID          uint     `gorm:"primaryKey"`
	HostID      uint     `gorm:"not null"`
	Name        string   `gorm:"type:varchar(50);not null"`
	Description string   ``
	Address     string   `gorm:"type:varchar(150);not null"`
	MinGuests   uint     `gorm:"not null"`
	MaxGuests   uint     `gorm:"not null"`
	Photos      []string `gorm:"type:text;serializer:json"`
	Commodities []string `gorm:"type:text;serializer:json"`

	// AvailabilityListID refers to the latest list of times when the room is available.
	// If there is no availability list, then this is `nil`.
	AvailabilityListID *uint

	// PriceListID refers to the latest list of prices of the room.
	// If there is no price list, then this is `nil`.
	PriceListID *uint
}

// RoomAvailabilityList is a list of dates when a specific room is available for booking.
type RoomAvailabilityList struct {
	ID            uint                   `gorm:"primaryKey"`
	RoomID        uint                   `gorm:"not null;index"`
	Room          Room                   ``
	EffectiveFrom time.Time              `gorm:"not null"`
	Items         []RoomAvailabilityItem `gorm:"many2many:room_availability_list_items;"`
}

// RoomAvailabilityItem defines a date range when a room is available (or not available).
type RoomAvailabilityItem struct {
	ID       uint                   `gorm:"primaryKey"`
	Lists    []RoomAvailabilityList `gorm:"many2many:room_availability_list_items;"`
	DateFrom time.Time              `gorm:"not null"`
	DateTo   time.Time              `gorm:"not null"`

	// Available determines if this item works as a "union" or a "disjoint".
	// When true, the owning room availability list is expanded by this date
	// range (normal behavior) When false, the owning room availablity list is
	// shrunk by this date range.
	//
	// Effectively, this allows you to define a room availability list like so:
	//
	// [Jan 1, Dec 31, true], [Jan 1, Jan 7, false]
	//
	// Which means that the room is available for booking on all days except
	// from Jan 1st to Jan 7th.
	Available bool
}

type RoomPriceList struct {
	ID            uint            `gorm:"primaryKey"`
	RoomID        uint            `gorm:"not null;index"`
	Room          Room            ``
	EffectiveFrom time.Time       `gorm:"not null"`
	BasePrice     uint            `gorm:"not null"`
	Items         []RoomPriceItem `gorm:"many2many:room_price_list_items;"`
	// If PerGuest is true, then this price is defined per guest of the room. If
	// false, then this price is defined regardless of the number of guests a
	// reservation is made for.
	PerGuest bool `gorm:"not null"`
}

type RoomPriceItem struct {
	ID       uint            `gorm:"primaryKey"`
	Lists    []RoomPriceList `gorm:"many2many:room_price_list_items;"`
	DateFrom time.Time       `gorm:"not null"`
	DateTo   time.Time       `gorm:"not null"`
	Price    uint            `gorm:"not null"`
}
