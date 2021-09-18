package server

import (
	"net/http"

	"github.com/guricerin/stop-now-smoking/entity"
	. "github.com/guricerin/stop-now-smoking/util"
	"github.com/julienschmidt/httprouter"
)

// GET /
func (s *Server) index(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	accessLog(req)
	user, _, err := s.fetchAccountFromCookie(req)
	if err == nil {
		// ログイン済み
		vm := ViewModel{
			LoginUser: toUserViewModel(user),
		}
		writeHtml(w, vm, "layout", "navbar.prv", "index")
	} else {
		// 未ログイン
		writeHtml(w, nil, "layout", "navbar.pub", "index")
	}
}

// GET /login
func (s *Server) showLogin(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	accessLog(req)
	writeHtml(w, nil, "layout", "navbar.pub", "login")
}

// POST /login
func (s *Server) authenticate(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	accessLog(req)
	err := req.ParseForm()
	if err != nil {
		Elog.Printf("ParseForm() error: %v", err)
		return
	}
	// todo: 8文字以上の英数字
	plainPassword := req.PostFormValue("password")
	accountId := req.PostFormValue("account_id")
	user, err := s.userStore.RetrieveByAccountId(accountId)
	if err != nil {
		// todo: account_id がちがう。が、「account_id or password is wrong」と表示する
		Elog.Printf("@%s error: %v", accountId, err)
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
		//todo
		Ilog.Printf("@%s : login failed. account_id or password is wrong.", user.AccountId)
	}
}

// GET /logout
func (s *Server) logout(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	accessLog(req)
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
	accessLog(req)
	writeHtml(w, nil, "layout", "navbar.pub", "signup")
}

// POST /signup
func (s *Server) createUser(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	accessLog(req)
	err := req.ParseForm()
	if err != nil {
		Elog.Printf("ParseForm() error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hashedPassword, err := entity.EncryptPassword(req.FormValue("password"))
	if err != nil {
		Elog.Printf("EncryptPassword() error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := entity.User{
		Name:      req.PostFormValue("name"),
		AccountId: req.PostFormValue("account_id"),
		Password:  hashedPassword,
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
	accessLog(req)
	user, _, err := s.fetchAccountFromCookie(req)
	if err != nil {
		Wlog.Printf("guest access GET /delete-account: %v", err)
		return
	}

	vm := ViewModel{
		LoginUser: toUserViewModel(user),
	}
	writeHtml(w, vm, "layout", "navbar.prv", "delete-account")
}

// POST /delete-account
func (s *Server) deleteAccount(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	accessLog(req)
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
	accessLog(req)
	vm := s.userRsrcViewModel(req, ps)
	switch vm.LoginState {
	case RsrcNotFound:
		http.NotFound(w, req)
	case Guest:
		writeHtml(w, vm, "layout", "navbar.pub", "user-page")
	case LoginAndRsrcUser:
		writeHtml(w, vm, "layout", "navbar.prv", "user-page")
	case LoginButNotRsrcUser:
		writeHtml(w, vm, "layout", "navbar.prv", "user-page")
	}
}