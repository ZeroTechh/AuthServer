package code

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCode(t *testing.T) {
	assert := assert.New(t)
	cfg := zap.NewDevelopmentConfig()
	logger, _ := cfg.Build()
	c := New(logger)
	ctx := context.TODO()
	id := "id"
	email := "sonicroshan122@gmail.com"

	token, err := c.CreateAndSend(ctx, id, email)
	assert.NoError(err)
	assert.NotZero(token)

	valid, e, err := c.Verify(ctx, token)
	assert.Equal(id, e)
	assert.True(valid)
	assert.NoError(err)
}
