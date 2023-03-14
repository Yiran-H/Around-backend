package handler

import (
    "net/http" 

    jwtmiddleware "github.com/auth0/go-jwt-middleware" //前端request发送过来时在trigger handler之前做一些过滤和判断
	//spring security - filter chain
	jwt "github.com/form3tech-oss/jwt-go"
    "github.com/gorilla/mux"   
)

func InitRouter() *mux.Router {
	//判断token是否有效 有效就search post
	//是不是可信设备 机器人 访问频率要求
    jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySigningKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

    router := mux.NewRouter()
    // router.Handle("/upload", http.HandlerFunc(uploadHandler)).Methods("POST")
    // router.Handle("/search", http.HandlerFunc(searchHandler)).Methods("GET")
    router.Handle("/upload", jwtMiddleware.Handler(http.HandlerFunc(uploadHandler))).Methods("POST")
	router.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(searchHandler))).Methods("GET")
	router.Handle("/post/{id}", jwtMiddleware.Handler(http.HandlerFunc(deleteHandler))).Methods("DELETE")

	router.Handle("/signup", http.HandlerFunc(signupHandler)).Methods("POST")
	router.Handle("/signin", http.HandlerFunc(signinHandler)).Methods("POST")
    return router
}