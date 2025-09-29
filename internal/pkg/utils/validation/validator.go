package validation

// Rule represents a validation rule function
type Rule func(value interface{}) string

// Validator holds validation rules
type Validator struct {
	rules []Rule
}

// New creates a new validator
func New() *Validator {
	return &Validator{
		rules: make([]Rule, 0),
	}
}

// AddRule adds a custom rule function to the validator
func (v *Validator) AddRule(rule Rule) {
	v.rules = append(v.rules, rule)
}

// Validate runs all rules against the value and returns an array of error messages
func (v *Validator) Validate(value interface{}) []string {
	var errors []string

	for _, rule := range v.rules {
		if errorMsg := rule(value); errorMsg != "" {
			errors = append(errors, errorMsg)
		}
	}

	return errors
}
