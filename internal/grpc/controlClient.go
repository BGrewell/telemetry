package grpc

import (
	"context"
	"github.com/BGrewell/go-conversions"
	api "github.com/BGrewell/tgams/api/go"
	log "github.com/BGrewell/tgams/internal/logging"
	"github.com/BGrewell/tgams/internal/timesync"
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

	//TODO: TEMP
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	response, err := c.client.Ping(ctx, &api.PingRequest{Id: 1})
	if err != nil {
		state, failed := UnpackGrpcError(err)
		if !failed {
			log.ErrorWithFields(map[string]interface{}{
				"code":    state.Code().String(),
				"message": state.Message(),
				"details": state.Details(),
			}, "failed to ping control server")
		} else {
			log.ErrorWithFields(map[string]interface{}{"error": err}, "failed to ping control server")
		}
	}

	log.DebugWithFields(map[string]interface{}{"response": response.Status.String()}, "pinged control server")

	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	ts := &api.TimeSyncRequest{
		Id: 1,
		T1: time.Now().UnixNano(),
	}
	tsr, err := c.client.TimeSync(ctx, ts)
	tsr.T4 = time.Now().UnixNano()

	log.DebugWithFields(map[string]interface{}{"t1": tsr.T1, "t2": tsr.T2, "t3": tsr.T3, "t4": tsr.T4}, "synced time with control server")
	delay := conversions.ConvertNanosecondsToStringTime(timesync.CalcDelay(tsr))
	offset := conversions.ConvertNanosecondsToStringTime(timesync.CalcOffset(tsr))
	log.DebugWithFields(map[string]interface{}{"delay": delay, "offset": offset}, "calculated time sync delay and offset")
	return nil
}
