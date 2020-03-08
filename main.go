package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gobuffalo/packr/v2"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
)

//Index holds fields displayed on the index.html template
type Index struct {
	AppName          string
	AppInstanceIndex int
	AppInstanceGUID  string
	// Database         string
	Envars    []string
	Services  []Service
	SpaceName string
}

//Service holds the name and label of a service instance
type Service struct {
	Name  string
	Label string
}

func main() {

	index := Index{"Unknown", -1, "Unknown", []string{}, []Service{}, "Unknown"}

	//template := template.Must(template.ParseFiles("./templates/index.html", "./templates/kill.html"))
	var templatesBox = packr.New("Templates", "./templates")
	var staticBox = packr.New("Static", "./static")

	templateIndex, err := templatesBox.FindString("index.html")
	if err != nil {
		log.Fatal(err)
	}
	templateKill, err := templatesBox.FindString("kill.html")
	if err != nil {
		log.Fatal(err)
	}
	t := template.New("")
	t.Parse(templateIndex)
	t.Parse(templateKill)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(staticBox)))

	if cfenv.IsRunningOnCF() {
		appEnv, err := cfenv.Current()
		if err != nil {
			log.Fatal(err)
		}
		if appEnv.Name != "" {
			index.AppName = appEnv.Name
		}
		if appEnv.Index > -1 {
			index.AppInstanceIndex = appEnv.Index
		}
		if appEnv.InstanceID != "" {
			index.AppInstanceGUID = appEnv.InstanceID
		}
		if appEnv.SpaceName != "" {
			index.SpaceName = appEnv.SpaceName
		}
		for _, svcs := range appEnv.Services {
			for _, svc := range svcs {
				index.Services = append(index.Services, Service{svc.Name, svc.Label})
			}
		}
		for _, envar := range os.Environ() {
			if strings.HasPrefix(envar, "TRAINING_") {
				index.Envars = append(index.Envars, envar)
			}
		}

		// config := DBConfig{
		// 	Hostname: os.Getenv("HOSTNAME"),
		// 	Name:     os.Getenv("NAME"),
		// 	Password: os.Getenv("PASSWORD"),
		// 	Username: os.Getenv("USERNAME"),
		// }

	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", index); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// if err := template.ExecuteTemplate(w, "index.html", index); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }
	})

	http.HandleFunc("/kill", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "kill.html", index); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// if err := template.ExecuteTemplate(w, "kill.html", index); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }

	})

	http.HandleFunc("/killInstance", func(w http.ResponseWriter, r *http.Request) {
		os.Exit(1)
	})

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8080"
	}

	fmt.Println(http.ListenAndServe(":"+PORT, nil))
}
