package code

import (
	"context"

	emailProto "github.com/ZeroTechh/VelocityCore/proto/EmailService"
	codeProto "github.com/ZeroTechh/VelocityCore/proto/VerificationCodeService"
	"github.com/ZeroTechh/VelocityCore/services"
	"github.com/ZeroTechh/VelocityCore/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// New returns new Code.
func New(log *zap.Logger) *Code {
	c := Code{}
	c.init(log)
	return &c
}

// Code handles creation of code and sending it to email
type Code struct {
	codeClient  codeProto.VerificationCodeClient
	codeConn    *grpc.ClientConn
	emailClient emailProto.EmailClient
	emailConn   *grpc.ClientConn
}

// init initializes.
func (c *Code) init(log *zap.Logger) {
	c.codeConn = utils.CreateGRPCClient(services.VerificationCodeService, log)
	c.codeClient = codeProto.NewVerificationCodeClient(c.codeConn)
	c.emailConn = utils.CreateGRPCClient(services.EmailVerificationSrv, log)
	c.emailClient = emailProto.NewEmailClient(c.emailConn)
}

// CreateAndSend creates a code and sends it to the email
func (c Code) CreateAndSend(ctx context.Context, id, email string) (string, error) {
	codeResp, err := c.codeClient.Create(ctx, &codeProto.UserData{UserID: id})
	if err != nil {
		return "", errors.Wrap(err, "Error while creating verification code")
	}

	_, err = c.emailClient.SendSimpleEmail(ctx, &emailProto.EmailData{Email: email})
	err = errors.Wrap(err, "Error while sending email")
	return codeResp.GetToken(), err
}

// Verify verifies code and returns email
func (c Code) Verify(ctx context.Context, token string) (bool, string, error) {
	r, err := c.codeClient.Validate(ctx, &codeProto.TokenData{Token: token})
	err = errors.Wrap(err, "Error while validating")
	return r.GetValid(), r.GetUserID(), err
}

// Disconnect disconnects client
func (c *Code) Disconnect() {
	c.codeConn.Close()
	c.emailConn.Close()
}
