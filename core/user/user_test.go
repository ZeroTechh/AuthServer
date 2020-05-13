package user

import (
	"context"
	"math/rand"
	"testing"
	"time"

	extraProto "github.com/ZeroTechh/VelocityCore/proto/UserExtraService"
	mainProto "github.com/ZeroTechh/VelocityCore/proto/UserMainService"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

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

func TestUser(t *testing.T) {
	assert := assert.New(t)
	cfg := zap.NewDevelopmentConfig()
	logger, _ := cfg.Build()
	u := New(logger)
	ctx := context.TODO()

	// Testing Register.
	main, extra := mock()
	id, msg, err := u.Register(ctx, main, extra)
	assert.NotZero(id)
	assert.Zero(msg)
	assert.NoError(err)

	// Testing Register returns message for invalid data
	_, msg, err = u.Register(ctx, mainProto.Data{}, extra)
	assert.NoError(err)
	assert.NotZero(msg)

	assert.NoError(u.Activate(ctx, id))
	u.Disconnect()
}
