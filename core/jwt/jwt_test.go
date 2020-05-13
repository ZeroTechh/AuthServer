package jwt

import (
	"context"
	"testing"

	proto "github.com/ZeroTechh/VelocityCore/proto/JWTService"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestJWT(t *testing.T) {
	assert := assert.New(t)
	ctx := context.TODO()
	id := "id"
	scopes := []string{"read", "write"}
	cfg := zap.NewDevelopmentConfig()
	logger, _ := cfg.Build()
	j := New(logger)

	fresh, err := j.Fresh(ctx, id)
	assert.NoError(err)
	assert.NotZero(fresh)

	access, refresh, err := j.AccessAndRefresh(ctx, id, scopes)
	assert.NotZero(access)
	assert.NotZero(refresh)
	assert.NoError(err)

	access, refresh, msg, err := j.Refresh(ctx, refresh)
	assert.NotZero(access)
	assert.NotZero(refresh)
	assert.Zero(msg)
	assert.NoError(err)

	claims, msg, err := j.Validate(ctx, access, proto.TokenType_ACCESS)
	assert.NoError(err)
	assert.Zero(msg)
	assert.Equal(scopes, claims.Scopes)
	assert.Equal(id, claims.UserIdentity)

	j.Disconnect()
}
