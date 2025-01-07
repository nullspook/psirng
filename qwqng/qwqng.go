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

package qwqng

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

type Qwqng struct {
	qwqngx *C.qwqngx
}

func NewQwqng() *Qwqng {
	q := &Qwqng{}
	if err := q.initQwqngx(); err != nil {
		log.Println("failed to initialize qwqngx")
	}
	return q
}

func (q *Qwqng) initQwqngx() error {
	qwqngx, err := initQwqngx()
	if err != nil {
		return err
	}
	q.qwqngx = qwqngx
	return nil
}

func (q *Qwqng) Close() {
	C.qwqngx_free(q.qwqngx)
}

func (q *Qwqng) Clear() error {
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

func (q *Qwqng) RandBytes(dest []byte, length int32) error {
	if q.qwqngx == nil {
		if err := q.initQwqngx(); err != nil {
			log.Println(err)
		}
		return errors.New("qwqngx is initializing")
	}

	if C.qwqngx_randbytes(q.qwqngx, (*C.char)(unsafe.Pointer(&dest[0])), C.int(length)) != 0 {
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

func (q *Qwqng) RandIntegers(dest []int32, length, min, max int32) error {
	if q.qwqngx == nil {
		if err := q.initQwqngx(); err != nil {
			log.Println(err)
		}
		return errors.New("qwqngx is initializing")
	}

	if C.qwqngx_randintegers(q.qwqngx, (*C.int)(unsafe.Pointer(&dest[0])), C.int(length), C.int(min), C.int(max)) != 0 {
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

func (q *Qwqng) RandUniform(dest []float64, length int32, min, max float64) error {
	if q.qwqngx == nil {
		if err := q.initQwqngx(); err != nil {
			log.Println(err)
		}
		return errors.New("qwqngx is initializing")
	}

	if C.qwqngx_randuniform(q.qwqngx, (*C.double)(unsafe.Pointer(&dest[0])), C.int(length), C.double(min), C.double(max)) != 0 {
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

func (q *Qwqng) RandNormal(dest []float64, length int32, mean, stddev float64) error {
	if q.qwqngx == nil {
		if err := q.initQwqngx(); err != nil {
			log.Println(err)
		}
		return errors.New("qwqngx is initializing")
	}

	if C.qwqngx_randnormal(q.qwqngx, (*C.double)(unsafe.Pointer(&dest[0])), C.int(length), C.double(mean), C.double(stddev)) != 0 {
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
