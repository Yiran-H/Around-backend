package handler

import (
    //标准库的
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath" //+

    //你自己的
	"around/model"
	"around/service"

    //第三方的
	jwt "github.com/form3tech-oss/jwt-go"
    "github.com/gorilla/mux" 
	"github.com/pborman/uuid" //+
)
//返回给前端 前端知道怎么显示  后端存文件怎么存都行 跟后缀没关系
var (
	mediaTypes = map[string]string{
		".jpeg": "image",
		".jpg":  "image",
        ".JPG":  "image",
		".gif":  "image",
		".png":  "image",
		".mov":  "video",
		".mp4":  "video",
		".avi":  "video",
		".flv":  "video",
		".wmv":  "video",
	}
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
    // fmt.Println("Received one post request")
    // decoder := json.NewDecoder(r.Body)
    // var p model.Post
    // if err := decoder.Decode(&p); err != nil {
    //     panic(err)
    // }

    // fmt.Fprintf(w, "Post received: %s\n", p.Message)

    //request body不是json 而是form data 上面代码不能用
    fmt.Println("Received one upload request")

    user := r.Context().Value("user")//return token
    claims := user.(*jwt.Token).Claims//return payload
    username := claims.(jwt.MapClaims)["username"]

    //new object，想要new pointer 前面加& 这个是初始化 ；声明一个pointer用*
	p := model.Post{
		Id:      uuid.New(), //auto create global unique id
		// User:    r.FormValue("user"), //从request body得到的信息
		User:    username.(string), //.() type assertion (java's cast)
		//从token中得到的 因为如果你post其实不需要额外输入username request body不需要username了
		//从request信息的token中读就行了
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("media_file") //header:mata data; file:content
	if err != nil {
        //w: object 最后err msg写在response里面-w
		http.Error(w, "Media file is not available / Failed to read media_file from request", http.StatusBadRequest) //500 internal err
		fmt.Printf("Media file is not available %v\n", err)
		return
	}

	suffix := filepath.Ext(header.Filename) //根据名字返回文件后缀
	if t, ok := mediaTypes[suffix]; ok {
		p.Type = t
	} else {
		p.Type = "unknown"
	}

	err = service.SavePost(&p, file)
	if err != nil {
		http.Error(w, "Failed to save post to backend", http.StatusInternalServerError)
		fmt.Printf("Failed to save post to backend %v\n", err)
		return
	}

	// fmt.Fprintf(w, "Post received: %s\n", p.Message)
	fmt.Println("Post is saved successfully.")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one request for search")
	w.Header().Set("Content-Type", "application/json")

	user := r.URL.Query().Get("user")
	keywords := r.URL.Query().Get("keywords")

	var posts []model.Post
	var err error
	if user != "" {
		posts, err = service.SearchPostsByUser(user)
	} else {
		posts, err = service.SearchPostsByKeywords(keywords)
	}

	if err != nil {
		http.Error(w, "Failed to read post from backend", http.StatusInternalServerError)
		fmt.Printf("Failed to read post from backend %v.\n", err)
		return
	}

	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
		return
	}
	w.Write(js)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for delete")

    user := r.Context().Value("user") 
    claims := user.(*jwt.Token).Claims 
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

    if err := service.DeletePost(id, username); err != nil {
        http.Error(w, "Failed to delete post from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete post from backend %v\n", err)
        return
    }
    fmt.Println("Post is deleted successfully")
}