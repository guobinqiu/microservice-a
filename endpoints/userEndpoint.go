package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"microservice-a/models"
	"microservice-a/registry"
	"microservice-a/services"
	"strconv"
)

//1 获得解码后的请求（自定义请求对象）
//2 调用业务逻辑
//3 返回编码前的响应（自定义响应对象）
func MakeGetUserNameEndpoint(svc services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(models.UserRequest)
		name, _ := svc.GetName(r.Uid)
		return models.UserResponse{Name: name + strconv.Itoa(registry.Port)}, nil
	}
}

func MakeDelUserEndpoint(svc services.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(models.UserRequest)
		if err := svc.DelUser(r.Uid); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
