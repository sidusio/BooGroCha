package main

import (
	"log"

	"sidus.io/boogrocha/internal/srv"
)

func main() {
	s := srv.NewServer([]byte("abababababababababababababababab"))
	log.Fatal(s.Run(":8080"))
}
