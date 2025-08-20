package internal

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(room *Room) error
	Update(room *Room) error
	Delete(room *Room) error
	FindById(id uint) (*Room, error)
	FindByHost(hostId uint) ([]Room, error)
	FindAvailableRooms(location string, guestsNumber uint, dateFrom time.Time, dateTo time.Time, pageNumber uint, pageSize uint) ([]RoomResultDTO, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(room *Room) error {
	return r.db.Create(room).Error
}

func (r *repository) Update(room *Room) error {
	return r.db.Save(room).Error
}

func (r *repository) Delete(room *Room) error {
	return r.db.Delete(&Room{}, room.ID).Error
}

func (r *repository) FindById(id uint) (*Room, error) {
	var room Room
	err := r.db.Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *repository) FindByHost(hostId uint) ([]Room, error) {
	var rooms []Room
	err := r.db.Where("host_id = ?", hostId).Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *repository) FindAvailableRooms(location string, guestsNumber uint, dateFrom time.Time, dateTo time.Time, pageNumber uint, pageSize uint) (rooms []RoomResultDTO, totalHits int64, err error) {

	with := `
		WITH temp_table (room_id, available, interval) AS
		(
			SELECT DISTINCT ON
				(rooms.id) rooms.id,
				room_availability_items.available,
				(CAST(@dateTo AS DATE) - CAST(@dateFrom AS DATE)) AS interval
			FROM rooms
			INNER JOIN room_availability_list_items 
				ON rooms.availability_list_id = room_availability_list_items.room_availability_list_id
			INNER JOIN room_availability_items 
				ON room_availability_items.id = room_availability_list_items.room_availability_item_id
			WHERE room_availability_items.date_from <= @dateFrom 
				AND room_availability_items.date_to >= @dateTo
				AND min_guests <= @guestsNumber
				AND max_guests >= @guestsNumber
				AND (CASE WHEN @address = '' THEN LOWER(address) ELSE LOWER(@address) end) 
					LIKE CONCAT('%' || LOWER(address) || '%')
			ORDER BY rooms.id, room_availability_items.date_to - room_availability_items.date_from ASC
		)
	`

	count := `SELECT COUNT(*) FROM temp_table`
	err = r.db.Raw(with+count, map[string]interface{}{
		"address":      location,
		"dateFrom":     dateFrom,
		"dateTo":       dateTo,
		"guestsNumber": guestsNumber,
	}).Count(&totalHits).Error

	if err != nil {
		return nil, 0, err
	}

	find := `SELECT 
			rooms.id AS ID, 
			rooms.name AS Name, 
			rooms.description AS Description, 
			rooms.address AS Address,  
			rooms.photos AS Photos, 
			room_price_lists.base_price AS BasePrice,
			CASE
				WHEN per_guest = TRUE THEN temp_table.interval * room_price_lists.base_price * @guestsNumber
				ELSE temp_table.interval * room_price_lists.base_price
			END AS TotalPrice
		FROM temp_table, rooms, room_price_lists
		WHERE temp_table.room_id = rooms.id and temp_table.available = TRUE 
			AND rooms.price_list_id = room_price_lists.id
		LIMIT @limit OFFSET @offset
	`
	offset := int((pageNumber - 1) * pageSize)
	println("offset", offset, "page size", int(pageSize))
	err = r.db.Raw(with+find, map[string]interface{}{
		"address":      location,
		"dateFrom":     dateFrom,
		"dateTo":       dateTo,
		"guestsNumber": guestsNumber,
		"limit":        int(pageSize),
		"offset":       offset,
	}).Scan(&rooms).Error

	return rooms, totalHits, nil
}
