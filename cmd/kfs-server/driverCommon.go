package main

import (
	"errors"
	"fmt"
)

type TokenErrResp struct {
	ErrorDescription string `json:"error_description"`
	ErrorMsg         string `json:"error"`
}

var EmptyToken = errors.New("empty token")

func (e *TokenErrResp) Error() string {
	return fmt.Sprint(e.ErrorMsg, " : ", e.ErrorDescription)
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
