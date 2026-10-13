package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Service struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Welcome, SkillHub!",
		})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		var service Service

		defer r.Body.Close()

		err := json.NewDecoder(r.Body).Decode(&service)
		if err != nil {
			http.Error(w, "Error to read JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(service.Name) == "" {
			http.Error(w, "Error on JSON, Name must not be empty", http.StatusBadRequest)
			return
		}

		if service.Price <= 0 {
			http.Error(w, "Error on JSON, Price can not be negative number", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(service)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func main() {
	http.HandleFunc("/welcome", welcomeHandler)
	http.HandleFunc("/services", servicesHandler)

	log.Println("Servidor iniciado na porta 8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
