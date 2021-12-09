package grpc

import (
	"google.golang.org/grpc/status"
)

func UnpackGrpcError(err error) (st *status.Status, failed bool) {
	st, ok := status.FromError(err)
	return st, !ok
}
