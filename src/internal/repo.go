package internal

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(room *Room) error
	Update(room *Room) error
	Delete(room *Room) error
	FindById(id uint) (*Room, error)
	FindByHost(hostId uint) ([]Room, error)
	FindByFilters(guestsNumber uint, address string) ([]Room, error)
	DeleteRoomsByHostId(hostId uint) error
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

func (r *repository) DeleteRoomsByHostId(hostId uint) error {
	err := r.db.Model(&Room{}).Where("host_id = ?", hostId).Update("deleted", true).Error
	if err != nil {
		return err
	}
	return nil
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

func (r *repository) FindByFilters(guestsNumber uint, address string) ([]Room, error) {
	var rooms []Room
	query := r.db.Where("min_guests <= ? and max_guests >= ?", guestsNumber, guestsNumber)

	if address != "" {
		query = query.Where("TRIM(LOWER(address)) LIKE CONCAT('%' || TRIM(LOWER( ? )) || '%')", address)
	}

	err := query.Find(&rooms).Error
	if err != nil {
		return nil, query.Error
	}

	return rooms, nil
}
