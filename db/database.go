package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

func GetIDAndPassword(db *sql.DB, email string) (int64, string, error) {
	row := db.QueryRowContext(context.TODO(), "select id, password from users where email = $1 limit 1", email)
	var password string
	var id int64
	err := row.Scan(&id, &password)
	return id, password, err
}

func CreateUser(db *sql.DB, email string, password string) (*User, error) {
	row := db.QueryRowContext(context.TODO(), "insert into users(email, password) values ($1,$2) returning id,email", email, password)
	user := &User{}
	err := row.Scan(&user.ID, &user.Email)
	return user, err
}

type Blogs struct {
	ID          int64     `json:"id"`
	Author      string    `json:"author,omitempty"`
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
	rows, err := db.QueryContext(context.TODO(), "select blogs.id, users.email, title,description,created_at from blogs inner join users on blogs.user_id = users.id order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*Blogs{}

	for rows.Next() {
		item := &Blogs{}
		err = rows.Scan(&item.ID, &item.Author, &item.Title, &item.Description, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func SearchBlogs(db *sql.DB, search string) ([]*Blogs, error) {
	rows, err := db.QueryContext(context.TODO(), "select blogs.id, users.email, title,description,created_at from blogs inner join users on blogs.user_id = users.id where ts @@ to_tsquery('simple', $1 || ':*')", search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []*Blogs{}

	for rows.Next() {
		item := &Blogs{}
		err = rows.Scan(&item.ID, &item.Author, &item.Title, &item.Description, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func CreateBlog(db *sql.DB, userID int64, title string, description string) (*Blogs, error) {
	row := db.QueryRowContext(context.TODO(), "insert into blogs(user_id, title, description) values ($1,$2,$3) returning id, title, description, created_at", userID, title, description)
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
