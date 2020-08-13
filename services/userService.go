package services

import (
	"fmt"
	"microservice-a/myerror"
	"net/http"
)

type UserService struct{}

func (u UserService) GetName(userid int) (string, error) {
	if userid == 100 {
		return "guobin", nil
	}
	return "guest", nil
}

func (u UserService) DelUser(userid int) error {
	if userid == 100 {
		return myerror.New("userid为100的记录不可以删除", http.StatusForbidden)
	}
	fmt.Printf("userid为%d的记录删除成功\n", userid)
	return nil
}
