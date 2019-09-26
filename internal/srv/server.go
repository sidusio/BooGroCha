package srv

import (
	"net/http"
)

type server struct {
	credentialsSecret []byte
}

func NewServer(credentialsSecret []byte) *server {
	return &server{credentialsSecret: credentialsSecret}
}

func (s *server) Run(address string) error {
	return http.ListenAndServe(address, s.newRouter())
}
