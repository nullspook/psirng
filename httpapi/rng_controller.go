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
	"encoding/json"
	"github.com/gorilla/schema"
	"net/http"
	"psirng/models"
	"psirng/services"
)

type RngController struct {
	rngService *services.RngService
}

var (
	decoder  = schema.NewDecoder()
	validate = validator.New()
)

func init() {
	decoder.IgnoreUnknownKeys(true)
}

func NewRngController(rngService *services.RngService) *RngController {
	return &RngController{rngService: rngService}
}

func (c *RngController) RandBooleans(w http.ResponseWriter, r *http.Request) {
	var randBooleansRequest models.RandBooleansRequest
	if err := decoder.Decode(&randBooleansRequest, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Struct(randBooleansRequest); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	data, err := c.rngService.RandBooleans(randBooleansRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, data)
}

func (c *RngController) RandBytes(w http.ResponseWriter, r *http.Request) {
	var randBytesRequest models.RandBytesRequest
	if err := decoder.Decode(&randBytesRequest, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Struct(randBytesRequest); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	data, err := c.rngService.RandBytes(randBytesRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, data)
}

func (c *RngController) RandIntegers(w http.ResponseWriter, r *http.Request) {
	var randIntegersRequest models.RandIntegersRequest
	if err := decoder.Decode(&randIntegersRequest, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Struct(randIntegersRequest); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	data, err := c.rngService.RandIntegers(randIntegersRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, data)
}

func (c *RngController) RandUniform(w http.ResponseWriter, r *http.Request) {
	var randUniformRequest models.RandUniformRequest
	if err := decoder.Decode(&randUniformRequest, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Struct(randUniformRequest); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	data, err := c.rngService.RandUniform(randUniformRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, data)
}

func (c *RngController) RandNormal(w http.ResponseWriter, r *http.Request) {
	var randNormalRequest models.RandNormalRequest
	if err := decoder.Decode(&randNormalRequest, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Struct(randNormalRequest); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	data, err := c.rngService.RandNormal(randNormalRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, data)
}

func writeJsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
