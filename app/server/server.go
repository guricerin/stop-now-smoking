package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/guricerin/stop-now-smoking/entity"
	. "github.com/guricerin/stop-now-smoking/util"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	cfg          *Config
	router       http.Handler
	userStore    *userStore
	sessionStore *sessionStore
}

func NewServer(cfg *Config, db DbDriver) *Server {
	userStore := NewUserStore(db)
	sessionStore := NewSessionStore(db)
	server := Server{
		cfg:          cfg,
		userStore:    userStore,
		sessionStore: sessionStore,
	}
	server.setupRouter()
	return &server
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%s", s.cfg.ServerHost, s.cfg.ServerPort)
	server := http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	return server.ListenAndServe()
}

func (s *Server) setupRouter() {
	router := httprouter.New()

	router.GET("/", s.index)
	router.GET("/login", s.showLogin)
	router.POST("/login", s.authenticate)
	router.GET("/signup", s.showSignup)
	router.POST("/signup", s.createUser)

	s.router = router
}

func accessLog(req *http.Request) {
	Ilog.Printf("%s %s", req.Method, req.URL)
}

// GET /
func (s *Server) index(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	accessLog(req)
	writeHtml(w, nil, "layout", "navbar.pub", "index")
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

func (s *Server) fetchAccountFromCookie(req *http.Request) (user entity.User, sess entity.Session, err error) {
	cookie, err := req.Cookie("_cookie")
	if err != nil {
		return
	}

	uuid := cookie.Value
	sess, err = s.sessionStore.RetrieveByUuid(uuid)
	if err != nil {
		return
	}

	user, err = s.userStore.RetrieveById(sess.UserId)
	return
}

func (s *Server) createCookie(uuid string) http.Cookie {
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    uuid,
		HttpOnly: true, // JavaScriptなど非HTTPのAPIを禁止
	}
	return cookie
}

func (s *Server) deleteCookie(w http.ResponseWriter) {
	// すぐに寿命が尽きるクッキーで上書きすることで、結果的に削除したことになる
	cookie := http.Cookie{
		Name:    "_cookie",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, &cookie)
}
