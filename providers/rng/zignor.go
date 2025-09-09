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

import "math"

// J. Doornik, "An Improved Ziggurat Method to Generate Normal Random Samples", University of Oxford, 2005.

const (
	zignorC = 128
	zignorR = 3.442619855899
	zignorV = 9.91256303526217e-3
)

type zignor struct {
	adZigX [zignorC + 1]float64
	adZigR [zignorC]float64
	dRanU  func() (float64, error)
	iRanU  func() (uint32, error)
}

func newZignor(dRanU func() (float64, error), iRanU func() (uint32, error)) *zignor {
	z := &zignor{
		dRanU: dRanU,
		iRanU: iRanU,
	}

	f := math.Exp(-0.5 * zignorR * zignorR)
	z.adZigX[0] = zignorV / f
	z.adZigX[1] = zignorR
	z.adZigX[zignorC] = 0
	for i := 2; i < zignorC; i++ {
		z.adZigX[i] = math.Sqrt(-2 * math.Log(zignorV/z.adZigX[i-1]+f))
		f = math.Exp(-0.5 * z.adZigX[i] * z.adZigX[i])
	}
	for i := 0; i < zignorC; i++ {
		z.adZigR[i] = z.adZigX[i+1] / z.adZigX[i]
	}

	return z
}

func (z *zignor) dRanNormalZig() (float64, error) {
	var (
		dRanUValue float64
		iRanUValue uint32
		err        error
	)

	for {
		if dRanUValue, err = z.dRanU(); err != nil {
			return 0, err
		}
		u := 2*dRanUValue - 1

		if iRanUValue, err = z.iRanU(); err != nil {
			return 0, err
		}
		i := iRanUValue & 0x7F

		if math.Abs(u) < z.adZigR[i] {
			return u * z.adZigX[i], nil
		}
		if i == 0 {
			return z.dRanNormalTail(zignorR, u < 0)
		}

		x := u * z.adZigX[i]
		f0 := math.Exp(-0.5 * (z.adZigX[i]*z.adZigX[i] - x*x))
		f1 := math.Exp(-0.5 * (z.adZigX[i+1]*z.adZigX[i+1] - x*x))

		if dRanUValue, err = z.dRanU(); err != nil {
			return 0, err
		}
		if f1+dRanUValue*(f0-f1) < 1.0 {
			return x, nil
		}
	}
}

func (z *zignor) dRanNormalTail(dMin float64, iNegative bool) (float64, error) {
	var (
		dRanUValue float64
		err        error
	)

	for {
		if dRanUValue, err = z.dRanU(); err != nil {
			return 0, err
		}
		x := -math.Log(dRanUValue) / dMin

		if dRanUValue, err = z.dRanU(); err != nil {
			return 0, err
		}
		y := -math.Log(dRanUValue)

		if -2*y > x*x {
			if iNegative {
				return x - dMin, nil
			}
			return dMin - x, nil
		}
	}
}
