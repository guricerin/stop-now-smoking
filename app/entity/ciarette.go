package entity

import "time"

type Cigarette struct {
	Id int64
	// 吸った本数
	SmokedCount int
	UserId      int64
	CreatedAt   time.Time
}
