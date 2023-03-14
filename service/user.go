package service

import (
	"fmt"
	"reflect"

	"around/backend"
	"around/constants"
	"around/model"

	"github.com/olivere/elastic/v7"
)

/*
func helper(username string) (*elastic.searchResult, error) {
	query := elastic.NewTermQuery("username", username)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}
	return searchResult, nil 
	//如果想让username一样password不一样的用户也可以注册
	//check调用判断password是不是一样 add调用判断searchReasult是不是空
}
*/

//检查user存不存在 存在就返回login成功 否则失败
func CheckUser(username, password string) (bool, error) {
	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("username", username))
	query.Must(elastic.NewTermQuery("password", password))
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}
	//以上可以替换成：
	//searchResult, err := helper(username)
	//for loop

	//return searchResult.TotalHits() > 0, nil
	//下面的复杂 遍历searchResult
	//reflect.TypeOf(utype) : 对于数据格式的限制 searchResult不只是一个类型
	//因为非关系型数据库没有schema的限制 啥类型的数据都可以存在里面和MYSQL不一样 每条记录的column数是一样的
	//duplicate? username & password 数据库不会check 自己可以写code check
	var utype model.User
	for _, item := range searchResult.Each(reflect.TypeOf(utype)) {
		u := item.(model.User)
		if u.Password == password && u.Username == username {
			fmt.Printf("Login as %s\n", username)
			return true, nil
		}
	}
	return false, nil
}

//signup / bool check duplicate
func AddUser(user *model.User) (bool, error) {
	//check duplicate
	query := elastic.NewTermQuery("username", user.Username)
	searchResult, err := backend.ESBackend.ReadFromES(query, constants.USER_INDEX)
	if err != nil {
		return false, err
	}

	//have duplicate
	if searchResult.TotalHits() > 0 {
		return false, nil
	}
	//以上可以替换成：
	//searchResult, err := helper(username)
	//if totalhits...

	err = backend.ESBackend.SaveToES(user, constants.USER_INDEX, user.Username)
	if err != nil {
		return false, err
	}
	fmt.Printf("User is added: %s\n", user.Username)
	return true, nil
}