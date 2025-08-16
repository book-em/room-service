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
	// If there is not availability list item, then this is `nil`.
	AvailabilityListID *uint
}

// RoomAvailabilityList is a list of dates when a specific room is available for booking.
type RoomAvailabilityList struct {
	ID            uint                   `gorm:"primaryKey"`
	RoomID        uint                   `gorm:"not null;index"`
	Room          Room                   ``
	EffectiveFrom time.Time              `gorm:"not null"`
	Items         []RoomAvailabilityItem `gorm:"foreignKey:AvailabilityListID"`
}

// RoomAvailabilityItem defines a date range when a room is available (or not available).
type RoomAvailabilityItem struct {
	ID                 uint      `gorm:"primaryKey"`
	AvailabilityListID uint      `gorm:"not null;index"`
	DateFrom           time.Time `gorm:"not null"`
	DateTo             time.Time `gorm:"not null"`

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
