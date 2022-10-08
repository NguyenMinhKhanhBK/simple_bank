package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	grpcGatewayClientIPHeader  = "x-forwarded-for"
	grpcUserAgentHeader        = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (m *Metadata) GetUserAgent() string {
	if m == nil {
		return ""
	}
	return m.UserAgent
}

func (m *Metadata) GetClientIP() string {
	if m == nil {
		return ""
	}
	return m.ClientIP
}

func (s *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return mtdt
	}

	if gwUserAgents := md.Get(grpcGatewayUserAgentHeader); len(gwUserAgents) > 0 {
		mtdt.UserAgent = gwUserAgents[0]
	}

	if userAgents := md.Get(grpcUserAgentHeader); len(userAgents) > 0 {
		mtdt.UserAgent = userAgents[0]
	}

	if clientIPs := md.Get(grpcGatewayClientIPHeader); len(clientIPs) > 0 {
		mtdt.ClientIP = clientIPs[0]
	}

	return mtdt
}
