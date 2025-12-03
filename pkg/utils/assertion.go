package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gotomicro/cetus/l"
	"github.com/gotomicro/ego/core/elog"
	"github.com/jmespath/go-jmespath"
)

// Validator defines a validation rule
type Validator struct {
	Check   string
	Expect  interface{}
	Assert  string
	Message string
}

// ValidationResult stores the validation outcome
type ValidationResult struct {
	Validator   Validator
	CheckValue  interface{}
	CheckResult string
}

// GetCheckValue extracts a field using jmespath
func GetCheckValue(responseBody []byte, check string) (interface{}, error) {
	// If check is provided but the response is empty, return an error
	if check != "" && len(responseBody) == 0 {
		errMsg := fmt.Sprintf("response body is empty, but check field is: %s", check)
		elog.Error("Unable to get field value", l.S("check", check), l.S("error", errMsg))
		return nil, fmt.Errorf(errMsg)
	}

	// If both the response and check are empty, return nil
	if len(responseBody) == 0 {
		return nil, nil
	}

	var data interface{}
	err := json.Unmarshal(responseBody, &data)
	if err != nil {
		elog.Error("failed to parse response JSON", l.S("error", err.Error()))
		return nil, fmt.Errorf("failed to parse response JSON: %v", err)
	}

	// Evaluate the jmespath expression
	result, err := jmespath.Search(check, data)
	if err != nil {
		elog.Error("jmespath execution failed", l.S("check", check), l.S("error", err.Error()))
		return nil, fmt.Errorf("jmespath execution failed: %v", err)
	}

	// Return an error if the field is missing
	if result == nil {
		errMsg := fmt.Sprintf("field not found: %s", check)
		elog.Error("failed to fetch check value", l.S("check", check), l.S("error", errMsg))
		return nil, fmt.Errorf(errMsg)
	}

	return result, nil
}

// Assertions Supported assertion types
var Assertions = map[string]func(actual interface{}, expected interface{}) (bool, string){
	"eq":                equal,
	"equals":            equal,
	"equal":             equal,
	"lt":                lessThan,
	"less_than":         lessThan,
	"le":                lessOrEqual,
	"less_or_equals":    lessOrEqual,
	"gt":                greaterThan,
	"greater_than":      greaterThan,
	"ge":                greaterOrEqual,
	"greater_or_equals": greaterOrEqual,
	"ne":                notEqual,
	"not_equal":         notEqual,
	"contains":          contains,
	"not_contains":      notContains,
}

// Assertion helpers
func equal(actual, expected interface{}) (bool, string) {
	actualStr, expectedStr := toString(actual), toString(expected)
	if actualStr == expectedStr {
		return true, ""
	}
	return false, fmt.Sprintf("expected %v, got %v", expected, actual)
}

func lessThan(actual, expected interface{}) (bool, string) {
	actualFloat, expectedFloat := toComparableFloat(actual), toComparableFloat(expected)
	if actualFloat < expectedFloat {
		return true, ""
	}
	return false, fmt.Sprintf("expected less than %v, got %v", expected, actual)
}

func lessOrEqual(actual, expected interface{}) (bool, string) {
	actualFloat, expectedFloat := toComparableFloat(actual), toComparableFloat(expected)
	if actualFloat <= expectedFloat {
		return true, ""
	}
	return false, fmt.Sprintf("expected <= %v, got %v", expected, actual)
}

func greaterThan(actual, expected interface{}) (bool, string) {
	actualFloat, expectedFloat := toComparableFloat(actual), toComparableFloat(expected)
	if actualFloat > expectedFloat {
		return true, ""
	}
	return false, fmt.Sprintf("expected > %v, got %v", expected, actual)
}

func greaterOrEqual(actual, expected interface{}) (bool, string) {
	actualFloat, expectedFloat := toComparableFloat(actual), toComparableFloat(expected)
	if actualFloat >= expectedFloat {
		return true, ""
	}
	return false, fmt.Sprintf("expected >= %v, got %v", expected, actual)
}

func notEqual(actual, expected interface{}) (bool, string) {
	actualStr, expectedStr := toString(actual), toString(expected)
	if actualStr != expectedStr {
		return true, ""
	}
	return false, fmt.Sprintf("expected != %v, got %v", expected, actual)
}

func contains(actual, expected interface{}) (bool, string) {
	actualStr, expectedStr := toString(actual), toString(expected)
	if containsStr(actualStr, expectedStr) {
		return true, ""
	}
	return false, fmt.Sprintf("expected to contain %v, got %v", expectedStr, actualStr)
}

func notContains(actual, expected interface{}) (bool, string) {
	actualStr, expectedStr := toString(actual), toString(expected)
	if !containsStr(actualStr, expectedStr) {
		return true, ""
	}
	return false, fmt.Sprintf("actual value contains %v; actual: %v", expectedStr, actualStr)
}

// Helper: does str contain substr?
func containsStr(str, substr string) bool {
	return strings.Contains(str, substr)
}

// Validate checks whether the assertion holds
func Validate(checkValue interface{}, assertion string, expect interface{}) (string, error) {
	// Select the assertion function
	assertFunc, ok := Assertions[assertion]
	if !ok {
		errMsg := fmt.Sprintf("unsupported assertion: %s", assertion)
		elog.Error("invalid assertion", l.S("assertion", assertion), l.S("error", errMsg))
		return errMsg, fmt.Errorf("unsupported assertion: %s", assertion)
	}

	// Execute the assertion
	success, errMsg := assertFunc(checkValue, expect)
	if !success {
		elog.Error("assertion failed",
			l.S("assertion", assertion),
			l.S("error", errMsg),
			l.A("checkValue", checkValue),
			l.A("expect", expect))
		return fmt.Sprintf("assertion failed: %s", errMsg), nil
	}

	elog.Info("assertion succeeded",
		l.S("assertion", assertion),
		l.A("checkValue", checkValue),
		l.A("expect", expect))

	return "assertion succeeded", nil
}
