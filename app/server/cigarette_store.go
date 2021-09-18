package server

import "time"

type cigaretteTable struct {
	Id          int64
	SmokedCount int
	UserId      int64
	CreatedAd   time.Time
}

type cigaretteStore struct {
	db DbDriver
}

func NewCigaretteStore(db DbDriver) *cigaretteStore {
	return &cigaretteStore{db: db}
}
