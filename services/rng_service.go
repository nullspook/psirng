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
	"context"
	"psirng/models"
	"psirng/providers/rng"
	"runtime"
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

func (s *RngService) RandBooleans(ctx context.Context, request models.RandBooleansRequest) (*models.RandBooleansResponse, error) {
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

func (s *RngService) RandBytes(ctx context.Context, request models.RandBytesRequest) (*models.RandBytesResponse, error) {
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

func (s *RngService) RandIntegers(ctx context.Context, request models.RandIntegersRequest) (*models.RandIntegersResponse, error) {
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

func (s *RngService) RandUniform(ctx context.Context, request models.RandUniformRequest) (*models.RandUniformResponse, error) {
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

func (s *RngService) RandNormal(ctx context.Context, request models.RandNormalRequest) (*models.RandNormalResponse, error) {
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

func (s *RngService) RandBooleansBiasAmplified(ctx context.Context, request models.RandBooleansBiasAmplifiedRequest) (*models.RandBooleansBiasAmplifiedResponse, error) {
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

func (s *RngService) StreamBooleans(ctx context.Context, request models.StreamBooleansRequest) (<-chan models.StreamBooleansResponseChunk, <-chan error) {
	dataChannel := make(chan models.StreamBooleansResponseChunk)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(dataChannel)
		defer close(errorChannel)

		chunkSize := *request.ChunkSize

		byteBufferLength := (chunkSize + 7) / 8
		byteBuffer := make([]byte, byteBufferLength)

		s.mutex.Lock()

		if err := s.rng.ClearBuffer(); err != nil {
			errorChannel <- err
			s.mutex.Unlock()
			return
		}

		s.mutex.Unlock()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.rng.RandBytes(byteBuffer); err != nil {
					errorChannel <- err
					return
				}

				data := make([]bool, chunkSize)
				for i := 0; i < int(chunkSize); i++ {
					data[i] = (byteBuffer[i>>3] >> (i & 7) & 1) != 0
				}

				dataChannel <- models.StreamBooleansResponseChunk{Data: data}
			}
		}
	}()

	return dataChannel, errorChannel
}

func (s *RngService) StreamBytes(ctx context.Context, request models.StreamBytesRequest) (<-chan models.StreamBytesResponseChunk, <-chan error) {
	dataChannel := make(chan models.StreamBytesResponseChunk)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(dataChannel)
		defer close(errorChannel)

		chunkSize := *request.ChunkSize
		byteBuffer := make([]byte, chunkSize)

		s.mutex.Lock()

		if err := s.rng.ClearBuffer(); err != nil {
			errorChannel <- err
			s.mutex.Unlock()
			return
		}

		s.mutex.Unlock()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.rng.RandBytes(byteBuffer); err != nil {
					errorChannel <- err
					return
				}

				data := make([]byte, chunkSize)
				copy(data, byteBuffer)

				dataChannel <- models.StreamBytesResponseChunk{Data: data}
			}

			runtime.Gosched()
		}
	}()

	return dataChannel, errorChannel
}

func (s *RngService) StreamIntegers(ctx context.Context, request models.StreamIntegersRequest) (<-chan models.StreamIntegersResponseChunk, <-chan error) {
	dataChannel := make(chan models.StreamIntegersResponseChunk)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(dataChannel)
		defer close(errorChannel)

		chunkSize := *request.ChunkSize
		min := *request.Min
		max := *request.Max

		for {
			select {
			case <-ctx.Done():
				return
			default:
				intBuffer := make([]int32, chunkSize)

				s.mutex.Lock()

				if err := s.rng.ClearBuffer(); err != nil {
					errorChannel <- err
					s.mutex.Unlock()
					return
				}

				if err := s.rng.RandIntegers(intBuffer, min, max); err != nil {
					errorChannel <- err
					s.mutex.Unlock()
					return
				}

				s.mutex.Unlock()

				data := make([]int32, chunkSize)
				copy(data, intBuffer)

				dataChannel <- models.StreamIntegersResponseChunk{Data: data}
			}

			runtime.Gosched()
		}
	}()

	return dataChannel, errorChannel
}

