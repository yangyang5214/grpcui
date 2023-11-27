package pkg

import (
	"net/http"
	"os"
	"path"
)

type GrpcWeb struct {
	curl *GrpcCurl
	pwd  string
}

func NewGrpcWeb(curl *GrpcCurl) *GrpcWeb {
	pwd, _ := os.Getwd()
	return &GrpcWeb{
		curl: curl,
		pwd:  pwd,
	}
}

func (s *GrpcWeb) home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(s.pwd, "./index.html"))
}

func (s *GrpcWeb) handleMethod(writer http.ResponseWriter, request *http.Request) {

}

func (s *GrpcWeb) faviconIco(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join(s.pwd, "./static/favicon.png"))
}

func (s *GrpcWeb) Handler(writer http.ResponseWriter, request *http.Request) {
	switch request.RequestURI {
	case "/":
		s.home(writer, request)
	case "/favicon.ico":
		s.faviconIco(writer, request)
	default:
		s.handleMethod(writer, request)
	}
}
