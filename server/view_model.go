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

type FollowViewModel struct {
	UserId           int64 // 喫煙回数を取得したいためのフィールド
	Name             string
	AccountId        string
	SmokedCountToday uint
}

type RsrcUserViewModel struct {
	Name                    string
	AccountId               string
	FavoriteBrand           string
	TotalSmokedCountAllDate uint
	TotalSmokedCountToday   uint
	TotalSmokedByDate       totalSmokedByDateViewModel
	Follows                 []FollowViewModel
	Followers               []FollowViewModel
	IsFollowedByLoginUser   bool
}

func toLoginUserViewModel(u entity.User) LoginUserViewModel {
	vm := LoginUserViewModel{
		Name:      u.Name,
		AccountId: u.AccountId,
	}
	return vm
}

func toFollowViewModel(u entity.User) FollowViewModel {
	return FollowViewModel{
		UserId:    u.Id,
		Name:      u.Name,
		AccountId: u.AccountId,
	}
}

func toFollowViewModels(us []entity.User) []FollowViewModel {
	res := make([]FollowViewModel, len(us))
	for i, u := range us {
		res[i] = toFollowViewModel(u)
	}
	return res
}

func toRsrcUserViewModel(u entity.User) RsrcUserViewModel {
	favBrand := u.FavoriteBrand
	if favBrand == "" {
		favBrand = "<なし>"
	}

	vm := RsrcUserViewModel{
		Name:          u.Name,
		AccountId:     u.AccountId,
		FavoriteBrand: favBrand,
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

// アカウントID検索結果につかう
type SearchedUserViewModel struct {
	Name      string
	AccountId string
}
type SearchedUsersViewModel struct {
	Query   string
	Results []SearchedUserViewModel
}

func toSearchedUsersViewModel(query string, users []entity.User) SearchedUsersViewModel {
	res := SearchedUsersViewModel{
		Query:   query,
		Results: make([]SearchedUserViewModel, len(users)),
	}
	for i, user := range users {
		res.Results[i] = SearchedUserViewModel{
			Name:      user.Name,
			AccountId: user.AccountId,
		}
	}
	return res
}

type ErrorViewModel struct {
	HasError bool
	Msg      string
}

func toErrorViewModel(msg string) ErrorViewModel {
	return ErrorViewModel{
		HasError: true,
		Msg:      msg,
	}
}

type ViewModel struct {
	LoginState    LoginState
	LoginUser     LoginUserViewModel
	RsrcUser      RsrcUserViewModel
	SearchedUsers SearchedUsersViewModel
	Error         ErrorViewModel
}
