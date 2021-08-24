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

var port string = "8080" // have to port-forward this

func random_fn(category string) (string, int) {
	files, _ := ioutil.ReadDir("./imgs/" + category)
	_selection := rand.Intn(len(files)-1) + 1
	return files[_selection].Name(), _selection
}

type img_response struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func default_enpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "afo suck cock")
}

func fox_endpoint(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	_fn, _id := random_fn("fox")
	_url := "https://" + r.Host + "/img/fox/" + _fn
	json.NewEncoder(w).Encode([]img_response{
		{Id: strconv.Itoa(_id), Url: _url},
	})
}
func yiff_endpoint(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	_fn, _id := random_fn("fox")
	_url := "https://" + r.Host + "/img/yiff/" + _fn
	json.NewEncoder(w).Encode([]img_response{
		{Id: strconv.Itoa(_id), Url: _url},
	})
}

var num_fox int
var num_yiff int

func main() {
	files, _ := ioutil.ReadDir("./imgs/fox")
	num_fox = len(files)
	fmt.Println("Currently serving", num_fox, "pictures of foxes :3")
	files, _ = ioutil.ReadDir("./imgs/yiff")
	num_yiff = len(files)
	fmt.Println("Currently serving", num_yiff, "yiff and counting owo")

	/* router */
	router := mux.NewRouter().StrictSlash(true)
	/* static */
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./imgs/"))))
	router.Handle("/", http.FileServer(http.Dir("./static/")))
	/* dynamic */
	router.HandleFunc("/fox", fox_endpoint)
	router.HandleFunc("/yiff", yiff_endpoint)
	/* 404 */
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "afo suck cock")
	})
	/* 405 */
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Method not allowed")
	})

	/* listen */
	log.Fatal(http.ListenAndServe(":"+port, router))
	fmt.Println("Listening on port", port)
}
