package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go-blog/db"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"

	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

var psql *sql.DB
var tokenUretici *jwtauth.JWTAuth
var secretPassword string = "pineapple"

func main() {
	tokenUretici = jwtauth.New("HS256", []byte(secretPassword), nil)

	psql = connectDB()

	r := chi.NewRouter()

	r.Route("/blog", func(r chi.Router) {
		//r.Use(TOKEN_KONTROL_FONSKIYONU)

		r.Post("/", AddBlog)
		r.Get("/", GetBlogs)
		r.Delete("/", DeleteBlog)
		r.Put("/", UpdateBlog)
	})
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", RegisterUser)
		r.Post("/login", LoginUser)
	})

	http.ListenAndServe(":8080", r)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	user := &db.User{}
	json.NewDecoder(r.Body).Decode(user)

	sifrelenmisPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.Write([]byte("Şifreleme hatası: " + err.Error()))
		return
	}

	userResult, err := db.CreateUser(psql, user.Email, string(sifrelenmisPassword))
	if err != nil {
		w.Write([]byte("User oluşturma hatası: " + err.Error()))
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(userResult)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	user := &db.User{}
	json.NewDecoder(r.Body).Decode(user)

	userID, passwordInDatabase, err := db.GetIDAndPassword(psql, user.Email)
	if err != nil {
		w.Write([]byte("Email veya şifre yanlış"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordInDatabase), []byte(user.Password))
	if err != nil {
		w.Write([]byte("Email veya şifre yanlış"))
		return
	}

	veriler := map[string]interface{}{"id": userID, "email": user.Email}
	expiration := time.Now().Add(time.Hour * 24 * 365)
	jwtauth.SetExpiry(veriler, expiration)

	_, tokenStr, err := tokenUretici.Encode(veriler)
	if err != nil {
		w.Write([]byte("Token üretim hatası : " + err.Error()))
		return
	}

	w.Write([]byte(tokenStr))
}

func GetUserID(r *http.Request) (int64, error) {
	var userID float64
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToLower(bearer[0:6]) == "bearer" {
		token, _ := jwtauth.VerifyToken(tokenUretici, bearer[7:])
		veriler, err := token.AsMap(r.Context())
		if err != nil {
			return 0, err
		} else {
			userID = veriler["id"].(float64)
		}
	} else {
		return 0, errors.New("TOKEN YOK")
	}

	return int64(userID), nil
}

func AddBlog(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserID(r)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	//BLOG ÜRETME KISMI AŞAĞIDA
	createBlog := &db.Blogs{}
	json.NewDecoder(r.Body).Decode(createBlog)

	blog, err := db.CreateBlog(psql, userID, createBlog.Title, createBlog.Description)
	if err != nil {
		w.Write([]byte("Blog oluşturma hatası: " + err.Error()))
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(blog)
}

func GetBlogs(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("s")
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
	w.Header().Set("content-type", "application/json")
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
