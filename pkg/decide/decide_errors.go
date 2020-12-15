/****************************************************************************
 * Copyright 2020, Optimizely, Inc. and contributors                        *
 *                                                                          *
 * Licensed under the Apache License, Version 2.0 (the "License");          *
 * you may not use this file except in compliance with the License.         *
 * You may obtain a copy of the License at                                  *
 *                                                                          *
 *    http://www.apache.org/licenses/LICENSE-2.0                            *
 *                                                                          *
 * Unless required by applicable law or agreed to in writing, software      *
 * distributed under the License is distributed on an "AS IS" BASIS,        *
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 * See the License for the specific language governing permissions and      *
 * limitations under the License.                                           *
 ***************************************************************************/

// Package decide has error definitions for decide api
package decide

import "errors"

type decideError string

const (
	// SDKNotReady when sdk is not ready
	SDKNotReady decideError = "Optimizely SDK not configured properly yet"
	// FlagKeyInvalid when invalid flag key is provided
	FlagKeyInvalid decideError = `No flag was found for key "%s".`
	// VariableValueInvalid when invalid variable value is provided
	VariableValueInvalid decideError = `Variable value for key "%s" is invalid or wrong type.`
)

// GetError returns error for decide error type
func GetError(errorType decideError) error {
	return errors.New(string(errorType))
}