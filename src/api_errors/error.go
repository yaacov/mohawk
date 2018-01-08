// Copyright 2016,2017,2018 Yaacov Zamir <kobi.zamir@gmail.com>
// and other contributors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package apiErrors API errors
package apiErrors

import (
	"fmt"
	"net/http"
)

// Error represent an error that occurred and its status code.
type Error struct {
	code    int
	message string
}

func (e Error) Error() string {
	return e.message
}

// JSON writes Error as JSON to http.ResponseWriter
func (e Error) JSON(w http.ResponseWriter) {
	w.WriteHeader(e.code)
	w.Write([]byte(fmt.Sprintf(`{"code":"%d","message":"%s"}`, e.code, e.message)))
}

// BadRequest creates new Error with status code 400 from error.
func BadRequest(err error) Error {
	return Error{code: http.StatusBadRequest, message: err.Error()}
}

// InternalError creates new Error with status code 500 from error.
func InternalError(err error) Error {
	return Error{code: http.StatusInternalServerError, message: err.Error()}
}