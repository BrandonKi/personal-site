package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type templateData struct {
	Content template.HTML
}

var (
	baseTemplate *template.Template
	pageMap      map[string]templateData
)

func main() {
	// service := os.Getenv("K_SERVICE")
	// if service == "" {
	// 	service = "???"
	// }

	// revision := os.Getenv("K_REVISION")
	// if revision == "" {
	// 	revision = "???"
	// }
	pageMap = map[string]templateData{
		"/":          {Content: template.HTML(loadFile("html/landing.html"))},
		"/home":      {Content: template.HTML(loadFile("html/home.html"))},
		"/portfolio": {Content: template.HTML(loadFile("html/portfolio.html"))},
		"/resources": {Content: template.HTML(loadFile("html/resources.html"))},
		"/blog":      {Content: template.HTML(loadFile("html/blog.html"))},
		"/contact":   {Content: template.HTML(loadFile("html/contact.html"))},
	}

	baseTemplate, _ = template.ParseFiles("html/index.html")

	http.HandleFunc("/", baseHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func baseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path[:])
	if data, ok := pageMap[r.URL.Path[:]]; ok {
		if err := baseTemplate.Execute(w, data); err != nil {
			msg := http.StatusText(http.StatusInternalServerError)
			log.Printf("template.Execute: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
		}
	} else {
		handleError(w, r, http.StatusNotFound)
	}
}

func handleError(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 page not found")
	}
}

func loadFile(filename string) string {
	body, err := os.ReadFile(filename)
	if err != nil {
		return "ERROR"
	}
	return string(body)
}
