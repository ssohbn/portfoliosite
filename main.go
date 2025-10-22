package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// home page should just go to index
	mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request)  {
		http.Redirect(w, r, "/static/html/index.html", http.StatusFound)
	})

	log.Printf("running server on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
