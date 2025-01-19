package internalgrpc

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func loggingInterceptor() grpc.UnaryServerInterceptor {
	if err := os.MkdirAll("logs/", 0o666); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile("logs/grpc_app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		panic(err)
	}

	logger := log.New(logFile, "", 0)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		interface{}, error,
	) {
		start := time.Now()
		result, err := handler(ctx, req)
		end := time.Since(start)

		userAgent := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if mUserAgent, exists := md["user-agent"]; exists {
				userAgent = mUserAgent[0]
			}
		}

		ip := ""
		if p, ok := peer.FromContext(ctx); ok {
			ip = p.Addr.String()
		}

		logger.Printf("%s [%s] %s %s %s %d %s \"%s\"\n",
			ip, start.String(), "POST", info.FullMethod, "GRPC", status.Code(err), end.String(), userAgent)

		return result, err
	}
}
