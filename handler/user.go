package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
    "regexp"
    "time"

    "around/model"
    "around/service"

    jwt "github.com/form3tech-oss/jwt-go"//jwt 重命名
)

var mySigningKey = []byte("secret") //密钥 ssh非对称加密 private/public 

func signinHandler(w http.ResponseWriter, r *http.Request) { 
    //如果request pass by copy需要把所有内容都copy 内容很多 就算不修改request性能不好 没有传地址好
    //response为啥不传*
    fmt.Println("Received one signin request")
    w.Header().Set("Content-Type", "text/plain")

    if r.Method == "OPTIONS" {
        return
    }

    //  Get User information from client json解析
    decoder := json.NewDecoder(r.Body)
    var user model.User
    if err := decoder.Decode(&user); err != nil {//&user decode以后 后面还要接着用 要切切实实改user的值 只是传copy 外面var user没有改
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

    exists, err := service.CheckUser(user.Username, user.Password)
    if err != nil {
        http.Error(w, "Failed to read user from Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to read user from Elasticsearch %v\n", err)
        return
    }

    if !exists {
        http.Error(w, "User doesn't exists or wrong password", http.StatusUnauthorized)
        fmt.Printf("User doesn't exists or wrong password\n")
        return
    }

    //claims = payload
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //password加进去很危险
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString(mySigningKey)
    if err != nil {
        http.Error(w, "Failed to generate token", http.StatusInternalServerError)
        fmt.Printf("Failed to generate token %v\n", err)
        return
    }

    w.Write([]byte(tokenString))
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one signup request")
    // w.Header().Set("Content-Type", "text/plain") //response body有东西 但在这里response body啥都没返回

    decoder := json.NewDecoder(r.Body)
    var user model.User
    if err := decoder.Decode(&user); err != nil {
        http.Error(w, "Cannot decode user data from client", http.StatusBadRequest)
        fmt.Printf("Cannot decode user data from client %v\n", err)
        return
    }

    //必须是小写字母或者数字 在前端输入后就可以check 不对就不能点button 不用等前端发过来再判断
    if user.Username == "" || user.Password == "" || regexp.MustCompile(`^[a-z0-9]$`).MatchString(user.Username) {
        http.Error(w, "Invalid username or password", http.StatusBadRequest)
        fmt.Printf("Invalid username or password\n")
        return
    }

    success, err := service.AddUser(&user)
    if err != nil {
        http.Error(w, "Failed to save user to Elasticsearch", http.StatusInternalServerError)
        fmt.Printf("Failed to save user to Elasticsearch %v\n", err)
        return
    }

    if !success {
        http.Error(w, "User already exists", http.StatusBadRequest)
        fmt.Println("User already exists")
        return
    }
    fmt.Printf("User added successfully: %s.\n", user.Username)
}