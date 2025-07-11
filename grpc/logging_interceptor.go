package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"log"
	"net"
	"strings"
)

func UnaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Printf("gRPC %s from %s", info.FullMethod, getRemoteAddress(ctx))
		resp, err := handler(ctx, req)
		return resp, err
	}
}

func getRemoteAddress(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if forwardedFor, ok := md["x-forwarded-for"]; ok && len(forwardedFor) > 0 {
			return strings.TrimSpace(strings.Split(forwardedFor[0], ",")[0])
		}
		if realIp, ok := md["x-real-ip"]; ok && len(realIp) > 0 {
			return realIp[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		if host, _, err := net.SplitHostPort(p.Addr.String()); err == nil {
			return host
		} else {
			return p.Addr.String()
		}
	}

	return "unknown"
}
