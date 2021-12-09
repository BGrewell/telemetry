package grpc

import (
	"context"
	api "github.com/BGrewell/tgams/api/go"
	log "github.com/BGrewell/tgams/internal/logging"
	"google.golang.org/grpc"
	"time"
)

type ControlClient struct {
	client api.ControlClient
	conn   *grpc.ClientConn
}

func (c *ControlClient) Connect(address string, timeout int) (err error) {
	log.DebugWithFields(map[string]interface{}{"address": address, "timeout": timeout}, "connecting to control server")
	// TODO: add certs and get rid of 'WithInsecure()'
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	c.conn, err = grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.ErrorWithFields(map[string]interface{}{"error": err}, "failed to connect to control server")
		return err
	}
	log.Debug("connected to control server")

	c.client = api.NewControlClient(c.conn)
	log.Debug("created new grpc client with established connection to control server")
	return nil
}
