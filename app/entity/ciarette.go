package entity

import "time"

type Cigarette struct {
	Id int64
	// 吸った本数。
	// 負数を許容しているのは、誤記録修正のため。
	// 集計の際は合計が負になる場合に0に直すこと。
	SmokedCount int
	UserId      int64
	CreatedAt   time.Time
}

func TotalSmokedCount(cigarettes []Cigarette) int {
	var res = 0
	for _, cig := range cigarettes {
		res += cig.SmokedCount
	}
	if res < 0 {
		res = 0
	}
	return res
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func TotalSmokedCountByDate(cigarettes []Cigarette, t time.Time) int {
	var res = 0
	dt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, jst)
	for _, cig := range cigarettes {
		ctime := cig.CreatedAt
		cdt := time.Date(ctime.Year(), ctime.Month(), ctime.Day(), 0, 0, 0, 0, jst)
		// 年月日のみ一致かを判定
		if dt.Equal(cdt) {
			res += cig.SmokedCount
		}
	}
	if res < 0 {
		res = 0
	}
	return res
}
