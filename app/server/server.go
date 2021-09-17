package server

import (
	"net/http"

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

func (s *Server) Run() {
	server := http.Server{
		Addr:    ":8000",
		Handler: s.router,
	}
	server.ListenAndServe()
}

func (s *Server) setupRouter() {
	router := httprouter.New()

	router.GET("/", s.index)

	s.router = router
}

func (s *Server) index(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	writeHtml(w, nil, "layout", "navbar.pub", "index")
}
