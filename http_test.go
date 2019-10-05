package jsonrpc

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"net/rpc"
	"testing"
)

type Args struct {
	Message string
}

type Service int

func (s *Service) Echo(args *Args, reply *string) error {
	*reply = args.Message
	return nil
}

func TestHandler(t *testing.T) {
	service := new(Service)
	server := rpc.NewServer()
	server.Register(service)
	handler := Handler(server)

	requestData := `{"method":"Service.Echo","params":[{"Message": "Hello"}], "id": 0}`

	r := httptest.NewRequest("POST", "http://example.com", bytes.NewReader([]byte(requestData)))
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	response := w.Result()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Errorf("%s", err)
	}

	result := `{"id":0,"result":"Hello","error":null}`

	if string(body) != result+"\n" {
		t.Errorf("got: %s, want: %s", body, result)
	}
}
