// Copyright 2022 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/intstr"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func TestForEachStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect *ForEachState
		err    string
	}
	testCases := []testCase{
		{
			desp: "all field",
			data: `{"mode": "sequential"}`,
			expect: &ForEachState{
				Mode: ForEachModeTypeSequential,
			},
			err: ``,
		},
		{
			desp: "mode unset",
			data: `{}`,
			expect: &ForEachState{
				Mode: ForEachModeTypeParallel,
			},
			err: ``,
		},
		{
			desp:   "invalid json format",
			data:   `{"mode": 1}`,
			expect: nil,
			err:    `forEachState value '{"mode": 1}' is not supported, it must be an object or string`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v ForEachState
			err := json.Unmarshal([]byte(tc.data), &v)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, &v)
		})
	}
}

func TestForEachStateStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp  string
		state ForEachState
		err   string
	}
	testCases := []testCase{
		{
			desp: "normal test & sequential",
			state: ForEachState{
				BaseState: BaseState{
					Name: "1",
					Type: "2",
				},
				InputCollection: "3",
				Actions: []Action{
					{},
				},
				Mode: ForEachModeTypeSequential,
			},
			err: ``,
		},
		{
			desp: "normal test & parallel int",
			state: ForEachState{
				BaseState: BaseState{
					Name: "1",
					Type: "2",
				},
				InputCollection: "3",
				Actions: []Action{
					{},
				},
				Mode: ForEachModeTypeParallel,
				BatchSize: &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 1,
				},
			},
			err: ``,
		},
		{
			desp: "normal test & parallel string",
			state: ForEachState{
				BaseState: BaseState{
					Name: "1",
					Type: "2",
				},
				InputCollection: "3",
				Actions: []Action{
					{},
				},
				Mode: ForEachModeTypeParallel,
				BatchSize: &intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "1",
				},
			},
			err: ``,
		},
		{
			desp: "invalid parallel int",
			state: ForEachState{
				BaseState: BaseState{
					Name: "1",
					Type: "2",
				},
				InputCollection: "3",
				Actions: []Action{
					{},
				},
				Mode: ForEachModeTypeParallel,
				BatchSize: &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 0,
				},
			},
			err: `Key: 'ForEachState.BatchSize' Error:Field validation for 'BatchSize' failed on the 'gt0' tag`,
		},
		{
			desp: "invalid parallel string",
			state: ForEachState{
				BaseState: BaseState{
					Name: "1",
					Type: "2",
				},
				InputCollection: "3",
				Actions: []Action{
					{},
				},
				Mode: ForEachModeTypeParallel,
				BatchSize: &intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "0",
				},
			},
			err: `Key: 'ForEachState.BatchSize' Error:Field validation for 'BatchSize' failed on the 'gt0' tag`,
		},
		{
			desp: "invalid parallel string format",
			state: ForEachState{
				BaseState: BaseState{
					Name: "1",
					Type: "2",
				},
				InputCollection: "3",
				Actions: []Action{
					{},
				},
				Mode: ForEachModeTypeParallel,
				BatchSize: &intstr.IntOrString{
					Type:   intstr.String,
					StrVal: "a",
				},
			},
			err: `Key: 'ForEachState.BatchSize' Error:Field validation for 'BatchSize' failed on the 'gt0' tag`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.state)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}