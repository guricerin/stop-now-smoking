package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/guricerin/stop-now-smoking/entity"
	. "github.com/guricerin/stop-now-smoking/util"
	"github.com/julienschmidt/httprouter"
)

// GET /
func (s *Server) index(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	user, _, err := s.fetchAccountFromCookie(req)
	if err == nil {
		// ログイン済み
		vm := ViewModel{
			LoginUser: toLoginUserViewModel(user),
		}
		writeHtml(w, vm, "layout", "navbar.prv", "index")
	} else {
		// 未ログイン
		writeHtml(w, nil, "layout", "navbar.pub", "index")
	}
}

// GET /login
func (s *Server) showLogin(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	writeHtml(w, nil, "layout", "navbar.pub", "login")
}

// POST /login
func (s *Server) authenticate(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	err := req.ParseForm()
	if err != nil {
		Elog.Printf("ParseForm() error: %v", err)
		return
	}
	plainPassword := req.PostFormValue("password")
	accountId := req.PostFormValue("account_id")
	user, err := s.userStore.RetrieveByAccountId(accountId)
	if err != nil {
		// account_id がちがう。が、「account_id or password is wrong」と表示する
		Ilog.Printf("@%s error: %v", accountId, err)
		msg := "アカウントIDまたはパスワードが間違っています。"
		vm := ViewModel{
			Error: toErrorViewModel(msg),
		}
		writeHtml(w, vm, "layout", "navbar.pub", "login")
		return
	}

	if entity.VerifyPasswordHash(user.Password, plainPassword) {
		sess, err := s.sessionStore.Create(user)
		if err != nil {
			Elog.Printf("@%s : %s", user.AccountId, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			Ilog.Printf("@%s : login success", user.AccountId)
			cookie := s.createCookie(sess.Uuid)
			http.SetCookie(w, &cookie)
			http.Redirect(w, req, "/", http.StatusFound)
		}
	} else {
		Ilog.Printf("@%s : login failed. account_id or password is wrong.", user.AccountId)
		msg := "アカウントIDまたはパスワードが間違っています。"
		vm := ViewModel{
			Error: toErrorViewModel(msg),
		}
		writeHtml(w, vm, "layout", "navbar.pub", "login")
	}
}

// GET /logout
func (s *Server) logout(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	cookie, err := req.Cookie("_cookie")
	if err != http.ErrNoCookie {
		Ilog.Println("session delete")
		uuid := cookie.Value
		err = s.sessionStore.DeleteByUuid(uuid)
		if err != nil {
			Wlog.Printf("delete session by uuid error: %v", err)
		}
	}
	s.deleteCookie(w)
	http.Redirect(w, req, "/", http.StatusFound)
}

// GET /signup
func (s *Server) showSignup(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	writeHtml(w, nil, "layout", "navbar.pub", "signup")
}

// POST /signup
func (s *Server) createUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	err := req.ParseForm()
	if err != nil {
		Elog.Printf("ParseForm() error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accountName := req.FormValue("name")
	if !entity.VerifyAccountName(accountName) {
		msg := "アカウント名に使用可能な文字列は、8文字以上255文字以下です。"
		Ilog.Println(msg)
		vm := ViewModel{}
		vm.Error = toErrorViewModel(msg)
		writeHtml(w, vm, "layout", "navbar.pub", "signup")
		return
	}

	accountId := req.FormValue("account_id")
	if !entity.VerifyAccountId(accountId) {
		msg := "アカウントIDに使用可能な文字列は、8文字以上255文字以下の半角英数字とアンダーバー（_）です。"
		Ilog.Println(msg)
		vm := ViewModel{}
		vm.Error = toErrorViewModel(msg)
		writeHtml(w, vm, "layout", "navbar.pub", "signup")
		return
	}

	plainPassword := req.FormValue("password")
	if !entity.VerifyPlainPassword(plainPassword) {
		msg := "パスワードに使用可能な文字列は、8文字以上255文字以下の半角英数字です。"
		Ilog.Println(msg)
		vm := ViewModel{}
		vm.Error = toErrorViewModel(msg)
		writeHtml(w, vm, "layout", "navbar.pub", "signup")
		return
	}

	hashedPassword, err := entity.EncryptPassword(plainPassword)
	if err != nil {
		Elog.Printf("EncryptPassword() error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := entity.User{
		Name:      accountName,
		AccountId: accountId,
		Password:  hashedPassword,
	}
	if s.userStore.CheckAccountIdExists(user) {
		msg := fmt.Sprintf("%v というアカウントIDは既に使用されています。", user.AccountId)
		Ilog.Printf(msg)
		vm := ViewModel{}
		vm.Error = toErrorViewModel(msg)
		writeHtml(w, vm, "layout", "navbar.pub", "signup")
		return
	}
	user, err = s.userStore.Create(user)
	if err != nil {
		Elog.Printf("create user error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sess, err := s.sessionStore.Create(user)
	if err != nil {
		Elog.Printf("create session error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := s.createCookie(sess.Uuid)
	http.SetCookie(w, &cookie)
	http.Redirect(w, req, "/", http.StatusFound)
}

// GET /delete-account
func (s *Server) showDeleteAccount(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	user, _, err := s.fetchAccountFromCookie(req)
	if err != nil {
		Wlog.Printf("guest access GET /delete-account: %v", err)
		return
	}

	vm := ViewModel{
		LoginUser: toLoginUserViewModel(user),
	}
	writeHtml(w, vm, "layout", "navbar.prv", "delete-account")
}

// POST /delete-account
func (s *Server) deleteAccount(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	user, sess, err := s.fetchAccountFromCookie(req)
	if err != nil {
		Wlog.Printf("guest access POST /delete-account: %v", err)
		return
	}
	err = req.ParseForm()
	if err != nil {
		Elog.Printf("ParseForm() error: %v", err)
		return
	}

	r := req.FormValue("delete")
	switch r {
	case "yes":
		if err = s.userStore.DeleteById(user.Id); err != nil {
			Elog.Printf("%v", err)
			return
		}
		if err = s.sessionStore.DeleteByUuid(sess.Uuid); err != nil {
			Elog.Printf("%v", err)
			return
		}
		if err = s.cigaretteStore.DeleteAllByUserId(user.Id); err != nil {
			Elog.Printf("%v", err)
			return
		}
		if err = s.followStore.DeleteAllByAccountId(user); err != nil {
			Elog.Printf("%v", err)
			return
		}

		Ilog.Printf("@%s: delete account", user.AccountId)
		s.deleteCookie(w)
		http.Redirect(w, req, "/", http.StatusFound)
	default:
		Ilog.Printf("@%s: not delete account", user.AccountId)
		http.Redirect(w, req, "/", http.StatusFound)
	}
}

// GET /users/:account_id
func (s *Server) userPage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	vm, rsrcUser := s.userRsrcViewModel(req, ps)
	switch vm.LoginState {
	case RsrcNotFound:
		http.NotFound(w, req)
		return
	}

	cigarettes, err := s.cigaretteStore.RetrieveAllByUserId(rsrcUser.Id)
	if err != nil {
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vm.RsrcUser.TotalSmokedCountAllDate = entity.TotalSmokedCountAllDate(cigarettes)
	now := time.Now()
	vm.RsrcUser.TotalSmokedCountToday = entity.TotalSmokedCountByDate(cigarettes, now)
	cigarettesByWeek, err := s.cigaretteStore.RetrieveAllByUserIdAndBetweenDate(rsrcUser.Id, now.AddDate(0, 0, -6), now)
	if err != nil {
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vm.RsrcUser.TotalSmokedByDate = totalsSmokedByDateViewModel(cigarettesByWeek)

	// URLクエリで指定された日付範囲内のデータフェッチ
	startDate, endDate, err := s.parseStartAndEndDateQuery(req)
	if err == nil {
		// 日付範囲内のデータフェッチ
		Dlog.Printf("%v, %v", startDate, endDate)
		cigarettesByDate, err := s.cigaretteStore.RetrieveAllByUserIdAndBetweenDate(rsrcUser.Id, startDate, endDate)
		if err != nil {
			Elog.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		vm.RsrcUser.TotalSmokedByDate = totalsSmokedByDateViewModel(cigarettesByDate)
	}

	switch vm.LoginState {
	case Guest:
		writeHtml(w, vm, "layout", "navbar.pub", "user-page")
	case LoginAndRsrcUser:
		writeHtml(w, vm, "layout", "navbar.prv", "user-page")
	case LoginButNotRsrcUser:
		writeHtml(w, vm, "layout", "navbar.prv", "user-page")
	default:
		err := fmt.Errorf("unexhausted LogState enum: %v", vm.LoginState)
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /users/:account_id/setting
func (s *Server) showUserSetting(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	vm, _ := s.userRsrcViewModel(req, ps)
	switch vm.LoginState {
	case LoginAndRsrcUser:
		// ok
		writeHtml(w, vm, "layout", "navbar.prv", "user-setting")
	case Guest:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	case LoginButNotRsrcUser:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	case RsrcNotFound:
		http.NotFound(w, req)
	default:
		err := fmt.Errorf("unexhausted LogState enum: %v", vm.LoginState)
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /users/:account_id/setting
func (s *Server) editUserSetting(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	vm, rsrcUser := s.userRsrcViewModel(req, ps)
	switch vm.LoginState {
	case LoginAndRsrcUser:
		err := req.ParseForm()
		if err != nil {
			Elog.Printf("%v", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		accountName := req.FormValue("account_name")
		if !entity.VerifyAccountName(accountName) {
			msg := "アカウント名に使用可能な文字列は、8文字以上255文字以下です。"
			Ilog.Println(msg)
			vm.Error = toErrorViewModel(msg)
			writeHtml(w, vm, "layout", "navbar.prv", "user-setting")
			return
		}

		favBrand := req.FormValue("favorite_brand")
		if !entity.VerifyFavoriteBrand(favBrand) {
			msg := "好きな銘柄に使用可能な文字列は、0文字以上255文字以下です。"
			Ilog.Println(msg)
			vm.Error = toErrorViewModel(msg)
			writeHtml(w, vm, "layout", "navbar.prv", "user-setting")
			return
		}

		newRsrcUser := entity.User{
			Id:            rsrcUser.Id,
			Name:          accountName,
			AccountId:     rsrcUser.AccountId,
			FavoriteBrand: favBrand,
			Password:      rsrcUser.Password,
		}
		err = s.userStore.Update(newRsrcUser)
		if err != nil {
			Elog.Printf("%v", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		// ok
		url := fmt.Sprintf("/users/%s", newRsrcUser.AccountId)
		http.Redirect(w, req, url, http.StatusFound)
		writeHtml(w, vm, "layout", "navbar.prv", "user-setting")
	case Guest:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	case LoginButNotRsrcUser:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	case RsrcNotFound:
		http.NotFound(w, req)
	default:
		err := fmt.Errorf("unexhausted LogState enum: %v", vm.LoginState)
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /users/:account_id/edit-cigarette-today
func (s *Server) editCigaretteToday(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	vm, loginUser := s.userRsrcViewModel(req, ps)
	switch vm.LoginState {
	case LoginAndRsrcUser:
		err := req.ParseForm()
		if err != nil {
			Elog.Printf("%v", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		smoked_count, err := strconv.Atoi(req.FormValue("smoked_count"))
		if err != nil {
			Elog.Printf("%v", err)
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
		if smoked_count < 0 {
			Elog.Printf("%v", err)
			http.Error(w, "403 forbidden", http.StatusForbidden)
			return
		}

		cigarette := entity.Cigarette{
			SmokedCount: uint(smoked_count),
			UserId:      loginUser.Id,
			CreatedAt:   time.Now(),
		}
		// duplicate on key を使わないのは
		// UniqueIdを知るために結局Retrieveする必要があるから。
		exist, err := s.cigaretteStore.ExistByUserIdAndDate(cigarette)
		if !exist || err != nil {
			Dlog.Printf("not exist, so create cigarette.: %v", err)
			err := s.cigaretteStore.Create(cigarette)
			if err != nil {
				Elog.Printf("%v", err)
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			Dlog.Printf("exist, so update cigarette.")
			err = s.cigaretteStore.UpdateByUserIdAndDate(cigarette)
			if err != nil {
				Elog.Printf("%v", err)
				http.Error(w, "500 internal server error", http.StatusInternalServerError)
				return
			}
		}
		url := fmt.Sprintf("/users/%s", vm.LoginUser.AccountId)
		http.Redirect(w, req, url, http.StatusFound)
	case RsrcNotFound:
		http.NotFound(w, req)
	case Guest:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	case LoginButNotRsrcUser:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	default:
		err := fmt.Errorf("unexhausted LogState enum: %v", vm.LoginState)
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /users/:account_id/follow/:dst_account_id
func (s *Server) follow(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	dstAccountId := ps.ByName("dst_account_id")
	vm, loginUser := s.userRsrcViewModel(req, ps)
	switch vm.LoginState {
	case LoginAndRsrcUser:
		if dstAccountId == loginUser.AccountId { // 自分自身はフォローできない
			str := fmt.Sprintf("can't follow yourself: account_id = %v", dstAccountId)
			err := errors.New(str)
			Elog.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		err := s.followStore.Create(loginUser.AccountId, dstAccountId)
		if err != nil {
			Elog.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		url := req.Header.Get("Referer") // 直前のURL
		http.Redirect(w, req, url, http.StatusFound)
	case RsrcNotFound:
		http.NotFound(w, req)
	case Guest:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	case LoginButNotRsrcUser:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	default:
		err := fmt.Errorf("unexhausted LogState enum: %v", vm.LoginState)
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST /users/:account_id/unfollow/:dst_account_id
func (s *Server) unfollow(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	dstAccountId := ps.ByName("dst_account_id")
	vm, loginUser := s.userRsrcViewModel(req, ps)
	switch vm.LoginState {
	case LoginAndRsrcUser:
		if dstAccountId == loginUser.AccountId { // 自分自身はフォロー解除できない
			str := fmt.Sprintf("can't unfollow yourself: account_id = %v", dstAccountId)
			err := errors.New(str)
			Elog.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		err := s.followStore.Delete(loginUser.AccountId, dstAccountId)
		if err != nil {
			Elog.Printf("%v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		url := req.Header.Get("Referer") // 直前のURL
		http.Redirect(w, req, url, http.StatusFound)
	case RsrcNotFound:
		http.NotFound(w, req)
	case Guest:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	case LoginButNotRsrcUser:
		http.Error(w, "403 forbidden", http.StatusForbidden)
	default:
		err := fmt.Errorf("unexhausted LogState enum: %v", vm.LoginState)
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /users/:account_id/follows
func (s *Server) showFollows(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	vm, _ := s.userRsrcViewModel(req, ps)
	switch vm.LoginState {
	case RsrcNotFound:
		http.NotFound(w, req)
		return
	}

	err := s.fetchSmokedCountTodayForFollows(&vm)
	if err != nil {
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch vm.LoginState {
	case Guest:
		writeHtml(w, vm, "layout", "navbar.pub", "follows-list")
	case LoginAndRsrcUser:
		writeHtml(w, vm, "layout", "navbar.prv", "follows-list")
	case LoginButNotRsrcUser:
		writeHtml(w, vm, "layout", "navbar.prv", "follows-list")
	default:
		err := fmt.Errorf("unexhausted LogState enum: %v", vm.LoginState)
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GET /users/:account_id/followers
func (s *Server) showFollowers(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	vm, _ := s.userRsrcViewModel(req, ps)
	switch vm.LoginState {
	case RsrcNotFound:
		http.NotFound(w, req)
		return
	}

	err := s.fetchSmokedCountTodayForFollowers(&vm)
	if err != nil {
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch vm.LoginState {
	case Guest:
		writeHtml(w, vm, "layout", "navbar.pub", "followers-list")
	case LoginAndRsrcUser:
		writeHtml(w, vm, "layout", "navbar.prv", "followers-list")
	case LoginButNotRsrcUser:
		writeHtml(w, vm, "layout", "navbar.prv", "followers-list")
	default:
		err := fmt.Errorf("unexhausted LogState enum: %v", vm.LoginState)
		Elog.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GET /search-account
func (s *Server) searchAccount(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	s.accessLog(req)
	qAccountId := req.URL.Query().Get("account_id")
	resultUsers, err := s.userStore.SearchAllByAccountId(qAccountId)
	if err != nil {
		Elog.Printf("%v", err)
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}

	vm := ViewModel{
		SearchedUsers: toSearchedUsersViewModel(qAccountId, resultUsers),
	}
	loginUser, _, err := s.fetchAccountFromCookie(req)
	if err != nil {
		writeHtml(w, vm, "layout", "navbar.pub", "search-account-result")
	} else {
		vm.LoginUser = LoginUserViewModel{
			Name:      loginUser.Name,
			AccountId: loginUser.AccountId,
		}
		writeHtml(w, vm, "layout", "navbar.prv", "search-account-result")
	}
}
