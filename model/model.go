package model

type Post struct {
    Id      string `json:"id"`
    User    string `json:"user"`
    Message string `json:"message"`
    Url     string `json:"url"`
    Type    string `json:"type"`
}

type User struct {
	Username string `json:"username"`//``:不需要转义 因为string也有双引号 可能需要转义
	Password string `json:"password"`
	Age      int64  `json:"age"`
	Gender   string `json:"gender"`
}