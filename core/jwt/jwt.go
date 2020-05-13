package jwt

import (
	"context"
	"fmt"

	proto "github.com/ZeroTechh/VelocityCore/proto/JWTService"
	"github.com/ZeroTechh/VelocityCore/services"
	"github.com/ZeroTechh/VelocityCore/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// New returns new JWT.
func New(log *zap.Logger) *JWT {
	j := JWT{}
	j.init(log)
	return &j
}

// JWT will handle creating and validating of access, refresh and fresh tokens.
type JWT struct {
	client proto.JWTClient
	conn   *grpc.ClientConn
}

// init initializes.
func (j *JWT) init(log *zap.Logger) {
	j.conn = utils.CreateGRPCClient(services.JWTService, log)
	j.client = proto.NewJWTClient(j.conn)
}

// Fresh returns fresh access token.
func (j JWT) Fresh(ctx context.Context, id string) (string, error) {
	r, err := j.client.FreshToken(ctx, &proto.JWTData{UserIdentity: id})
	err = errors.Wrap(err, "Error while creating fresh token")
	fmt.Println(err)
	return r.Token, err
}

// AccessAndRefresh returns access and refresh token
func (j JWT) AccessAndRefresh(
	ctx context.Context, id string, scopes []string) (string, string, error) {

	r, err := j.client.AccessAndRefreshTokens(ctx, &proto.JWTData{
		UserIdentity: id, Scopes: scopes,
	})
	err = errors.Wrap(err, "Error while creating access and refresh token")
	return r.GetAcccessToken(), r.GetRefreshToken(), err
}

// Refresh returns new access and refresh token based on old refresh token.
func (j JWT) Refresh(
	ctx context.Context, token string) (string, string, string, error) {

	r, err := j.client.RefreshTokens(ctx, &proto.Token{Token: token})
	err = errors.Wrap(err, "Error while refreshing token")
	return r.GetAcccessToken(), r.GetRefreshToken(), r.GetMessage(), err
}

// Validate validates a token.
func (j JWT) Validate(
	ctx context.Context,
	token string,
	tokenType proto.TokenType) (proto.Claims, string, error) {

	r, err := j.client.ValidateToken(ctx, &proto.ValidRequest{
		Type: tokenType, Token: token,
	})
	err = errors.Wrap(err, "Error while validating token")
	return *r, r.GetMessage(), err
}

// Disconnect will disconnect from service.
func (j *JWT) Disconnect() {
	j.conn.Close()
}
