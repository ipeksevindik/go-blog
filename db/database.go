package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type Blogs struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (blog *Blogs) ToJson() ([]byte, error) {
	data, err := json.Marshal(blog)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (blog *Blogs) FromJson(jsonString string) error {
	err := json.Unmarshal([]byte(jsonString), blog)
	if err != nil {
		return err
	}
	return nil
}

func GetBlogs(db *sql.DB) ([]*Blogs, error) {
	rows, err := db.QueryContext(context.TODO(), "select id, title,description,created_at from blogs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*Blogs{}

	for rows.Next() {
		item := &Blogs{}
		err = rows.Scan(&item.ID, &item.Title, &item.Description, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func SearchBlogs(db *sql.DB, search string) ([]*Blogs, error) {
	rows, err := db.QueryContext(context.TODO(), "select id, title,description,created_at from blogs where ts @@ to_tsquery('simple', $1 || ':*')", search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*Blogs{}

	for rows.Next() {
		item := &Blogs{}
		err = rows.Scan(&item.ID, &item.Title, &item.Description, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func CreateBlog(db *sql.DB, title string, description string) (*Blogs, error) {
	row := db.QueryRowContext(context.TODO(), "insert into blogs(title, description) values ($1,$2) returning id, title, description, created_at", title, description)
	blog := &Blogs{}
	err := row.Scan(&blog.ID, &blog.Title, &blog.Description, &blog.CreatedAt)
	return blog, err
}

func DeleteBlog(db *sql.DB, blogID int) (int64, error) {
	row := db.QueryRowContext(context.TODO(), "delete from blogs where id = $1 returning id", blogID)
	var id int64
	err := row.Scan(&id)
	return id, err
}

func UpdateBlog(db *sql.DB, blogID int, title string, description string) (int64, error) {
	row := db.QueryRowContext(context.TODO(), "update blogs set title = $1, description =$2 where id = $3 returning id", title, description, blogID)
	var id int64
	err := row.Scan(&id)
	return id, err
}
