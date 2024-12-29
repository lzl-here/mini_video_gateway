package repo

import "gorm.io/gorm"

type DBRepo struct {
	db *gorm.DB
}

func NewDBRepo(db *gorm.DB) *DBRepo {
	return &DBRepo{db: db}
}
