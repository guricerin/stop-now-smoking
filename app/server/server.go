package server

import (
	"net/http"

	u "github.com/guricerin/stop-now-smoking/util"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router http.Handler
}

func NewServer() *Server {
	server := Server{}
	server.setupRouter()
	return &server
}

func (s *Server) Run() error {
	server := http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}
	return server.ListenAndServe()
}

func (s *Server) setupRouter() {
	router := httprouter.New()

	router.GET("/", s.index)

	s.router = router
}

func (s *Server) index(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	u.Ilog.Printf("index")
	writeHtml(w, nil, "layout", "navbar.pub", "index")
}
