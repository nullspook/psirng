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

/*
#cgo LDFLAGS: -lqwqngx
#include "qwqngx.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"log"
	"unsafe"
)

func initQwqngx() (*C.qwqngx, error) {
	var qwqngx *C.qwqngx
	if C.qwqngx_init(&qwqngx) != 0 {
		return nil, errors.New("failed to initialize qwqngx")
	}
	return qwqngx, nil
}

type QwqngxRngProvider struct {
	qwqngx *C.qwqngx
	baseRngProvider
}

func NewQwqngxRngProvider() (*QwqngxRngProvider, error) {
	q := &QwqngxRngProvider{}

	base, err := newBaseRngProvider(q)
	if err != nil {
		return nil, err
	}

	q.baseRngProvider = *base

	if err := q.initQwqngx(); err != nil {
		return nil, err
	}

	return q, nil
}

func (q *QwqngxRngProvider) initQwqngx() error {
	qwqngx, err := initQwqngx()
	if err != nil {
		return err
	}
	q.qwqngx = qwqngx
	return nil
}

func (q *QwqngxRngProvider) Close() {
	C.qwqngx_free(q.qwqngx)
}

func (q *QwqngxRngProvider) ClearBuffer() error {
	if q.qwqngx == nil {
		if err := q.initQwqngx(); err != nil {
			log.Println(err)
		}
		return errors.New("qwqngx is initializing")
	}

	if C.qwqngx_clear(q.qwqngx) != 0 {
		statusString := C.GoString(C.qwqngx_status_string(q.qwqngx))
		C.qwqngx_free(q.qwqngx)
		q.qwqngx = nil
		if err := q.initQwqngx(); err != nil {
			log.Println(err)
		}
		return errors.New(statusString)
	}
	return nil
}

func (q *QwqngxRngProvider) RandBytes(dest []byte) error {
	if q.qwqngx == nil {
		if err := q.initQwqngx(); err != nil {
			log.Println(err)
		}
		return errors.New("qwqngx is initializing")
	}

	if C.qwqngx_randbytes(q.qwqngx, (*C.char)(unsafe.Pointer(&dest[0])), C.int(len(dest))) != 0 {
		statusString := C.GoString(C.qwqngx_status_string(q.qwqngx))
		C.qwqngx_free(q.qwqngx)
		q.qwqngx = nil
		if err := q.initQwqngx(); err != nil {
			log.Println(err)
		}
		return errors.New(statusString)
	}
	return nil
}
