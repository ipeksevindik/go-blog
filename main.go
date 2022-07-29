package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

var db *Database = &Database{Blogs: []*Blog{}}

func main() {

	str := ReadFile()
	err := db.AddAllBlogs(str)
	if err != nil {
		log.Fatal("Dosya Okunamadı")
	}

	r := chi.NewRouter()

	r.Route("/blog", func(r chi.Router) {
		r.Post("/", AddBlog)
		r.Get("/", GetBlogs)
		r.Delete("/", DeleteBlog)
		r.Put("/", UpdateBlog)
	})

	go http.ListenAndServe(":80", r)

	fmt.Scanln()

	str, err = db.GetAllBlogsAsString()
	if err != nil {
		log.Fatal("Json kaydetme hatası")
	} else {
		WriteFile([]byte(str))
	}
}

func AddBlog(w http.ResponseWriter, r *http.Request) {
	createBlog := &Blog{}
	json.NewDecoder(r.Body).Decode(createBlog)
	blog := db.CreateBlog(createBlog.Title, createBlog.Description)
	WriteFile([]byte{})
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func GetBlogs(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("s")
	w.Header().Set("content-type", "application/json")

	if search == "" {
		json.NewEncoder(w).Encode(db.GetBlogs())
	} else {
		json.NewEncoder(w).Encode(db.SearchBlogs(search))
	}

}

func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Write([]byte("Girdiğiniz id bir sayı değildi"))
		return
	}

	err = db.DeleteBlog(id)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("Blog " + idStr + " başarıyla silindi."))
	}

}

func UpdateBlog(w http.ResponseWriter, r *http.Request) {

	updateBlog := &Blog{}
	json.NewDecoder(r.Body).Decode(updateBlog)
	blog := db.UpdateBlog(updateBlog.ID, updateBlog.Title, updateBlog.Description)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func ReadFile() string {
	content, err := os.ReadFile("bloglar.txt")
	if err != nil {
		log.Fatal(err)
	}
	contentStr := string(content)

	return contentStr
}

func WriteFile(data []byte) {
	err := os.WriteFile("bloglar.txt", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
