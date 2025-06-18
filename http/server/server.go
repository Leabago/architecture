package main

import (
	"fmt"
	"net/http"
	"os"
)

func hello(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(rw, "Method is not supported", http.StatusNotFound)
		return
	}

	fmt.Fprintln(rw, "hello")
}

func postHello(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(rw, "Method is not supported", http.StatusNotFound)
		return
	}

	fmt.Fprintln(rw, "postHello")
}

func main() {

	port := ":8080"
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/post", postHello)

	srv := &http.Server{
		Addr: port,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// if err := http.ListenAndServe(port, nil); err != nil {
	// 	fmt.Println(err)
	// }
}
