package main

import (
	"io"
	"log"
	"net/http"

	"github.com/Komosa/httpmap"
)

func main() {
	m := httpmap.New()
	m.Get("hello/:user", httpmap.HandlerFunc(func(w http.ResponseWriter, _ *http.Request, named httpmap.Named) {
		io.WriteString(w, "Hello, "+named.Param("user")+"!")
	}))

	err := http.ListenAndServe(":8001", m)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// ask for localhost:8001/hello/:juliett

}
