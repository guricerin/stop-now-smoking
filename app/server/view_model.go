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
	Name      string
	AccountId string
}

type RsrcUserViewModel struct {
	Name                    string
	AccountId               string
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
		Name:      u.Name,
		AccountId: u.AccountId,
	}
}

func toFollowViewModels(us []entity.User) []FollowViewModel {
	// res := make([]FollowViewModel, len(us)) // todo: こっちだと余計な空データが先頭に混入する。原因不明
	res := make([]FollowViewModel, 0)
	for _, u := range us {
		res = append(res, toFollowViewModel(u))
	}
	return res
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
	for _, user := range users {
		res.Results = append(res.Results, SearchedUserViewModel{
			Name:      user.Name,
			AccountId: user.AccountId,
		})
	}
	return res
}

type ViewModel struct {
	LoginState    LoginState
	LoginUser     LoginUserViewModel
	RsrcUser      RsrcUserViewModel
	SearchedUsers SearchedUsersViewModel
}
