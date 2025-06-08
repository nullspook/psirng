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

package healthcheckers

import (
	"context"
	"psirng/models"
	"psirng/services"
)

type RngHealthChecker struct {
	rngService *services.RngService
}

var randBytesCheckerRequest = models.RandBytesRequest{Length: &[]uint32{uint32(1)}[0]}

func NewRngHealthChecker(rng *services.RngService) *RngHealthChecker {
	return &RngHealthChecker{rngService: rng}
}

func (hc *RngHealthChecker) HealthCheck(_ context.Context) error {
	if _, err := hc.rngService.RandBytes(randBytesCheckerRequest); err != nil {
		if _, err := hc.rngService.RandBytes(randBytesCheckerRequest); err != nil {
			return err
		}
	}
	return nil
}
