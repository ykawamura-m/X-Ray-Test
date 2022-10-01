package model

const (
	MYSQL  = 1
	DYNAMO = 2
)

// ユーザー情報
type User struct {
	ID    string `form:"id"`
	Name  string `form:"name"`
	Email string `form:"email"`
	Tel   string `form:"tel"`
	DB    int    `form:"db"`
}
