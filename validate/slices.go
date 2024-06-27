package validate

import (
	"fmt"
	"reflect"

	p "github.com/anthdm/superkit/validate/primitives"
)

type sliceValidator struct {
	Rules      []p.Rule
	IsOptional bool
}

func Slice(schema fieldValidator) *sliceValidator {
	return &sliceValidator{
		Rules: []p.Rule{
			{
				Name:         "sliceItemsMatchSchema",
				RuleValue:    schema,
				ErrorMessage: "all items should match the schema",
				ValidateFunc: func(set p.Rule) bool {
					rv := reflect.ValueOf(set.FieldValue)
					if rv.Kind() != reflect.Slice {
						return false
					}
					s, ok := set.RuleValue.(fieldValidator)
					if !ok {
						return false
					}
					for idx := 0; idx < rv.Len(); idx++ {
						v := rv.Index(idx).Interface()
						_, ok := s.Validate(v)
						if !ok {
							return false
						}
					}
					return true
				},
			},
		},
	}
}

// GLOBAL METHODS

func (v *sliceValidator) Refine(ruleName string, errorMsg string, validateFunc p.RuleValidateFunc) *sliceValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         ruleName,
			ErrorMessage: errorMsg,
			ValidateFunc: validateFunc,
		},
	)

	return v
}

func (v *sliceValidator) Optional() *sliceValidator {
	v.IsOptional = true
	return v
}

func (v *sliceValidator) Validate(fieldValue any) ([]string, bool) {
	return p.GenericValidator(fieldValue, v.Rules, v.IsOptional)
}

// UNIQUE METHODS

// TODO
// some & every -> pass a validator

func (v *sliceValidator) NotEmpty() *sliceValidator {
	v.Rules = append(v.Rules,
		sliceMin(1, "should not be empty"),
	)
	return v
}

// Minimum number of items
func (v *sliceValidator) Min(n int) *sliceValidator {
	v.Rules = append(v.Rules,
		sliceMin(n, fmt.Sprintf("should be at least %d items long", n)),
	)
	return v
}

// Maximum number of items
func (v *sliceValidator) Max(n int) *sliceValidator {
	v.Rules = append(v.Rules,
		sliceMax(n, fmt.Sprintf("should be at maximum %d items long", n)),
	)
	return v
}

// Exact number of items
func (v *sliceValidator) Len(n int) *sliceValidator {
	v.Rules = append(v.Rules,
		sliceLength(n, fmt.Sprintf("should be exactly %d items long", n)),
	)
	return v
}

func (v *sliceValidator) Contains(val any) *sliceValidator {
	v.Rules = append(v.Rules,
		p.Rule{
			Name:         "contains",
			RuleValue:    val,
			ErrorMessage: fmt.Sprintf("should contain %v", val),
			ValidateFunc: func(set p.Rule) bool {
				rv := reflect.ValueOf(set.FieldValue)
				if rv.Kind() != reflect.Slice {
					return false
				}
				for idx := 0; idx < rv.Len(); idx++ {
					v := rv.Index(idx).Interface()

					if reflect.DeepEqual(v, val) {
						return true
					}
				}

				return false
			},
		},
	)
	return v
}

func sliceMin(n int, errMsg string) p.Rule {
	return p.Rule{
		Name:         "sliceMin",
		RuleValue:    n,
		ErrorMessage: errMsg,
		ValidateFunc: func(set p.Rule) bool {
			rv := reflect.ValueOf(set.FieldValue)
			if rv.Kind() != reflect.Slice {
				return false
			}
			return rv.Len() >= n
		},
	}
}
func sliceMax(n int, errMsg string) p.Rule {
	return p.Rule{
		Name:         "sliceMax",
		RuleValue:    n,
		ErrorMessage: errMsg,
		ValidateFunc: func(set p.Rule) bool {
			rv := reflect.ValueOf(set.FieldValue)
			if rv.Kind() != reflect.Slice {
				return false
			}
			return rv.Len() <= n
		},
	}
}
func sliceLength(n int, errMsg string) p.Rule {
	return p.Rule{
		Name:         "sliceLength",
		RuleValue:    n,
		ErrorMessage: errMsg,
		ValidateFunc: func(set p.Rule) bool {
			rv := reflect.ValueOf(set.FieldValue)
			if rv.Kind() != reflect.Slice {
				return false
			}
			return rv.Len() == n
		},
	}
}
