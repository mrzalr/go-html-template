package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"
	"os"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/parse-num", parseNum)

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	http.ListenAndServe(port, nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		indexPath := path.Join("template", "index.html")
		tmpl, err := template.ParseFiles(indexPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func parseNum(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		payload := map[string]string{}

		dc := json.NewDecoder(r.Body)
		err := dc.Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		numbers := strings.Split(payload["numbers"], "\n")
		count := len(numbers)
		invalid := 0

		result := []string{}
		for _, num := range numbers {
			if len(num) == 0 {
				invalid++
				continue
			}

			_, err := strconv.ParseInt(num, 10, 64)
			if err != nil {
				invalid++
				continue
			}

			result = append(result, fmt.Sprintf("62%s", num))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"numbers": strings.Join(result, "\n"),
			"count":   count,
			"invalid": invalid,
		})

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
