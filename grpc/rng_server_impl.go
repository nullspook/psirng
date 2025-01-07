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
	"psirng/models"
	"psirng/services"
)

type RngServerImpl struct {
	rngService *services.RngService
	RngServer
}

func NewRngServerImpl(rngService *services.RngService) *RngServerImpl {
	return &RngServerImpl{rngService: rngService}
}

func (s *RngServerImpl) RandBooleans(_ context.Context, request *RandBooleansRequest) (*RandBooleansResponse, error) {
	data, err := s.rngService.RandBooleans(models.RandBooleansRequest{
		Length: &request.Length,
	})
	if err != nil {
		return nil, err
	}
	return &RandBooleansResponse{Data: data}, nil
}

func (s *RngServerImpl) RandBytes(_ context.Context, request *RandBytesRequest) (*RandBytesResponse, error) {
	data, err := s.rngService.RandBytes(models.RandBytesRequest{
		Length: &request.Length,
	})
	if err != nil {
		return nil, err
	}
	return &RandBytesResponse{Data: data}, nil
}

func (s *RngServerImpl) RandIntegers(_ context.Context, request *RandIntegersRequest) (*RandIntegersResponse, error) {
	data, err := s.rngService.RandIntegers(models.RandIntegersRequest{
		Length: &request.Length,
		Min:    &request.Min,
		Max:    &request.Max,
	})
	if err != nil {
		return nil, err
	}
	return &RandIntegersResponse{Data: data}, nil
}

func (s *RngServerImpl) RandUniform(_ context.Context, request *RandUniformRequest) (*RandUniformResponse, error) {
	data, err := s.rngService.RandUniform(models.RandUniformRequest{
		Length: &request.Length,
		Min:    &request.Min,
		Max:    &request.Max,
	})
	if err != nil {
		return nil, err
	}
	return &RandUniformResponse{Data: data}, nil
}

func (s *RngServerImpl) RandNormal(_ context.Context, request *RandNormalRequest) (*RandNormalResponse, error) {
	data, err := s.rngService.RandNormal(models.RandNormalRequest{
		Length: &request.Length,
		Mean:   &request.Mean,
		StdDev: &request.Stddev,
	})
	if err != nil {
		return nil, err
	}
	return &RandNormalResponse{Data: data}, nil
}
