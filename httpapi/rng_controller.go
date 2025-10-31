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
	"github.com/go-playground/validator/v10"
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

	resp, err := c.rngService.RandBooleans(r.Context(), randBooleansRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, resp)
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

	resp, err := c.rngService.RandBytes(r.Context(), randBytesRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, resp)
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

	resp, err := c.rngService.RandIntegers(r.Context(), randIntegersRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, resp)
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

	resp, err := c.rngService.RandUniform(r.Context(), randUniformRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, resp)
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

	resp, err := c.rngService.RandNormal(r.Context(), randNormalRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, resp)
}

func (c *RngController) RandBooleansBiasAmplified(w http.ResponseWriter, r *http.Request) {
	var randBooleansBiasAmplifiedRequest models.RandBooleansBiasAmplifiedRequest
	if err := decoder.Decode(&randBooleansBiasAmplifiedRequest, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := validate.Struct(randBooleansBiasAmplifiedRequest); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	resp, err := c.rngService.RandBooleansBiasAmplified(r.Context(), randBooleansBiasAmplifiedRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, resp)
}

func (c *RngController) StreamBooleans(w http.ResponseWriter, r *http.Request) {
	var req models.StreamBooleansRequest
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	dataChan, errChan := c.rngService.StreamBooleans(r.Context(), req)
	writeNDJsonStreamResponse(w, dataChan, errChan)
}

func (c *RngController) StreamBooleansBiasAmplified(w http.ResponseWriter, r *http.Request) {
	var req models.StreamBooleansBiasAmplifiedRequest
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	dataChan, errChan := c.rngService.StreamBooleansBiasAmplified(r.Context(), req)
	writeNDJsonStreamResponse(w, dataChan, errChan)
}

func (c *RngController) StreamBytes(w http.ResponseWriter, r *http.Request) {
	var req models.StreamBytesRequest
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	dataChan, errChan := c.rngService.StreamBytes(r.Context(), req)
	writeNDJsonStreamResponse(w, dataChan, errChan)
}

func (c *RngController) StreamIntegers(w http.ResponseWriter, r *http.Request) {
	var req models.StreamIntegersRequest
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	dataChan, errChan := c.rngService.StreamIntegers(r.Context(), req)
	writeNDJsonStreamResponse(w, dataChan, errChan)
}

func (c *RngController) StreamNormal(w http.ResponseWriter, r *http.Request) {
	var req models.StreamNormalRequest
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	dataChan, errChan := c.rngService.StreamNormal(r.Context(), req)
	writeNDJsonStreamResponse(w, dataChan, errChan)
}

func (c *RngController) StreamUniform(w http.ResponseWriter, r *http.Request) {
	var req models.StreamUniformRequest
	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	dataChan, errChan := c.rngService.StreamUniform(r.Context(), req)
	writeNDJsonStreamResponse(w, dataChan, errChan)
}

func writeJsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeNDJsonStreamResponse[T any](w http.ResponseWriter, dataChan <-chan T, errChan <-chan error) {
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)

	for {
		select {
		case err, ok := <-errChan:
			if ok && err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		case chunk, ok := <-dataChan:
			if !ok {
				return
			}
			if err := enc.Encode(chunk); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
