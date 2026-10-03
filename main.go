package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type bookmark struct {
	name string
	link string
}

var links = make([]bookmark, 0, 5)

var (
	errNotName          = errors.New("NOT_NAME")
	errNotLink          = errors.New("NOT_LINK")
	errInvalidLink      = errors.New("INVALID_LINK")
	errInvalidForm      = errors.New("INVALID_FORM")
	errMethodNotAllowed = errors.New("METHOD_NOT_ALLOWED")
)

func main() {

	http.HandleFunc("/", home)
	http.HandleFunc("/bookmarks", bookmarks)
	http.HandleFunc("/hello", sayHello)
	http.HandleFunc("/add", addBookmark)

	log.Println("Server run on port http://localhost:3030")

	if err := http.ListenAndServe(":3030", nil); err != nil {
		log.Fatal("Error:", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, "Hello, my name is Misha"); err != nil {
		log.Println("Error:", err)
	}
}

func bookmarks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if len(links) == 0 {
			if _, err := fmt.Fprintln(w, "У вас не создано закладок"); err != nil {
				log.Println("Error:", err)
			}
		}
		for index, value := range links {
			if _, err := fmt.Fprintf(w, "%d. %s: %s\n", index+1, value.name, value.link); err != nil {
				log.Println("Error:", err)
			}
		}
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, writeErr := fmt.Fprintln(w, "Error", errInvalidForm); writeErr != nil {
				log.Println("Error:", writeErr)
			}

			log.Println("Error:", err)
			return
		}

		form := r.PostForm
		name := form.Get("name")
		link := form.Get("link")

		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := fmt.Fprintln(w, "Error:", errNotName); err != nil {
				log.Println("Error:", err)
			}
			return
		}
		if link == "" {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := fmt.Fprintln(w, "Error:", errNotLink); err != nil {
				log.Println("Error:", err)
			}
			return
		}
		if _, err := url.ParseRequestURI(link); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := fmt.Fprintln(w, "Error:", errInvalidLink); err != nil {
				log.Println("Error:", err)
			}
			return
		}
		links = append(links, bookmark{name: name, link: link})
		if _, err := fmt.Fprintln(w, "Закладка добавлена"); err != nil {
			log.Println("Error:", err)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := fmt.Fprintln(w, "Error:", errMethodNotAllowed); err != nil {
			log.Println("Error:", err)
		}
	}
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	// можно записать короче: r.URL.Query().Get("name")

	if _, err := fmt.Fprintf(w, "Привет, %s", name); err != nil {
		log.Println("Error:", err)
	}
}

func addBookmark(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	link := query.Get("link")

	if len(name) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintln(w, errNotName); err != nil {
			log.Println("Error: ", err)
		}
		return
	}

	if len(link) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintln(w, errNotLink); err != nil {
			log.Println("Error: ", err)
		}
		return
	}
	if _, err := url.ParseRequestURI(link); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err := fmt.Fprintln(w, errInvalidLink); err != nil {
			log.Println("Error: ", err)
		}
		return
	}

	links = append(links, bookmark{name: name, link: link})
	if _, err := fmt.Fprint(w, "Закладка добавлена"); err != nil {
		log.Println("Error", err)
	}
}
