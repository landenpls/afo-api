package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type img_response struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func default_enpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "afo suck cock")
}

func fox_endpoint(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	_id := strconv.Itoa(rand.Intn(num_fox-1) + 1)
	_url := "http://" + r.Host + "/img/fox/fox_" + _id + ".jpg"
	json.NewEncoder(w).Encode([]img_response{
		{Id: _id, Url: _url},
	})
}

var num_fox int

func main() {
	files, _ := ioutil.ReadDir("./imgs/fox")
	num_fox = len(files)
	fmt.Println("Currently serving", num_fox, "pictures of foxes :3")

	/* router */
	router := mux.NewRouter().StrictSlash(true)
	/* static */
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./imgs/"))))
	router.Handle("/", http.FileServer(http.Dir("./static/")))
	/* dynamic */
	router.HandleFunc("/fox", fox_endpoint)
	/* 404 */
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "afo suck cock")
	})
	/* 405 */
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Method not allowed")
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}
