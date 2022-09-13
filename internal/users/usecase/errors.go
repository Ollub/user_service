package usecase

import "errors"

var UserExistsError = errors.New("user already exists")
var UserNotFoundError = errors.New("user not found")
var BadPasswordError = errors.New("passwords dont match")
