package gapi

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()

	result, err := handler(ctx, req)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	entry := logrus.WithFields(logrus.Fields{
		"protocol":    "gRPC",
		"method":      info.FullMethod,
		"duration":    time.Since(startTime),
		"status_code": statusCode,
	})

	if err != nil {
		entry.WithError(err).Error("received a gRPC request")
	} else {
		entry.Info("received a gRPC request")
	}

	return result, err
}
