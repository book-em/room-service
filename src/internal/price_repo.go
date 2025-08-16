package internal

import (
	"gorm.io/gorm"
)

type RoomPriceRepo interface {
	CreateList(priceList *RoomPriceList) error
	FindListById(id uint) (*RoomPriceList, error)
	FindListsByRoomId(roomId uint) ([]RoomPriceList, error)
	FindCurrentListOfRoom(roomId uint) (*RoomPriceList, error)
}

type roomPriceRepo struct{ db *gorm.DB }

func NewRoomPriceRepo(db *gorm.DB) RoomPriceRepo {
	return &roomPriceRepo{db}
}

func (r *roomPriceRepo) CreateList(priceList *RoomPriceList) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(priceList).Error; err != nil {
			return err
		}

		if err := tx.Model(&Room{}).
			Where("id = ?", priceList.RoomID).
			Update("price_list_id", priceList.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *roomPriceRepo) FindListById(id uint) (*RoomPriceList, error) {
	var list RoomPriceList
	err := r.db.Preload("Items").Where("id = ?", id).First(&list).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
}

func (r *roomPriceRepo) FindListsByRoomId(roomId uint) ([]RoomPriceList, error) {
	var lists []RoomPriceList
	err := r.db.Preload("Items").Where("room_id = ?", roomId).Find(&lists).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *roomPriceRepo) FindCurrentListOfRoom(roomId uint) (*RoomPriceList, error) {
	var latest RoomPriceList
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
