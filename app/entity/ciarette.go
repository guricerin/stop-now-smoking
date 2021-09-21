package entity

import "time"

type Cigarette struct {
	Id int64
	// 吸った本数。
	// 負数を許容しているのは、誤記録修正のため。
	// 集計の際は一日あたりの合計が負になる場合に0に直すこと。
	SmokedCount int
	UserId      int64
	CreatedAt   time.Time
}

func TotalSmokedCount(cigarettes []Cigarette) uint {
	var res = 0
	for _, cig := range cigarettes {
		res += cig.SmokedCount
	}
	if res < 0 {
		res = 0
	}
	return uint(res)
}

func TotalSmokedCountAllDate(cigarettes []Cigarette) uint {
	group := GroupCigarettesByDate(cigarettes)
	var res uint = 0
	for _, v := range group {
		res += TotalSmokedCount(v)
	}
	return res
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// 指定日付の喫煙記録のみを集計
func TotalSmokedCountByDate(cigarettes []Cigarette, t time.Time) uint {
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
	return uint(res)
}

// 日付ごとに喫煙記録をグルーピング
func GroupCigarettesByDate(cigarettes []Cigarette) map[time.Time][]Cigarette {
	res := make(map[time.Time][]Cigarette)
	for _, cig := range cigarettes {
		t := cig.CreatedAt
		dt := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, jst)
		_, ok := res[dt]
		if !ok {
			res[dt] = make([]Cigarette, 1)
		}
		res[dt] = append(res[dt], cig)
	}
	return res
}
