package server

import "github.com/guricerin/stop-now-smoking/entity"

type UserViewModel struct {
	Name      string
	AccountId string
}

func toUserViewModel(u entity.User) UserViewModel {
	vm := UserViewModel{
		Name:      u.Name,
		AccountId: u.AccountId,
	}
	return vm
}

type ViewModel struct {
	LoginState LoginState
	LoginUser  UserViewModel
	RsrcUser   UserViewModel
}
