package main

import (
	"fmt"
	"html/template"
	"net/http"

)

type Rsvp struct{
	Name, Email, Phone string
	WillAttend bool
}

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates(){
	// Загрузка шаблонов
	templateName := [5]string {"welcome", "form","thanks", "sorry", "list"}
	for index, name := range templateName{
		t, err := template.ParseFiles("templates/layout.html", "templates/" + name + ".html")
		if err != nil {
			panic(err)
		}else{
			templates[name] = t
			fmt.Println("loaded template: ",index, name)
		}
	}
}

func welcomeHandler(w http.ResponseWriter, r *http.Request){
	templates["welcome"].Execute(w, nil)
}

func listHandler(w http.ResponseWriter, r *http.Request){
	templates["list"].Execute(w, responses)
}

type formData struct{
	*Rsvp
	Errors []string
}

func formHandler(w http.ResponseWriter, r * http.Request){
	if r.Method == http.MethodGet {
		templates["form"].Execute(w, formData{
			Rsvp: &Rsvp{},
			Errors: []string{},
		})
	}else if r.Method == http.MethodPost {
		r.ParseForm()
		responseData := Rsvp {
			Name: r.Form["name"][0],
			Email: r.Form["email"][0],
			Phone: r.Form["phone"][0],
			WillAttend: r.Form["willattend"][0] == "true",
		}
		errors:= []string{}
		if responseData.Name == ""{
			errors = append(errors, "Please enter your name")
		}
		if responseData.Email == ""{
			errors = append(errors, "Please enter your email")
		}
		if responseData.Phone == ""{
			errors = append(errors, "Please enter yuor phone")
		}
		if len(errors) > 0 {
			templates["form"].Execute(w, formData{
				Rsvp: &responseData, Errors: errors,
			})
		}else{
			responses = append(responses, &responseData)
			if responseData.WillAttend {
				templates["thanks"].Execute(w, responseData.Name)
			} else {
				templates["sorry"].Execute(w, responseData.Name)
			}
	    }
	}
}

func main()  {
	loadTemplates()
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		fmt.Println(err)
	}
	
}