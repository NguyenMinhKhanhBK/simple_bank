package gapi

import (
	"context"
	"net/http"
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
		"status_code": int(statusCode),
		"status_text": statusCode.String(),
	})

	if err != nil {
		entry.WithError(err).Error("received a gRPC request")
	} else {
		entry.Info("received a gRPC request")
	}

	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rr *ResponseRecorder) WriteHeader(statusCode int) {
	rr.StatusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}
func (rr *ResponseRecorder) Write(body []byte) (int, error) {
	rr.Body = body
	return rr.ResponseWriter.Write(body)
}

func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rr := &ResponseRecorder{ResponseWriter: w, StatusCode: http.StatusOK}
		handler.ServeHTTP(rr, req)

		logEntry := logrus.WithFields(logrus.Fields{
			"protocol":    "HTTP",
			"method":      req.Method,
			"path":        req.RequestURI,
			"duration":    time.Since(startTime),
			"status_code": rr.StatusCode,
			"status_text": http.StatusText(rr.StatusCode),
		})

		if rr.StatusCode != http.StatusOK {
			logEntry.WithField("body", string(rr.Body)).Error()
		} else {
			logEntry.Info("received an HTTP request")
		}
	})

}
