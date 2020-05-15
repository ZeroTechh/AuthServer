package handler

import (
	"context"
	"net/http"

	"github.com/ZeroTechh/blaze"
	"go.uber.org/zap"
)

var success = map[string]string{keys.Str("success"): messages.Str("success")}

// Register adds user into db and returns access and refresh token.
func (h Handler) Register(w http.ResponseWriter, r *http.Request) {
	funcLog := blaze.NewFuncLog("Register", log, zap.Any("Req", r))
	funcLog.Started()

	main, extra, scopes, msg, err := registerData(r)
	if err != nil || msg != "" {
		respond(w, nil, msg, err, funcLog)
		return
	}

	id, msg, err := h.user.Register(context.TODO(), main, extra)
	if err != nil || msg != "" {
		respond(w, nil, msg, err, funcLog)
		return
	}

	_, err = h.code.CreateAndSend(context.TODO(), id, main.Email)
	if err != nil {
		respond(w, nil, "", err, funcLog)
		return
	}

	access, refresh, err := h.jwt.AccessAndRefresh(context.TODO(), id, scopes)
	output := map[string]string{
		keys.Str("accessToken"):  access,
		keys.Str("refreshToken"): refresh,
		keys.Str("userID"):       main.UserID,
	}
	respond(w, output, "", err, funcLog)
}

// Auth authenticates user and returns access and refresh token.
func (h Handler) Auth(w http.ResponseWriter, r *http.Request) {
	funcLog := blaze.NewFuncLog("Auth", log, zap.Any("Req", r))
	funcLog.Started()

	username, email, password, scopes, msg := authData(r)
	if msg != "" {
		respond(w, nil, msg, nil, funcLog)
		return
	}

	valid, id, err := h.user.Auth(context.TODO(), username, email, password)
	if err != nil || !valid {
		respond(w, nil, messages.Str("invalidCredentials"), err, funcLog)
		return
	}

	access, refresh, err := h.jwt.AccessAndRefresh(context.TODO(), id, scopes)
	output := map[string]string{
		keys.Str("accessToken"):  access,
		keys.Str("refreshToken"): refresh,
	}
	respond(w, output, "", err, funcLog)
}

// Verify verifies code and activates account.
func (h Handler) Verify(w http.ResponseWriter, r *http.Request) {
	funcLog := blaze.NewFuncLog("Auth", log, zap.Any("Req", r))
	funcLog.Started()

	code, msg := verificationData(r)
	if msg != "" {
		respond(w, nil, msg, nil, funcLog)
		return
	}

	valid, id, err := h.code.Verify(context.TODO(), code)
	if !valid || err != nil {
		respond(w, nil, messages.Str("invalidCode"), err, funcLog)
		return
	}

	err = h.user.Activate(context.TODO(), id)
	respond(w, success, "", err, funcLog)
}
