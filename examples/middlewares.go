package main

import (
	"net/http"
	"strings"

	"github.com/sriramr98/whiskey"
)

var (
	dummyToken  = "Abcde"
	dummyUserId = "k2380sd"
)

func AuthMiddleware(ctx whiskey.Context) error {
	UnAuthorizedErr := whiskey.NewHttpError(http.StatusUnauthorized, whiskey.BodyTypeJSON)

	authHeader, hasHeader := ctx.GetHeader("Authorization")
	if !hasHeader {
		return UnAuthorizedErr
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return UnAuthorizedErr
	}

	token := strings.Split(authHeader, "Bearer ")[1]
	if token == "" {
		return UnAuthorizedErr
	}

	if token != dummyToken {
		return UnAuthorizedErr
	}

	ctx.Set("userId", dummyUserId)
	return nil
}
