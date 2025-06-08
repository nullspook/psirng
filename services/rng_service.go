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
	"psirng/qwqng"
	"sync"
)

type RngService struct {
	qwqng *qwqng.Qwqng
	mutex *sync.Mutex
}

func NewRngService(qwqng *qwqng.Qwqng) *RngService {
	return &RngService{
		qwqng: qwqng,
		mutex: &sync.Mutex{},
	}
}

func (s *RngService) Close() {
	s.qwqng.Close()
}

func (s *RngService) RandBooleans(request models.RandBooleansRequest) ([]bool, error) {
	result := make([]bool, *request.Length)

	bufferLength := (*request.Length + 7) >> 3
	buffer := make([]byte, bufferLength)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.qwqng.Clear(); err != nil {
		return nil, err
	}

	if err := s.qwqng.RandBytes(buffer, int32(bufferLength)); err != nil {
		return nil, err
	}

	for i := 0; i < int(*request.Length); i++ {
		result[i] = (buffer[i>>3] >> (i & 7) & 1) != 0
	}

	return result, nil
}

func (s *RngService) RandBytes(request models.RandBytesRequest) ([]byte, error) {
	result := make([]byte, *request.Length)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.qwqng.Clear(); err != nil {
		return nil, err
	}

	if err := s.qwqng.RandBytes(result, int32(*request.Length)); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *RngService) RandIntegers(request models.RandIntegersRequest) ([]int32, error) {
	result := make([]int32, *request.Length)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.qwqng.Clear(); err != nil {
		return nil, err
	}

	if err := s.qwqng.RandIntegers(result, int32(*request.Length), *request.Min, *request.Max); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *RngService) RandUniform(request models.RandUniformRequest) ([]float64, error) {
	result := make([]float64, *request.Length)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.qwqng.Clear(); err != nil {
		return nil, err
	}

	if err := s.qwqng.RandUniform(result, int32(*request.Length), *request.Min, *request.Max); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *RngService) RandNormal(request models.RandNormalRequest) ([]float64, error) {
	result := make([]float64, *request.Length)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if err := s.qwqng.Clear(); err != nil {
		return nil, err
	}

	if err := s.qwqng.RandNormal(result, int32(*request.Length), *request.Mean, *request.StdDev); err != nil {
		return nil, err
	}

	return result, nil
}
