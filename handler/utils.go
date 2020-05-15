package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	extraProto "github.com/ZeroTechh/VelocityCore/proto/UserExtraService"
	mainProto "github.com/ZeroTechh/VelocityCore/proto/UserMainService"
	"github.com/pkg/errors"

	"github.com/ZeroTechh/blaze"
	"go.uber.org/zap"
)

var keys = config.Map("keys")

func fromUrl(r *http.Request, key string) string {
	data, ok := r.URL.Query()[keys.Str(key)]
	if !ok {
		return ""
	}
	return data[0]
}

func mainData(r *http.Request) (username, email, password, scopes string) {
	username = fromUrl(r, "username")
	email = fromUrl(r, "email")
	password = fromUrl(r, "password")
	scopes = fromUrl(r, "scopes")
	return
}

func extraData(r *http.Request) (firstName, lastName, gender, birthdayUTC string) {
	firstName = fromUrl(r, "firstName")
	lastName = fromUrl(r, "lastName")
	gender = fromUrl(r, "gender")
	birthdayUTC = fromUrl(r, "birthdayUTC")
	return
}

func registerData(r *http.Request) (mainProto.Data, extraProto.Data, []string, string, error) {
	username, email, pwd, scopes := mainData(r)
	first, last, gender, utc := extraData(r)
	for _, d := range []string{username, email, pwd, scopes, first, last, gender, utc} {
		if d == "" {
			return mainProto.Data{}, extraProto.Data{}, nil, messages.Str("dataNotProvided"), nil
		}
	}

	intUTC, err := strconv.ParseInt(utc, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Error while converting birthday to int")
		return mainProto.Data{}, extraProto.Data{}, nil, "", err
	}

	main := mainProto.Data{
		Username: username,
		Email:    email,
		Password: pwd,
	}
	extra := extraProto.Data{
		FirstName:   first,
		LastName:    last,
		Gender:      gender,
		BirthdayUTC: intUTC,
	}

	return main, extra, strings.Split(scopes, ","), "", nil
}

func authData(r *http.Request) (string, string, string, []string, string) {
	username, email, pwd, scopes := mainData(r)
	if (username == "" && email == "") || pwd == "" || scopes == "" {
		return "", "", "", nil, messages.Str("dataNotProvided")
	}
	return username, email, pwd, strings.Split(scopes, ","), ""
}

func verificationData(r *http.Request) (string, string) {
	code := fromUrl(r, "emailVerificationCode")
	if code == "" {
		return "", messages.Str("dataNotProvided")
	}
	return code, ""
}

func respond(w http.ResponseWriter, data map[string]string, msg string, err error, log *blaze.FuncLog) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		_ = log.Error(err)
		data = map[string]string{
			keys.Str("error"): messages.Str("internalServerError"),
		}
	}
	if msg != "" {
		data = map[string]string{
			keys.Str("message"): msg,
		}
	}
	log.Completed(zap.Any("Data", data))
	_ = json.NewEncoder(w).Encode(data)
}
