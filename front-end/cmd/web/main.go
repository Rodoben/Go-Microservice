package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 8081:")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFs embed.FS

func render(w http.ResponseWriter, t string) {

	log.Println("1111111111111111111111111111111111111render")

	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}
	log.Println("11111111122222222222222222222222222222111111111render")
	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	log.Println("TemplateSlice1:", templateSlice)

	for _, x := range partials {
		log.Println("X:", x)
		templateSlice = append(templateSlice, x)
	}
	log.Println("TemplateSlice2:", templateSlice)
	log.Println("111111113333333333333333333311111111111render")
	tmpl, err := template.ParseFS(templateFs, templateSlice...)
	if err != nil {
		log.Println("Error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("11111114444444444444444444411111111111render")
	var data struct {
		BrokerURL string
	}

	//data.BrokerURL = os.Getenv("BROKER_URL")
	data.BrokerURL = "http://localhost:8080"
	log.Println("111111111111555555555555555511111111111render", data.BrokerURL)
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("111111166666666666666666111111111111render11111111111", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("111111166666666666666666111111111111render")
}
