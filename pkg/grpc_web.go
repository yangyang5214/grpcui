package pkg

import (
	"encoding/json"
	"github.com/jhump/protoreflect/dynamic"
	"io"
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

	r := map[string]any{
		"services": services,
		"server":   s.curl.GetAddr(),
	}

	err = json.NewEncoder(writer).Encode(r)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

type SendRequest struct {
	Payload map[string]any `json:"payload,omitempty"`
	Method  string         `json:"method,omitempty"`
}

func (s *GrpcWeb) Send(writer http.ResponseWriter, request *http.Request) {
	headerContentType := request.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(writer, "Content Type is not application/json", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var sendRequest *SendRequest
	err = json.Unmarshal(body, &sendRequest)
	if err != nil {
		http.Error(writer, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	payload := sendRequest.Payload
	method := sendRequest.Method

	md, err := s.curl.GetMethodDescByName(method)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	dMsg := dynamic.NewMessage(md.GetInputType())
	err = JsonToMessage(payload, dMsg)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := s.client.send(request.Context(), md, dMsg)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = writer.Write([]byte(resp))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
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
	r := map[string]any{
		"payload": payload,
	}
	err = json.NewEncoder(writer).Encode(r)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
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

func (s *GrpcWeb) AllMethods(writer http.ResponseWriter, request *http.Request) {
	services, err := s.curl.ListServices()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	var methodWraps []*MethodWrap
	for _, service := range services {
		methods, err := s.curl.ListMethods(service)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		methodWraps = append(methodWraps, &MethodWrap{
			Service: service,
			Methods: methods,
		})
	}

	r := &AllMethod{
		Addr:    s.curl.GetAddr(),
		Methods: methodWraps,
	}
	err = json.NewEncoder(writer).Encode(r)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
