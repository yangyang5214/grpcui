package pkg

import (
	"encoding/json"
	"github.com/jhump/protoreflect/dynamic"
	"net/http"
)

type GrpcWeb struct {
	curl   *GrpcCurl
	client *GrpcSend
}

func NewGrpcWeb(curl *GrpcCurl) *GrpcWeb {
	return &GrpcWeb{
		curl:   curl,
		client: NewGrpcSend(curl.GetConn()),
	}
}

func (s *GrpcWeb) ListServices(writer http.ResponseWriter, request *http.Request) {
	services, err := s.curl.ListServices()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(writer).Encode(services)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *GrpcWeb) Send(writer http.ResponseWriter, request *http.Request) {
	payload := request.PostFormValue("payload")
	method := request.PostFormValue("method")
	md, err := s.curl.GetMethodDescByName(method)
	if err != nil {
		http.Error(writer, "method param required", http.StatusBadRequest)
		return
	}

	jsonData := make(map[string]any)
	err = json.Unmarshal([]byte(payload), &jsonData)
	dMsg := dynamic.NewMessage(md.GetInputType())
	err = JsonToMessage(jsonData, dMsg)
	if err != nil {
		http.Error(writer, "method param required", http.StatusBadRequest)
		return
	}
	resp, err := s.client.send(request.Context(), md, dMsg)
	if err != nil {
		http.Error(writer, "method param required", http.StatusBadRequest)
		return
	}
	_, err = writer.Write([]byte(resp))
	if err != nil {
		http.Error(writer, "method param required", http.StatusBadRequest)
		return
	}
}

func (s *GrpcWeb) FakeBody(writer http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()
	method := queryParams.Get("method")

	md, err := s.curl.GetMethodDescByName(method)
	if err != nil {
		http.Error(writer, "method param required", http.StatusBadRequest)
		return
	}
	payload := s.curl.genPayload(md)
	_, err = writer.Write([]byte(payload))
	if err != nil {
		http.Error(writer, "service param required", http.StatusBadRequest)
		return
	}
}

func (s *GrpcWeb) ListMethod(writer http.ResponseWriter, request *http.Request) {
	queryParams := request.URL.Query()
	serviceName := queryParams.Get("service")
	if serviceName == "" {
		http.Error(writer, "service param required", http.StatusBadRequest)
		return
	}
	methods, err := s.curl.ListMethods(serviceName)
	if err != nil {
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = json.NewEncoder(writer).Encode(methods)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
