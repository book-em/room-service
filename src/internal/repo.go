package internal

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(room *Room) error
	FindById(id uint) (*Room, error)
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

func (r *repository) FindById(id uint) (*Room, error) {
	var room Room
	err := r.db.Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}
