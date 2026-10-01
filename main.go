package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type bookmark struct {
	name string
	link string
}

var links = make([]bookmark, 0, 5)

var (
	errNotName     = errors.New("NOT_NAME")
	errNotLink     = errors.New("NOT_LINK")
	errInvalidLink = errors.New("INVALID_LINK")
)

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/bookmarks", bookmarks)
	http.HandleFunc("/hello", sayHello)
	http.HandleFunc("/add", addBookmark)

	fmt.Println("Server run on port http://localhost:3030")

	if err := http.ListenAndServe(":3030", nil); err != nil {
		fmt.Println("Error", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, "Hello, my name is Misha"); err != nil {
		fmt.Println("Error:", err)
	}
}

func bookmarks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if len(links) == 0 {
			if _, err := fmt.Fprintln(w, "У вас не создано закладок"); err != nil {
				fmt.Println("Error:", err)
			}
		}
		for index, value := range links {
			if _, err := fmt.Fprintf(w, "%d. %s: %s\n", index+1, value.name, value.link); err != nil {
				fmt.Println("Error:", err)
			}
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			fmt.Println("Error:", err)
		}

		form := r.PostForm
		name := form.Get("name")
		link := form.Get("link")

		if _, err := fmt.Fprintf(w, "Name: %s\nLink: %s\n", name, link); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	// можно записать короче: r.URL.Query().Get("name")

	if _, err := fmt.Fprintf(w, "Привет, %s", name); err != nil {
		fmt.Println("Error:", err)
	}
}

func addBookmark(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	link := query.Get("link")

	if len(name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintln(w, errNotName); err != nil {
			fmt.Println("Error: ", err)
		}
		return
	}

	if len(link) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintln(w, errNotLink); err != nil {
			fmt.Println("Error: ", err)
		}
		return
	}
	if _, err := url.ParseRequestURI(link); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintln(w, errInvalidLink); err != nil {
			fmt.Println("Error: ", err)
		}
		return
	}

	links = append(links, bookmark{name: name, link: link})
	if _, err := fmt.Fprint(w, "Закладка добавлена"); err != nil {
		fmt.Println("Error", err)
	}
}
