package service

import (
	"mime/multipart"
	"reflect"

	"around/backend"
	"around/constants"
	"around/model"

	"github.com/olivere/elastic/v7"
)

//select * from post where user = xxx
func SearchPostsByUser(user string) ([]model.Post, error) {
	query := elastic.NewTermQuery("user", user)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

//select * from post where message like "%richard%"
func SearchPostsByKeywords(keywords string) ([]model.Post, error) {
	query := elastic.NewMatchQuery("message", keywords)
	query.Operator("AND")
	if keywords == "" {
		query.ZeroTermsQuery("all")
	}
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.POST_INDEX)
	if err != nil {
		return nil, err
	}
	return getPostFromSearchResult(searchResult), nil
}

func getPostFromSearchResult(searchResult *elastic.SearchResult) []model.Post {
	var ptype model.Post
	var posts []model.Post

	for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
		p := item.(model.Post)
		posts = append(posts, p)
	}
	return posts
}
//从postman/http request body中读出来的file
func SavePost(post *model.Post, file multipart.File) error {
	medialink, err := backend.GCSBackend.SaveToGCS(file, post.Id)
	if err != nil {
		return err
	}
	post.Url = medialink

	err = backend.ESBackend.SaveToES(post, constants.POST_INDEX, post.Id)
	//if (err != nil) { //失败就删 但删除也会失败
	//	err := backend.GCSBackend.DeleteFromGCS(post.Id)
	//}
	/*
	func DeleteFromGCS(post.Id) {...} //网上抄类似
	*/

	//网盘自动检查 遍历所有文件 
	/*
	func Check() { //没必要一失败就回滚 定期检查就行 like GC 如果有错那就再试一次 best effort 对成功率没有要求
		// loop all files in GCS, search ES by url
	}
	*/
	return err
	//存了一个没有人refer的文件 数据库没有东西去refer网盘里的文件 网盘里存了很多无效文件 对user没关系 
	//如果存数据失败 会返回user一个error然后retry 不影响用户使用 会影响后台运行成本 是一个engineering问题
	//两个不同的存储模块没有办法去用mysql的transaction机制解决 需要重新开发library
	//
}

func DeletePost(id string, user string) error {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("id", id))
	query.Must(elastic.NewTermQuery("user", user))

	return backend.ESBackend.DeleteFromES(query, constants.POST_INDEX)
}