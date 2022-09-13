package delivery

import "github.com/Ollub/user_service/internal/users"

type LoginResp struct {
	Token  string `json:"token"`
	UserId uint32 `json:"userId"`
}

type LoginReq struct {
	Email    string
	Password string
}

type ListUsersResp struct {
	Users []*users.User `json:"users"`
}
