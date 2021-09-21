package server

import (
	"time"

	"github.com/guricerin/stop-now-smoking/entity"
)

type LoginUserViewModel struct {
	Name      string
	AccountId string
}

type totalSmokedByDateViewModel map[time.Time]uint

type RsrcUserViewModel struct {
	Name                    string
	AccountId               string
	TotalSmokedCountAllDate uint
	TotalSmokedCountToday   uint
	TotalSmokedByDate       totalSmokedByDateViewModel
}

func toLoginUserViewModel(u entity.User) LoginUserViewModel {
	vm := LoginUserViewModel{
		Name:      u.Name,
		AccountId: u.AccountId,
	}
	return vm
}

func toRsrcUserViewModel(u entity.User) RsrcUserViewModel {
	vm := RsrcUserViewModel{
		Name:      u.Name,
		AccountId: u.AccountId,
	}
	return vm
}

// todo: 日時が新しい順に並び替える
func totalsSmokedByDateViewModel(cigs []entity.Cigarette) totalSmokedByDateViewModel {
	groupByDate := entity.GroupCigarettesByDate(cigs)
	res := make(totalSmokedByDateViewModel, len(groupByDate))
	for k, v := range groupByDate {
		total := entity.TotalSmokedCount(v)
		res[k] = total
	}
	return res
}

type ViewModel struct {
	LoginState LoginState
	LoginUser  LoginUserViewModel
	RsrcUser   RsrcUserViewModel
}
