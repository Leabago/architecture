package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	r, err := http.Get("http://localhost:8080/hello")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))

}
