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

package httpapi

import (
	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-lib-go/healthz"
	"net/http"
)

type Router struct {
	MuxRouter        *mux.Router
	rngController    *RngController
	rngHealthChecker healthz.HealthChecker
}

func NewRouter(rngController *RngController, rngHealthChecker healthz.HealthChecker) *Router {
	router := Router{
		MuxRouter:        mux.NewRouter(),
		rngController:    rngController,
		rngHealthChecker: rngHealthChecker,
	}

	healthHandler := healthz.NewHealthHandler()
	_ = healthHandler.RegisterChecker("rng", rngHealthChecker)

	router.MuxRouter.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { healthHandler.ServeHTTP(w, r) }).Methods("GET")

	router.MuxRouter.HandleFunc("/api/randbooleans", rngController.RandBooleans).Methods("GET")
	router.MuxRouter.HandleFunc("/api/randbytes", rngController.RandBytes).Methods("GET")
	router.MuxRouter.HandleFunc("/api/randintegers", rngController.RandIntegers).Methods("GET")
	router.MuxRouter.HandleFunc("/api/randuniform", rngController.RandUniform).Methods("GET")
	router.MuxRouter.HandleFunc("/api/randnormal", rngController.RandNormal).Methods("GET")

	return &router
}
