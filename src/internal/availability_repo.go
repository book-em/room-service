package internal

import (
	"gorm.io/gorm"
)

type RoomAvailabilityRepo interface {
	CreateList(availList *RoomAvailabilityList) error
	FindListById(id uint) (*RoomAvailabilityList, error)
	FindListsByRoomId(roomId uint) ([]RoomAvailabilityList, error)
	FindCurrentListOfRoom(roomId uint) (*RoomAvailabilityList, error)
}
type roomAvailabilityRepo struct{ db *gorm.DB }

func NewRoomAvailabilityRepo(db *gorm.DB) RoomAvailabilityRepo { return &roomAvailabilityRepo{db} }

func (r *roomAvailabilityRepo) CreateList(availList *RoomAvailabilityList) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(availList).Error; err != nil {
			return err
		}

		if err := tx.Model(&Room{}).
			Where("id = ?", availList.RoomID).
			Update("availability_list_id", availList.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *roomAvailabilityRepo) FindListById(id uint) (*RoomAvailabilityList, error) {
	var li RoomAvailabilityList
	err := r.db.Preload("Items").Where("id = ?", id).First(&li).Error
	if err != nil {
		return nil, err
	}
	return &li, nil
}

func (r *roomAvailabilityRepo) FindListsByRoomId(roomId uint) ([]RoomAvailabilityList, error) {
	var lists []RoomAvailabilityList
	err := r.db.Preload("Items").Where("room_id = ?", roomId).Find(&lists).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *roomAvailabilityRepo) FindCurrentListOfRoom(roomId uint) (*RoomAvailabilityList, error) {
	var latest RoomAvailabilityList
	err := r.db.
		Preload("Items").
		Where("room_id = ?", roomId).
		Order("effective_from DESC").
		First(&latest).Error

	if err != nil {
		return nil, err
	}
	return &latest, nil
}
