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
	cfg            *Config
	router         http.Handler
	userStore      *userStore
	sessionStore   *sessionStore
	cigaretteStore *cigaretteStore
}

func NewServer(cfg *Config, db DbDriver) *Server {
	userStore := NewUserStore(db)
	sessionStore := NewSessionStore(db)
	cigaretteStore := NewCigaretteStore(db)
	server := Server{
		cfg:            cfg,
		userStore:      userStore,
		sessionStore:   sessionStore,
		cigaretteStore: cigaretteStore,
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

	// asset
	router.ServeFiles("/static/*filepath", http.Dir("static/"))

	router.GET("/", s.index)
	router.GET("/login", s.showLogin)
	router.POST("/login", s.authenticate)
	router.GET("/logout", s.logout)
	router.GET("/signup", s.showSignup)
	router.POST("/signup", s.createUser)
	router.GET("/delete-account", s.showDeleteAccount)
	router.POST("/delete-account", s.deleteAccount)
	router.GET("/users/:account_id", s.userPage)

	s.router = router
}

func accessLog(req *http.Request) {
	Dlog.Printf("%s %s", req.Method, req.URL)
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

// /users/:account_id
func (s *Server) userRsrcViewModel(req *http.Request, ps httprouter.Params) (vm ViewModel) {
	accountId := ps.ByName("account_id")
	rsrcUser, err := s.userStore.RetrieveByAccountId(accountId)
	if err != nil {
		Ilog.Printf("rsrc user not found: %v", err)
		vm.LoginState = RsrcNotFound
		return
	}
	vm.RsrcUser = toUserViewModel(rsrcUser)

	loginUser, _, err := s.fetchAccountFromCookie(req)
	if err != nil {
		Ilog.Printf("access user is guest: %v", err)
		vm.LoginState = Guest
		return
	}
	vm.LoginUser = toUserViewModel(loginUser)

	if loginUser == rsrcUser {
		vm.LoginState = LoginAndRsrcUser
	} else {
		vm.LoginState = LoginButNotRsrcUser
	}
	return
}
