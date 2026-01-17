package model

type Comment struct {
	ID     uint
	UserID uint
	User   User
	PostID uint
	Post   Post
}
