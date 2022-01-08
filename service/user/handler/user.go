package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	user "webProject/service/user/proto/user"
)

type User struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *User) SendSms(ctx context.Context, req *user.Request, rsp *user.Response) error {
	log.Log("Received User.SendSms request")

	return nil
}
