package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Categories struct {
	ID          int    `json:"id"`
	Nama        string `json:"nama"`
	Description string `json:"description"`
}

var category = []Categories{
	{ID: 1, Nama: "SD", Description: "Sekolah Dasar"},
	{ID: 2, Nama: "SMP", Description: "Sekolah Menegah Pertama"},
	{ID: 3, Nama: "SMA", Description: "Sekolah Menegah Atas"},
}

func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for _, p := range category {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Category Belum ada", http.StatusNotFound)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var updateCategory Categories
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	for i := range category {
		if category[i].ID == id {
			updateCategory.ID = id
			category[i] = updateCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}
	http.Error(w, "Category Belum ada", http.StatusNotFound)


}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for i, p := range category {
		if p.ID == id {
			category = append(category[:i], category[i+1:]...)
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}
	http.Error(w, "Category Belum ada", http.StatusNotFound)
}

func main(){
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoryByID(w,r)		
		} else if r.Method == "PUT" {
			updateCategory(w,r)
		} else if  r.Method == "DELETE" {
			deleteCategory(w,r)
		}
		
	}) 
	
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
			
		} else if r.Method == "POST" {
			var categoryBaru Categories
			err := json.NewDecoder(r.Body).Decode(&categoryBaru)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			categoryBaru.ID = len(category) + 1
			category = append(category, categoryBaru)

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(categoryBaru)

		}
	})

	fmt.Println("server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil{
		fmt.Println("gagal running server")
	}
}