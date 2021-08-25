package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var port string = "8080" // have to port-forward this

var wc http.Client = http.Client{}

func random_fn(category string) (string, int) {
	_files, _ := ioutil.ReadDir("./imgs/" + category)
	_selection := rand.Intn(len(_files)-1) + 1
	return _files[_selection].Name(), _selection
}

type reddit_post struct {
	title           string
	thumbnail       string
	link_flair_type string
	nsfw            bool
	pinned          bool
	is_mod_post     bool
	author          string
	image_url       string
}
type reddit_api_response struct {
	Data struct {
		Children []struct {
			Data struct {
				Title         string                `json:"title"`
				Thumbnail     string                `json:"thumbnail"`
				LinkFlairType string                `json:"link_flair_type"`
				Distinguished string                `json:"distinguished"`
				Author        string                `json:"author"`
				ImageUrl      string                `json:"url_url_overriden_by_dest"`
				Preview       map[int64]interface{} `json:"preview"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func reddit(subreddit string) []reddit_post {
	_url := "https://www.reddit.com/r/" + subreddit + ".json"
	_req, _err := http.NewRequest("GET", _url, nil)
	_req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36 OPR/78.0.4093.153")
	_req.Header.Add("accept", "application/json;q=0.9")
	_req.Header.Add("accept-encoding", "*")
	_resp, _err := wc.Do(_req)
	if _err != nil {
		// handle it or something
	}
	defer _resp.Body.Close()
	_body, _err := io.ReadAll(_resp.Body)
	var _json reddit_api_response
	json.Unmarshal(_body, &_json)
	_posts := _json.Data.Children
	_returns := []reddit_post{}
	for _, _child := range _posts {
		_post := _child.Data
		fmt.Println(_post.Preview[0].(map[string]string))
		_returns = append(_returns, reddit_post{
			title:       _post.Title,
			author:      _post.Author,
			thumbnail:   _post.Thumbnail,
			image_url:   _post.ImageUrl,
			is_mod_post: (_post.Distinguished == "moderator"),
		})
	}
	return _returns
}

type img_response struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

func default_enpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "afo suck cock")
}

func reddit_endpoint(w http.ResponseWriter, r *http.Request) {
	_mode := r.URL.Query().Get("type")
	_subreddit := "yiff"
	if _mode == "gay" {
		_subreddit = "gfur"
	}
	json.NewEncoder(w).Encode(reddit(_subreddit))
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
	_fn, _id := random_fn("yiff")
	_url := "https://" + r.Host + "/img/yiff/" + _fn
	json.NewEncoder(w).Encode([]img_response{
		{Id: strconv.Itoa(_id), Url: _url},
	})
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL, "FROM", r.Header.Get("x-forwarded-for"), " AS \""+r.Header.Get("user-agent")+"\"")
		auth := r.URL.Query().Get("api_key")
		if auth != "" {
			fmt.Println("WITH AUTH KEY ", auth)
		}
		h.ServeHTTP(w, r)
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
	/* dynamic */
	router.HandleFunc("/fox", fox_endpoint)
	router.HandleFunc("/yiff", yiff_endpoint)
	router.HandleFunc("/reddit", reddit_endpoint)
	/* static */
	router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./imgs/"))))
	router.PathPrefix("/demo/").Handler(http.StripPrefix("/demo/", http.FileServer(http.Dir("./static/demo/"))))
	router.Handle("/", http.FileServer(http.Dir("./static/")))
	/* 404 */
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "afo suck cock")
	})
	/* 405 */
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Method not allowed")
	})
	/* listen */
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(":"+port, middleware(router)))
}
