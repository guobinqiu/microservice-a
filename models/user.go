package models

type UserRequest struct {
	Uid int `json:"uid"`
}

type UserResponse struct {
	Name string `json:"name"`
}