func (s *RngService) StreamNormal(ctx context.Context, request models.StreamNormalRequest) (<-chan models.StreamNormalResponseChunk, <-chan error) {
	dataChannel := make(chan models.StreamNormalResponseChunk)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(dataChannel)
		defer close(errorChannel)

		chunkSize := *request.ChunkSize
		mean := *request.Mean
		stdDev := *request.StdDev

		for {
			select {
			case <-ctx.Done():
				return
			default:
				floatBuffer := make([]float64, chunkSize)

				s.mutex.Lock()

				if err := s.rng.ClearBuffer(); err != nil {
					errorChannel <- err
					s.mutex.Unlock()
					return
				}

				if err := s.rng.RandNormal(floatBuffer, mean, stdDev); err != nil {
					errorChannel <- err
					s.mutex.Unlock()
					return
				}

				s.mutex.Unlock()

				data := make([]float64, chunkSize)
				copy(data, floatBuffer)

				dataChannel <- models.StreamNormalResponseChunk{Data: data}
			}

			runtime.Gosched()
		}
	}()

	return dataChannel, errorChannel
}

func (s *RngService) StreamUniform(ctx context.Context, request models.StreamUniformRequest) (<-chan models.StreamUniformResponseChunk, <-chan error) {
	dataChannel := make(chan models.StreamUniformResponseChunk)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(dataChannel)
		defer close(errorChannel)

		chunkSize := *request.ChunkSize
		min := *request.Min
		max := *request.Max

		for {
			select {
			case <-ctx.Done():
				return
			default:
				floatBuffer := make([]float64, chunkSize)

				s.mutex.Lock()

				if err := s.rng.ClearBuffer(); err != nil {
					errorChannel <- err
					s.mutex.Unlock()
					return
				}

				if err := s.rng.RandUniform(floatBuffer, min, max); err != nil {
					errorChannel <- err
					s.mutex.Unlock()
					return
				}

				s.mutex.Unlock()

				data := make([]float64, chunkSize)
				copy(data, floatBuffer)

				dataChannel <- models.StreamUniformResponseChunk{Data: data}
			}

			runtime.Gosched()
		}
	}()

	return dataChannel, errorChannel
}

func (s *RngService) StreamBooleansBiasAmplified(ctx context.Context, request models.StreamBooleansBiasAmplifiedRequest) (<-chan models.StreamBooleansResponseChunk, <-chan error) {
	dataChannel := make(chan models.StreamBooleansResponseChunk)
	errorChannel := make(chan error, 1)

	go func() {
		defer close(dataChannel)
		defer close(errorChannel)

		chunkSize := *request.ChunkSize

		// Level 0: bound = 1 (no amplification)
		// Level 1: bound = 2 (reverse von Neumann)
		// Level 2: bound = 3
		// Level 3: bound = 4
		// ...
		bound := int(*request.AmplificationLevel + 1)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				s.mutex.Lock()

				if err := s.rng.ClearBuffer(); err != nil {
					errorChannel <- err
					s.mutex.Unlock()
					return
				}

				var byteBuffer [1]byte
				bitIndex := 8
				randBoolean := func() (bool, error) {
					if bitIndex == 8 {
						if err := s.rng.RandBytes(byteBuffer[:]); err != nil {
							s.mutex.Unlock()
							return false, err
						}
						bitIndex = 0
					}
					bit := (byteBuffer[0] >> bitIndex) & 1
					bitIndex++
					return bit != 0, nil
				}

				data := make([]bool, chunkSize)
				for i := 0; i < int(chunkSize); i++ {
					position := 0

					for {
						b, err := randBoolean()
						if err != nil {
							s.mutex.Unlock()
							return
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

				s.mutex.Unlock()

				dataChannel <- models.StreamBooleansResponseChunk{Data: data}
			}

			runtime.Gosched()
		}
	}()

	return dataChannel, errorChannel
}
