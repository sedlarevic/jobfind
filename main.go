package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func crawl(w http.ResponseWriter, r *http.Request) {

}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", crawl)
}
