/*
 * Copyright (C) 2025 NullSpook
 *
 * This file is part of psirng.
 *
 * psirng is free software: you can redistribute it and/or modify it under the
 * terms of the GNU Affero General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option)
 * any later version.
 *
 * psirng is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
 * FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License for
 * more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with psirng.  If not, see <https://www.gnu.org/licenses/>.
 */

package grpc

import (
	"context"
	"github.com/hyperledger/fabric-lib-go/healthz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type HealthServerImpl struct {
	rngHealthChecker healthz.HealthChecker
	grpc_health_v1.HealthServer
}

func NewHealthServerImpl(rngHealthChecker healthz.HealthChecker) *HealthServerImpl {
	return &HealthServerImpl{
		rngHealthChecker: rngHealthChecker,
	}
}

func (s *HealthServerImpl) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	if in.Service == "rng" {
		if err := s.rngHealthChecker.HealthCheck(ctx); err != nil {
			return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING}, status.Errorf(codes.Unavailable, err.Error())
		}
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
	} else {
		return nil, status.Error(codes.NotFound, "unknown service")
	}
}

func (s *HealthServerImpl) Watch(_ *grpc_health_v1.HealthCheckRequest, _ grpc_health_v1.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
