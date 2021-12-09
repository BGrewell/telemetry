package grpc

import (
	"context"
	"fmt"
	api "github.com/BGrewell/tgams/api/go"
	log "github.com/BGrewell/tgams/internal/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"sync"
	"time"
)

var (
	serverOnce     sync.Once
	serverInstance *controlServer
)

func GetControlServer(listenAddr string, listenPort int) *controlServer {
	serverOnce.Do(func() {
		serverInstance = &controlServer{
			listenAddr: listenAddr,
			port:       listenPort,
		}
	})
	return serverInstance
}

type controlServer struct {
	api.UnimplementedControlServer
	listenAddr string
	port       int
	requestId  uint64
	responseId uint64
	gs         *grpc.Server
}

// getRequestId returns the next request id
func (s *controlServer) getRequestId() uint64 {
	s.requestId++
	return s.requestId
}

// getResponseId returns the next response id
func (s *controlServer) getResponseId() uint64 {
	s.responseId++
	return s.responseId
}

// ServeAsync starts the grpc api server in a goroutine
func (s *controlServer) ServeAsync() {
	go s.Serve()
}

// Serve starts the grpc api server
func (s *controlServer) Serve() {
	log.DebugWithFields(
		map[string]interface{}{
			"listenAddr": s.listenAddr,
			"port":       s.port,
		}, "setting up grpc listener")

	// Setup listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.listenAddr, s.port))
	if err != nil {
		log.FatalWithFields(
			map[string]interface{}{
				"listenAddr": s.listenAddr,
				"port":       s.port,
				"err":        err,
			}, "failed to start grpc listener")
		return
	}

	// Setup and register grpc server
	s.gs = grpc.NewServer()
	log.Debug("created new grpc server")

	api.RegisterControlServer(s.gs, s)
	log.Debug("registered ControlServer to handle grpc calls")

	// Start grpc server
	log.Debug("servicing grpc endpoint")
	err = s.gs.Serve(listener)
	if err != nil {
		log.FatalWithFields(
			map[string]interface{}{
				"listenAddr": s.listenAddr,
				"port":       s.port,
				"err":        err,
			}, "grpc server failed")
		return
	}
}

func (s *controlServer) Shutdown() {
	log.Debug("shutting down grpc server")
	s.gs.GracefulStop()
	log.Debug("grpc server shut down")
}

func (s *controlServer) Ping(ctx context.Context, req *api.PingRequest) (response *api.PingResponse, err error) {
	log.Debug("received grpc ping request")
	return &api.PingResponse{
		RequestId: req.Id,
		Id:        s.getResponseId(),
		Status:    api.PingResponse_OK,
	}, nil
}

func (s *controlServer) TimeSync(ctx context.Context, req *api.TimeSyncRequest) (response *api.TimeSyncResponse, err error) {
	//TODO: For accuracy this should be done on a socket where we can request hardware timestamps
	log.Debug("received grpc timesync request")
	t2 := time.Now().UnixNano()
	return &api.TimeSyncResponse{
		RequestId: req.Id,
		Id:        s.getResponseId(),
		T1:        req.T1,
		T2:        t2,
		T3:        time.Now().UnixNano(),
	}, nil
}

func (s *controlServer) StartTelemetry(context.Context, *api.StartTelemetryRequest) (*api.StartTelemetryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTelemetry() not implemented")
}

func (s *controlServer) StopTelemetry(context.Context, *api.StopTelemetryRequest) (*api.StopTelemetryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopTelemetry() not implemented")
}

func (s *controlServer) GetTelemetry(context.Context, *api.GetTelemetryRequest) (*api.GetTelemetryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTelemetry() not implemented")
}
