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

package rng

import (
	"encoding/binary"
	"math"
)

type base struct {
	zignor zignor
	Provider
}

func newBase(self Provider) (*base, error) {
	rp := &base{
		Provider: self,
	}

	rp.zignor = *newZignor(
		func() (float64, error) {
			var randomFloat64Buffer [1]float64
			if err := rp.RandUniform(randomFloat64Buffer[:], 0, 1); err != nil {
				return 0, err
			}
			return randomFloat64Buffer[0], nil
		},
		func() (uint32, error) {
			var randomInt32Buffer [1]int32
			if err := rp.RandIntegers(randomInt32Buffer[:], math.MinInt32, math.MaxInt32); err != nil {
				return 0, err
			}
			return uint32(randomInt32Buffer[0]), nil
		},
	)

	return rp, nil
}

func (b *base) RandIntegers(dest []int32, min, max int32) error {
	if min == max {
		// only one possible value
		for i := 0; i < len(dest); i++ {
			dest[i] = min
		}
		return nil
	}

	if min == math.MinInt32 && max == math.MaxInt32 {
		// full int32 range
		byteBuffer := make([]byte, len(dest)*4)
		if err := b.RandBytes(byteBuffer); err != nil {
			return err
		}
		for i := 0; i < len(dest); i++ {
			offset := i * 4
			dest[i] = int32(binary.LittleEndian.Uint32(byteBuffer[offset : offset+4]))
		}
		return nil
	}

	bound := uint32(uint64(int64(max)-int64(min)) + 1)
	mask := bound - 1

	if bound&mask == 0 {
		// bound is a power of two
		byteBuffer := make([]byte, len(dest)*4)
		if err := b.RandBytes(byteBuffer); err != nil {
			return err
		}
		for i := 0; i < len(dest); i++ {
			offset := i * 4
			random32 := binary.LittleEndian.Uint32(byteBuffer[offset : offset+4])
			dest[i] = min + int32(random32&mask)
		}
		return nil
	}

	var randomUInt32Buffer [4]byte
	for i := 0; i < len(dest); i++ {
		// M.E. O'Neill. "Efficiently Generating a Number in a Range". PCG, A Better Random Number Generator.
		// https://www.pcg-random.org/posts/bounded-rands.html
		if err := b.RandBytes(randomUInt32Buffer[:]); err != nil {
			return err
		}
		random32 := binary.LittleEndian.Uint32(randomUInt32Buffer[:])
		m := uint64(random32) * uint64(bound)
		l := uint32(m)
		if l < bound {
			t := -bound
			if t >= bound {
				t -= bound
				if t >= bound {
					t %= bound
				}
			}
			for l < t {
				if err := b.RandBytes(randomUInt32Buffer[:]); err != nil {
					return err
				}
				random32 = binary.LittleEndian.Uint32(randomUInt32Buffer[:])
				m = uint64(random32) * uint64(bound)
				l = uint32(m)
			}
		}
		dest[i] = min + int32(m>>32)
	}
	return nil
}

func (b *base) RandUniform(dest []float64, min, max float64) error {
	scale := max - min
	byteBuffer := make([]byte, len(dest)*8)

	if err := b.RandBytes(byteBuffer); err != nil {
		return err
	}

	for i := 0; i < len(dest); i++ {
		random64 := binary.LittleEndian.Uint64(byteBuffer[i*8 : i*8+8])
		dest[i] = min + float64(random64>>11)*0x1.0p-53*scale
	}

	return nil
}

func (b *base) RandNormal(dest []float64, mean, stddev float64) error {
	var (
		dRanNormalZigValue float64
		err                error
	)

	for i := 0; i < len(dest); i++ {
		if dRanNormalZigValue, err = b.zignor.dRanNormalZig(); err != nil {
			return err
		}
		dest[i] = mean + stddev*dRanNormalZigValue
	}

	return nil
}
