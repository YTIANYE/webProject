package handler

import (
	"context"

	"github.com/micro/go-micro/util/log"

	getCaptcha "webProject/service/getCaptcha/proto/getCaptcha"
)

type GetCaptcha struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetCaptcha) Call(ctx context.Context, req *getCaptcha.Request, rsp *getCaptcha.Response) error {
	log.Log("Received GetCaptcha.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *GetCaptcha) Stream(ctx context.Context, req *getCaptcha.StreamingRequest, stream getCaptcha.GetCaptcha_StreamStream) error {
	log.Logf("Received GetCaptcha.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&getCaptcha.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *GetCaptcha) PingPong(ctx context.Context, stream getCaptcha.GetCaptcha_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&getCaptcha.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
