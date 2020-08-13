package transports

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"microservice-a/models"
	"net/http"
	"strconv"
)

//收到客户端请求进行解码(http请求转自定义请求对象)
func DecodeUserRequest(c context.Context, r *http.Request) (interface{}, error) {
	//if r.URL.Query().Get("uid") != "" {
	//	uid, _ := strconv.Atoi(r.URL.Query().Get("uid"))
	//	return models.UserRequest{Uid: uid}, nil
	//}

	params := mux.Vars(r)
	if v, ok := params["uid"]; ok {
		uid, _ := strconv.Atoi(v)
		return models.UserRequest{Uid: uid}, nil
	}

	return nil, errors.New("参数错误")
}

//编码响应发送给客户端(对象转json)
func EncodeGetUserNameResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func EncodeDelUserResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	if _, ok := response.(models.Error); ok {
		w.Header().Set("content-type", "application/json")
		return json.NewEncoder(w).Encode(response)
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
