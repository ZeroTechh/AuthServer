package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	extraProto "github.com/ZeroTechh/VelocityCore/proto/UserExtraService"
	mainProto "github.com/ZeroTechh/VelocityCore/proto/UserMainService"
	"github.com/stretchr/testify/assert"
)

var a *assert.Assertions

func randStr(length int) string {
	charset := "1234567890abcdefghijklmnopqrstuvwxyz"
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// mock returns mock user extra data for testing
func mock() (mainProto.Data, extraProto.Data) {
	randomStr := randStr(10)
	main := mainProto.Data{
		Username: randomStr,
		Email:    randomStr + "@gmail.com",
		Password: randomStr,
	}
	extra := extraProto.Data{
		UserID:      randomStr,
		FirstName:   randomStr,
		LastName:    randomStr,
		Gender:      "male",
		BirthdayUTC: int64(864466669),
	}
	return main, extra
}

func getResponse(data map[string]string, function func(w http.ResponseWriter, r *http.Request)) map[string]string {
	req, err := http.NewRequest("GET", "/auth", nil)
	a.NoError(err)

	q := req.URL.Query()
	for key, value := range data {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(function)
	handler.ServeHTTP(rr, req)

	var output map[string]string
	a.NoError(json.Unmarshal(rr.Body.Bytes(), &output))
	return output
}

func TestHandler(t *testing.T) {
	a = assert.New(t)
	h := New()

	// Testing Register.
	main, extra := mock()
	registerReq := map[string]string{
		keys.Str("username"):    main.Username,
		keys.Str("email"):       main.Email,
		keys.Str("password"):    main.Password,
		keys.Str("firstName"):   extra.FirstName,
		keys.Str("lastName"):    extra.LastName,
		keys.Str("gender"):      extra.Gender,
		keys.Str("birthdayUTC"): strconv.FormatInt(extra.BirthdayUTC, 10),
		keys.Str("scopes"):      "read",
	}

	registerResp := getResponse(registerReq, h.Register)
	a.NotZero(registerResp[keys.Str("accessToken")])
	a.NotZero(registerResp[keys.Str("refreshToken")])
	id := registerResp[keys.Str("userID")]

	// Testing Register returns message for invalid data.
	registerResp = getResponse(nil, h.Register)
	a.NotZero(registerResp[keys.Str("message")])

	// Testing Register returns message for already used username.
	registerResp = getResponse(registerReq, h.Register)
	a.NotZero(registerResp[keys.Str("message")])

	// Testing Auth.
	authReq := map[string]string{
		keys.Str("username"): main.Username,
		keys.Str("password"): main.Password,
		keys.Str("scopes"):   "read",
	}

	authResp := getResponse(authReq, h.Auth)
	a.NotZero(authResp[keys.Str("accessToken")])
	a.NotZero(authResp[keys.Str("refreshToken")])

	// Testing Auth returns message for invalid data.
	authResp = getResponse(nil, h.Auth)
	a.NotZero(authResp[keys.Str("message")])

	// Testing Auth returns message for invalid credentials
	authReq[keys.Str("username")] = "invalid username"
	authResp = getResponse(authReq, h.Auth)
	a.NotZero(authResp[keys.Str("message")])

	// Testing Verify.
	code, _ := h.code.CreateAndSend(context.TODO(), id, main.Email)
	verifyReq := map[string]string{
		keys.Str("emailVerificationCode"): code,
	}

	verifyResp := getResponse(verifyReq, h.Verify)
	fmt.Println(verifyResp)
	a.NotZero(verifyResp[keys.Str("success")])

	// Testing Verify returns message for invalid token.
	verifyReq[keys.Str("emailVerificationCode")] = "invalid token"
	verifyResp = getResponse(verifyReq, h.Verify)
	a.NotZero(verifyResp[keys.Str("message")])

	// Testing Verify returns message for invalid data.
	verifyResp = getResponse(nil, h.Verify)
	a.NotZero(verifyResp[keys.Str("message")])
}
