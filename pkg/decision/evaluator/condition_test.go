/****************************************************************************
 * Copyright 2019-2021, Optimizely, Inc. and contributors                   *
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

package evaluator

import (
	"fmt"
	"testing"

	"github.com/WolffunService/experiment/pkg/decide"
	"github.com/WolffunService/experiment/pkg/entities"
	"github.com/WolffunService/experiment/pkg/logging"
	"github.com/stretchr/testify/suite"
)

type ConditionTestSuite struct {
	suite.Suite
	mockLogger         *MockLogger
	conditionEvaluator *CustomAttributeConditionEvaluator
	options            decide.Options
	reasons            decide.DecisionReasons
}

func (s *ConditionTestSuite) SetupTest() {
	s.mockLogger = new(MockLogger)
	s.conditionEvaluator = NewCustomAttributeConditionEvaluator(s.mockLogger)
	s.options = decide.Options{}
	s.reasons = decide.NewDecisionReasons(&s.options)
}

func (s *ConditionTestSuite) TestCustomAttributeConditionEvaluator() {
	condition := entities.Condition{
		Match: "exact",
		Value: "foo",
		Name:  "string_foo",
		Type:  "custom_attribute",
	}

	// Test condition passes
	user := entities.UserContext{
		Attributes: map[string]interface{}{
			"string_foo": "foo",
		},
	}

	condTreeParams := entities.NewTreeParameters(&user, map[string]entities.Audience{})
	result, _, _ := s.conditionEvaluator.Evaluate(condition, condTreeParams, &s.options)
	s.Equal(result, true)

	// Test condition fails
	user = entities.UserContext{
		Attributes: map[string]interface{}{
			"string_foo": "not_foo",
		},
	}
	result, _, _ = s.conditionEvaluator.Evaluate(condition, condTreeParams, &s.options)
	s.Equal(result, false)
}

func (s *ConditionTestSuite) TestCustomAttributeConditionEvaluatorWithoutMatchType() {
	condition := entities.Condition{
		Value: "foo",
		Name:  "string_foo",
		Type:  "custom_attribute",
	}

	// Test condition passes
	user := entities.UserContext{
		Attributes: map[string]interface{}{
			"string_foo": "foo",
		},
	}

	condTreeParams := entities.NewTreeParameters(&user, map[string]entities.Audience{})
	result, _, _ := s.conditionEvaluator.Evaluate(condition, condTreeParams, &s.options)
	s.Equal(result, true)

	// Test condition fails
	user = entities.UserContext{
		Attributes: map[string]interface{}{
			"string_foo": "not_foo",
		},
	}
	result, _, _ = s.conditionEvaluator.Evaluate(condition, condTreeParams, &s.options)
	s.Equal(result, false)
}

func (s *ConditionTestSuite) TestCustomAttributeConditionEvaluatorWithInvalidMatchType() {
	condition := entities.Condition{
		Value: "foo",
		Name:  "string_foo",
		Type:  "custom_attribute",
		Match: "invalid",
	}

	// Test condition fails
	user := entities.UserContext{
		Attributes: map[string]interface{}{
			"string_foo": "foo",
		},
	}

	condTreeParams := entities.NewTreeParameters(&user, map[string]entities.Audience{})
	s.mockLogger.On("Warning", fmt.Sprintf(logging.UnknownMatchType.String(), ""))
	s.options.IncludeReasons = true
	result, reasons, _ := s.conditionEvaluator.Evaluate(condition, condTreeParams, &s.options)
	s.Equal(result, false)
	messages := reasons.ToReport()
	s.Len(messages, 1)
	s.Equal(`invalid Condition matcher "invalid"`, messages[0])
	s.mockLogger.AssertExpectations(s.T())
}

func (s *ConditionTestSuite) TestCustomAttributeConditionEvaluatorWithUnknownType() {
	condition := entities.Condition{
		Value: "foo",
		Name:  "string_foo",
		Type:  "",
	}

	// Test condition fails
	user := entities.UserContext{
		Attributes: map[string]interface{}{
			"string_foo": "foo",
		},
	}

	condTreeParams := entities.NewTreeParameters(&user, map[string]entities.Audience{})
	s.mockLogger.On("Warning", fmt.Sprintf(logging.UnknownConditionType.String(), ""))
	s.options.IncludeReasons = true
	result, reasons, _ := s.conditionEvaluator.Evaluate(condition, condTreeParams, &s.options)
	s.Equal(result, false)
	messages := reasons.ToReport()
	s.Len(messages, 1)
	s.Equal(`unable to evaluate condition of type ""`, messages[0])
	s.mockLogger.AssertExpectations(s.T())
}

func (s *ConditionTestSuite) TestCustomAttributeConditionEvaluatorForGeSemver() {
	conditionEvaluator := CustomAttributeConditionEvaluator{}
	condition := entities.Condition{
		Match: "semver_ge",
		Value: "2.9",
		Name:  "string_foo",
		Type:  "custom_attribute",
	}

	// Test condition passes
	user := entities.UserContext{
		Attributes: map[string]interface{}{
			"string_foo": "2.9.1",
		},
	}

	condTreeParams := entities.NewTreeParameters(&user, map[string]entities.Audience{})
	result, _, _ := conditionEvaluator.Evaluate(condition, condTreeParams, &s.options)
	s.Equal(result, true)
}

func (s *ConditionTestSuite) TestCustomAttributeConditionEvaluatorForGeSemverBeta() {
	conditionEvaluator := CustomAttributeConditionEvaluator{}
	condition := entities.Condition{
		Match: "semver_ge",
		Value: "3.7.0",
		Name:  "string_foo",
		Type:  "custom_attribute",
	}

	// Test condition passes
	user := entities.UserContext{
		Attributes: map[string]interface{}{
			"string_foo": "3.7.1-beta",
		},
	}

	condTreeParams := entities.NewTreeParameters(&user, map[string]entities.Audience{})
	result, _, _ := conditionEvaluator.Evaluate(condition, condTreeParams, &s.options)
	s.Equal(true, result)
}

func (s *ConditionTestSuite) TestCustomAttributeConditionEvaluatorForGeSemverInvalid() {
	conditionEvaluator := CustomAttributeConditionEvaluator{}
	condition := entities.Condition{
		Match: "semver_ge",
		Value: "3.7.0",
		Name:  "string_foo",
		Type:  "custom_attribute",
	}

	// Test condition passes
	user := entities.UserContext{
		Attributes: map[string]interface{}{
			"string_foo": "3.7.2.2",
		},
	}

	condTreeParams := entities.NewTreeParameters(&user, map[string]entities.Audience{})
	_, _, err := conditionEvaluator.Evaluate(condition, condTreeParams, &s.options)
	s.NotNil(err)
}

func TestConditionTestSuite(t *testing.T) {
	suite.Run(t, new(ConditionTestSuite))
}
