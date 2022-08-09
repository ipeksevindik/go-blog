package main

import (
	"database/sql"
	"encoding/json"
	"go-blog/db"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

var psql *sql.DB

func main() {
	psql = connectDB()

	r := chi.NewRouter()

	r.Route("/blog", func(r chi.Router) {
		r.Post("/", AddBlog)
		r.Get("/", GetBlogs)
		r.Delete("/", DeleteBlog)
		r.Put("/", UpdateBlog)
	})

	http.ListenAndServe(":8080", r)
}

func AddBlog(w http.ResponseWriter, r *http.Request) {
	createBlog := &db.Blogs{}
	json.NewDecoder(r.Body).Decode(createBlog)

	blog, err := db.CreateBlog(psql, createBlog.Title, createBlog.Description)
	if err != nil {
		w.Write([]byte("Blog oluşturma hatası: " + err.Error()))
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func GetBlogs(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("s")
	w.Header().Set("content-type", "application/json")

	var blogs []*db.Blogs
	var err error
	if search == "" {
		blogs, err = db.GetBlogs(psql)
	} else {
		blogs, err = db.SearchBlogs(psql, search)
	}
	if err != nil {
		w.Write([]byte("Blog çekme hatası: " + err.Error()))
		return
	}

	json.NewEncoder(w).Encode(blogs)
}

func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.Write([]byte("Girdiğiniz id bir sayı değildi"))
		return
	}

	_, err = db.DeleteBlog(psql, id)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("Blog " + idStr + " başarıyla silindi."))
	}
}

func UpdateBlog(w http.ResponseWriter, r *http.Request) {
	updateBlog := &db.Blogs{}
	json.NewDecoder(r.Body).Decode(updateBlog)
	blog, err := db.UpdateBlog(psql, int(updateBlog.ID), updateBlog.Title, updateBlog.Description)
	if err != nil {
		log.Println("Güncelleme hatası !")
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func connectDB() *sql.DB {
	conn, err := sql.Open("postgres", "postgresql://ipek:123456@localhost:5432/go-blog?sslmode=disable")
	if err != nil {
		log.Fatalf("error creating database : %v \n", err)
	}
	return conn
}
