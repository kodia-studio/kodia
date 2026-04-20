// Package policy provides an Attribute-Based Access Control (ABAC) engine for Kodia.
package policy

// Attribute represents a key-value attribute of a subject (user), object (resource), or environment.
type Attributes map[string]interface{}

// Condition defines a function that evaluates attributes to decide access.
type Condition func(subject, object, environment Attributes) bool

// Policy defines a rule with conditions.
type Policy struct {
	Name       string
	Description string
	Effect     Effect
	Condition  Condition
}

// Effect represents the result of a policy evaluation.
type Effect string

const (
	EffectAllow Effect = "ALLOW"
	EffectDeny  Effect = "DENY"
)

// Evaluator evaluates a set of policies against attributes.
type Evaluator struct {
	policies []Policy
}

// NewEvaluator creates a new policy evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{policies: make([]Policy, 0)}
}

// AddPolicy adds a policy to the evaluator.
func (e *Evaluator) AddPolicy(p Policy) {
	e.policies = append(e.policies, p)
}

// Evaluate checks if the subject is allowed to perform the action on the object in the environment.
// Returns true only if at least one ALLOW policy matches and no DENY policies match.
func (e *Evaluator) Evaluate(subject, object, environment Attributes) bool {
	allowed := false
	for _, p := range e.policies {
		if p.Condition(subject, object, environment) {
			if p.Effect == EffectDeny {
				return false // Explicit deny
			}
			if p.Effect == EffectAllow {
				allowed = true
			}
		}
	}
	return allowed
}
