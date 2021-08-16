package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"text/template"

	"github.com/gorilla/mux" // mux = router
)

type Item struct {
	Id int
}

var items map[int]Item

type Items []Item

func (it Items) Len() int {
	return len(it)
}

func (it Items) Swap(i, j int) {
	it[i], it[j] = it[j], it[i]
}

func (it Items) Less(i, j int) bool {
	return it[i].Id < it[j].Id
}

func MakeWebHandler() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/", GetMainHandler).Methods("GET")
	mux.HandleFunc("/items", GetItemListHandler).Methods("GET")

	fileServer := http.FileServer(http.Dir("./frontend/build/static/"))
	mux.Handle("/static/", http.StripPrefix("static/", fileServer))
	items = make(map[int]Item)
	items[1] = Item{}
	items[2] = Item{}

	return mux
}

func GetItemListHandler(w http.ResponseWriter, r *http.Request) {
	list := make(Items, 0)
	for _, item := range items {
		list = append(list, item)
	}

	sort.Sort(list)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func GetMainHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("frontend/build/index.html")
	t.Execute(w, nil)
}

func main() {
	http.ListenAndServe(":3000", MakeWebHandler())
}
