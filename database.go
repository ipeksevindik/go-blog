package main

import (
	"encoding/json"
	"errors"
	"fmt"

	s "strings"
	"time"
)

type Blog struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (blog *Blog) ToJson() ([]byte, error) {
	data, err := json.Marshal(blog)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (blog *Blog) FromJson(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), blog)
	if err != nil {
		return err
	}
	return nil
}

type Database struct {
	Blogs []*Blog
}

// 1-) dosyayı string olarak oku
// 2-) eline geçen stringi "/n" karakteri ile split et
// 3-) split sonucu oluşan []string dizisinin üzerinde for loop yap
// 4-) her bir ayrı string'i fromJson ile Blog structuna çevir
// 5-) eline geçen blog structını db.Blogs = append(db.Blogs, [eline geçen blog]) şeklinde ekle

func (db *Database) GetAllBlogsAsString() (string, error) {
	bloglarStr := []string{}
	for _, blog := range db.Blogs {
		blogBytes, err := blog.ToJson()
		if err != nil {
			return "", nil
		}
		bloglarStr = append(bloglarStr, string(blogBytes))
	}
	return s.Join(bloglarStr, "\n"), nil
}

func (db *Database) AddAllBlogs(content string) error {
	arr := s.Split(content, "\n")

	for _, blogStr := range arr {
		if blogStr == "" {
			continue
		}
		blog := &Blog{}
		err := blog.FromJson(blogStr)
		if err != nil {
			return err
		} else {
			db.Blogs = append(db.Blogs, blog)
		}

	}

	return nil
}

func (db *Database) GetBlogs() []*Blog {
	return db.Blogs
}

func (db *Database) SearchBlogs(search string) []*Blog {
	result := []*Blog{}

	for _, blog := range db.Blogs {
		if s.Contains(blog.Title, search) || s.Contains(blog.Description, search) {
			result = append(result, blog)
		}
	}

	return result
}

func (db *Database) CreateBlog(title string, description string) *Blog {
	blog := &Blog{
		ID:          len(db.Blogs),
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
	}

	db.Blogs = append(db.Blogs, blog)

	return blog
}

func (db *Database) DeleteBlog(id int) error {
	idx := -1
	for i, blog := range db.Blogs {
		if blog.ID == id {
			idx = i
			break
		}
	}

	if idx != -1 {
		db.Blogs = append(db.Blogs[:idx], db.Blogs[idx+1:]...)
		return nil
	} else {
		return errors.New("böyle bir blog yok")
	}
}

func (db *Database) UpdateBlog(id int, title string, description string) *Blog {

	idx := -1
	for i, blog := range db.Blogs {
		if blog.ID == id {
			idx = i
			break
		}
	}

	if idx != -1 {
		db.Blogs[idx].Title = title
		db.Blogs[idx].Description = description
		return db.Blogs[idx]
	} else {
		return nil
	}

}

func (db *Database) Print() {
	for _, blog := range db.Blogs {
		fmt.Println(blog)
	}
}
