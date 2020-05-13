package user

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	extraProto "github.com/ZeroTechh/VelocityCore/proto/UserExtraService"
	mainProto "github.com/ZeroTechh/VelocityCore/proto/UserMainService"
	metaProto "github.com/ZeroTechh/VelocityCore/proto/UserMetaService"
	"github.com/ZeroTechh/VelocityCore/services"
	"github.com/ZeroTechh/VelocityCore/utils"
	"github.com/pkg/errors"
)

func New(log *zap.Logger) *User {
	u := User{}
	u.init(log)
	return &u
}

// User handles user data.
type User struct {
	mainConn    *grpc.ClientConn
	extraConn   *grpc.ClientConn
	metaConn    *grpc.ClientConn
	mainClient  mainProto.UserMainClient
	extraClient extraProto.UserExtraClient
	metaClient  metaProto.UserMetaClient
}

// init initializes.
func (u *User) init(log *zap.Logger) {
	u.mainConn = utils.CreateGRPCClient(services.UserMainService, log)
	u.extraConn = utils.CreateGRPCClient(services.UserExtraService, log)
	u.metaConn = utils.CreateGRPCClient(services.UserMetaService, log)

	u.mainClient = mainProto.NewUserMainClient(u.mainConn)
	u.extraClient = extraProto.NewUserExtraClient(u.extraConn)
	u.metaClient = metaProto.NewUserMetaClient(u.metaConn)
}

// Auth athenticates user.
func (u User) Auth(
	ctx context.Context, username, email, password string) (bool, string, error) {
	r, err := u.mainClient.Auth(ctx, &mainProto.AuthRequest{
		Username: username, Email: email, Password: password})
	err = errors.Wrap(err, "Error while authenticating user")
	return r.GetValid(), r.GetUserID(), err
}

// validate validates user data.
func (u User) validate(
	ctx context.Context, main mainProto.Data, extra extraProto.Data) (bool, error) {

	mainResp, err := u.mainClient.Validate(ctx, &main)
	if err != nil {
		return false, errors.Wrap(err, "Error while validating main data")
	}

	extraResp, err := u.extraClient.Validate(ctx, &extra)
	err = errors.Wrap(err, "Error while validating extra data")

	valid := mainResp.GetValid() && extraResp.GetValid()
	return valid, err
}

// Register adds user into database.
func (u User) Register(
	ctx context.Context, main mainProto.Data, extra extraProto.Data) (string, string, error) {
	valid, err := u.validate(ctx, main, extra)
	if err != nil || !valid {
		return "", "INVALID DATA", err
	}

	mainResp, err := u.mainClient.Add(ctx, &main)
	if err != nil || mainResp.GetMessage() != "" {
		err = errors.Wrap(err, "Error while adding main data into db")
		return "", mainResp.GetMessage(), err
	}
	id := mainResp.GetUserID()

	extra.UserID = id
	extraResp, err := u.extraClient.Add(ctx, &extra)
	if err != nil || extraResp.GetMessage() != "" {
		err = errors.Wrap(err, "Error while adding extra data into db")
		return "", mainResp.GetMessage(), err
	}

	metaResp, err := u.metaClient.Add(ctx, &metaProto.Identifier{UserID: id})
	err = errors.Wrap(err, "Error while adding extra data into db")
	return id, metaResp.GetMessage(), err
}

// Activate activates user account
func (u User) Activate(ctx context.Context, id string) error {
	_, err := u.metaClient.Activate(ctx, &metaProto.Identifier{UserID: id})
	return errors.Wrap(err, "Error while activating user account")
}

// Diconnect disconnects.
func (u *User) Disconnect() {
	u.mainConn.Close()
	u.extraConn.Close()
	u.metaConn.Close()
}
