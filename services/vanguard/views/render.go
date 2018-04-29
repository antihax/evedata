package views

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"time"
)

var Templates *template.Template

const LayoutPath string = "templates/layout/layout.html"

var (
	templates            map[string]*template.Template
	templateIncludeFiles []string
)

func init() {
	templates = make(map[string]*template.Template)

	mainTemplate := template.Must(template.New("base").Funcs(template.FuncMap{
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values) == 0 {
				return nil, errors.New("invalid dict call")
			}

			dict := make(map[string]interface{})

			for i := 0; i < len(values); i++ {
				key, isset := values[i].(string)
				if !isset {
					if reflect.TypeOf(values[i]).Kind() == reflect.Map {
						m := values[i].(map[string]interface{})
						for i, v := range m {
							dict[i] = v
						}
					} else {
						return nil, errors.New("dict values must be maps")
					}
				} else {
					i++
					if i == len(values) {
						return nil, errors.New("specify the key for non array values")
					}
					dict[key] = values[i]
				}

			}
			return dict, nil
		},
	}).ParseFiles("templates/layout/layout.html"))

	pageFiles, err := filepath.Glob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	includeFiles, err := filepath.Glob("templates/includes/*.html")
	if err != nil {
		log.Fatal(err)
	}

	templateIncludeFiles = append(includeFiles, "templates/layout/layout.html")

	for _, file := range pageFiles {
		fileName := filepath.Base(file)
		files := append(includeFiles, file)
		templates[fileName], err = mainTemplate.Clone()
		if err != nil {
			log.Fatal(err)
		}
		templates[fileName] = template.Must(templates[fileName].ParseFiles(files...))
	}

}

func renderJSON(w http.ResponseWriter, v interface{}, cacheTime time.Duration) error {
	cache(w, cacheTime)
	return json.NewEncoder(w).Encode(v)
}

func renderTemplate(w http.ResponseWriter, name string, cacheTime time.Duration, data interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, fmt.Sprintf("The template %s does not exist.", name),
			http.StatusInternalServerError)
		return errors.New(fmt.Sprintf("The template %s does not exist.", name))
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	cache(w, cacheTime)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return nil
}

func renderStatic(w http.ResponseWriter, r *http.Request, area string, page page) {
	p := newPage(r, page.Title)
	setCache(w, 60*60*24*31)

	mainTemplate := template.New("base")
	mainTemplate, err := mainTemplate.ParseFiles("templates/layout/layout.html")
	if err != nil {
		log.Fatal(err)
	}
	t, err := mainTemplate.Clone()
	if err != nil {
		log.Fatal(err)
	}

	templates := append(templateIncludeFiles, "templates/"+area+"/"+page.Template, "templates/"+area+"/base.html")
	t = template.Must(
		t.ParseFiles(templates...),
	)

	if err := t.ExecuteTemplate(w, "base", p); err != nil {
		httpErrCode(w, err, http.StatusInternalServerError)
		return
	}
	return
}

func cache(w http.ResponseWriter, cacheTime time.Duration) {
	if cacheTime == 0 {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
	} else {
		w.Header().Set("Cache-Control", "max-age:"+strconv.Itoa(int(cacheTime.Seconds()))+", public")
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		w.Header().Set("Expires", time.Now().UTC().Add(cacheTime).Format(http.TimeFormat))
	}
}
