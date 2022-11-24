package main

import (
	"fmt"
	"github.com/ThePratikSah/gomongoapi/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	err := http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Serving at PORT: 3000")
}
