package model

type Ad struct {
	ID          int     `gorm:"primaryKey" json: "id"`
	Title       string  `gorm:"not null" json: "title"`
	Description string  `json: "description"`
	ImageURL    string  `json: "image_url"`
	Price       float64 `json: "price"`
	AuthorID    int     `gorm:"not null" json: "author_id"`
	Author      User    `gorm:"foreignKey:AuthorID" json: "author"`
}

type User struct {
	ID       uint   `gorm:"primary_key" json: "id"`
	Login    string `gorm:"unique;not null" json: "login"`
	Password string `gorm: "not null" json: "-"`
}
