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

package services

import (
	"psirng/models"
	"psirng/providers/rng"
	"sync"
)

type RngService struct {
	rng   rng.Provider
	mutex *sync.Mutex
}

func NewRngService(rng rng.Provider) *RngService {
	return &RngService{
		rng:   rng,
		mutex: &sync.Mutex{},
	}
}

func (s *RngService) Close() {
	s.rng.Close()
}

func (s *RngService) RandBooleans(request models.RandBooleansRequest) (*models.RandBooleansResponse, error) {
	data := make([]bool, *request.Length)

	byteBufferLength := (*request.Length + 7) / 8
	byteBuffer := make([]byte, byteBufferLength)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.rng.ClearBuffer(); err != nil {
		return nil, err
	}

	if err := s.rng.RandBytes(byteBuffer); err != nil {
		return nil, err
	}

	for i := 0; i < int(*request.Length); i++ {
		data[i] = (byteBuffer[i>>3] >> (i & 7) & 1) != 0
	}

	return &models.RandBooleansResponse{Data: data}, nil
}

func (s *RngService) RandBytes(request models.RandBytesRequest) (*models.RandBytesResponse, error) {
	data := make([]byte, *request.Length)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.rng.ClearBuffer(); err != nil {
		return nil, err
	}

	if err := s.rng.RandBytes(data); err != nil {
		return nil, err
	}

	return &models.RandBytesResponse{Data: data}, nil
}

func (s *RngService) RandIntegers(request models.RandIntegersRequest) (*models.RandIntegersResponse, error) {
	data := make([]int32, *request.Length)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.rng.ClearBuffer(); err != nil {
		return nil, err
	}

	if err := s.rng.RandIntegers(data, *request.Min, *request.Max); err != nil {
		return nil, err
	}

	return &models.RandIntegersResponse{Data: data}, nil
}

func (s *RngService) RandUniform(request models.RandUniformRequest) (*models.RandUniformResponse, error) {
	data := make([]float64, *request.Length)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.rng.ClearBuffer(); err != nil {
		return nil, err
	}

	if err := s.rng.RandUniform(data, *request.Min, *request.Max); err != nil {
		return nil, err
	}

	return &models.RandUniformResponse{Data: data}, nil
}

func (s *RngService) RandNormal(request models.RandNormalRequest) (*models.RandNormalResponse, error) {
	data := make([]float64, *request.Length)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.rng.ClearBuffer(); err != nil {
		return nil, err
	}

	if err := s.rng.RandNormal(data, *request.Mean, *request.StdDev); err != nil {
		return nil, err
	}

	return &models.RandNormalResponse{Data: data}, nil
}

func (s *RngService) RandBooleansBiasAmplified(request models.RandBooleansBiasAmplifiedRequest) (*models.RandBooleansBiasAmplifiedResponse, error) {
	data := make([]bool, *request.Length)

	// Level 0: bound = 1 (no amplification)
	// Level 1: bound = 2 (reverse von Neumann)
	// Level 2: bound = 3
	// Level 3: bound = 4
	// ...
	bound := int(*request.AmplificationLevel + 1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.rng.ClearBuffer(); err != nil {
		return nil, err
	}

	var byteBuffer [1]byte
	bitIndex := 8
	randBoolean := func() (bool, error) {
		if bitIndex == 8 {
			if err := s.rng.RandBytes(byteBuffer[:]); err != nil {
				return false, err
			}
			bitIndex = 0
		}
		bit := (byteBuffer[0] >> bitIndex) & 1
		bitIndex++
		return bit != 0, nil
	}

	for i := 0; i < int(*request.Length); i++ {
		position := 0

		for {
			b, err := randBoolean()
			if err != nil {
				return nil, err
			}

			if b {
				position++
			} else {
				position--
			}

			if position == bound || position == -bound {
				break
			}
		}

		data[i] = position > 0
	}

	return &models.RandBooleansBiasAmplifiedResponse{Data: data}, nil
}
