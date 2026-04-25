package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewEvaluator creates a new evaluator
func TestNewEvaluator(t *testing.T) {
	e := NewEvaluator()

	require.NotNil(t, e)
	assert.Equal(t, 0, len(e.policies))
}

// TestAddPolicy adds policies to evaluator
func TestAddPolicy(t *testing.T) {
	e := NewEvaluator()

	policy := Policy{
		Name: "test",
		Condition: func(s, o, env Attributes) bool {
			return true
		},
	}

	e.AddPolicy(policy)

	assert.Equal(t, 1, len(e.policies))
}

// TestEvaluateAllowPolicy tests evaluation with ALLOW policy
func TestEvaluateAllowPolicy(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "allow-all",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			return true
		},
	})

	result := e.Evaluate(nil, nil, nil)

	assert.True(t, result)
}

// TestEvaluateDenyPolicy tests evaluation with DENY policy
func TestEvaluateDenyPolicy(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "deny-all",
		Effect: EffectDeny,
		Condition: func(s, o, env Attributes) bool {
			return true
		},
	})

	result := e.Evaluate(nil, nil, nil)

	assert.False(t, result)
}

// TestEvaluateDenyOverridesAllow tests that DENY overrides ALLOW
func TestEvaluateDenyOverridesAllow(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "allow",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			return true
		},
	})

	e.AddPolicy(Policy{
		Name:   "deny",
		Effect: EffectDeny,
		Condition: func(s, o, env Attributes) bool {
			return true
		},
	})

	result := e.Evaluate(nil, nil, nil)

	assert.False(t, result, "DENY should override ALLOW")
}

// TestEvaluateNoMatchingPolicy tests when no policies match
func TestEvaluateNoMatchingPolicy(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "non-matching",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			return false
		},
	})

	result := e.Evaluate(nil, nil, nil)

	assert.False(t, result, "should deny when no policies match")
}

// TestEvaluateSubjectAttributes tests evaluation with subject attributes
func TestEvaluateSubjectAttributes(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "admin-only",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			role, ok := s["role"]
			return ok && role == "admin"
		},
	})

	// Admin should be allowed
	adminSubject := Attributes{"role": "admin"}
	result := e.Evaluate(adminSubject, nil, nil)
	assert.True(t, result)

	// Non-admin should be denied
	userSubject := Attributes{"role": "user"}
	result = e.Evaluate(userSubject, nil, nil)
	assert.False(t, result)
}

// TestEvaluateObjectAttributes tests evaluation with object attributes
func TestEvaluateObjectAttributes(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "can-read-public",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			if visibility, ok := o["visibility"]; ok {
				return visibility == "public"
			}
			return false
		},
	})

	// Public object should be allowed
	publicObject := Attributes{"visibility": "public"}
	result := e.Evaluate(nil, publicObject, nil)
	assert.True(t, result)

	// Private object should be denied
	privateObject := Attributes{"visibility": "private"}
	result = e.Evaluate(nil, privateObject, nil)
	assert.False(t, result)
}

// TestEvaluateEnvironmentAttributes tests evaluation with environment attributes
func TestEvaluateEnvironmentAttributes(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "only-during-business-hours",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			if hour, ok := env["hour"]; ok {
				h := hour.(int)
				return h >= 9 && h < 17
			}
			return false
		},
	})

	// During business hours (10am)
	result := e.Evaluate(nil, nil, Attributes{"hour": 10})
	assert.True(t, result)

	// Outside business hours (8pm)
	result = e.Evaluate(nil, nil, Attributes{"hour": 20})
	assert.False(t, result)
}

// TestEvaluateComplexPolicy tests complex policy with multiple attribute checks
func TestEvaluateComplexPolicy(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "owner-can-edit",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			// Subject must be owner
			subjectID, subjectOk := s["user_id"]
			ownerID, ownerOk := o["owner_id"]
			if !subjectOk || !ownerOk {
				return false
			}

			// Resource must be editable
			editable, editableOk := o["editable"]
			if !editableOk {
				return false
			}

			return subjectID == ownerID && editable == true
		},
	})

	// Owner editing editable resource
	subject := Attributes{"user_id": "user-1"}
	object := Attributes{"owner_id": "user-1", "editable": true}
	result := e.Evaluate(subject, object, nil)
	assert.True(t, result)

	// Non-owner trying to edit
	subject = Attributes{"user_id": "user-2"}
	result = e.Evaluate(subject, object, nil)
	assert.False(t, result)

	// Owner trying to edit non-editable
	subject = Attributes{"user_id": "user-1"}
	object = Attributes{"owner_id": "user-1", "editable": false}
	result = e.Evaluate(subject, object, nil)
	assert.False(t, result)
}

// TestEvaluateMultiplePolicies tests multiple policy evaluation
func TestEvaluateMultiplePolicies(t *testing.T) {
	e := NewEvaluator()

	// Policy 1: Allow admins
	e.AddPolicy(Policy{
		Name:   "admin-allow",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			role, ok := s["role"]
			return ok && role == "admin"
		},
	})

	// Policy 2: Deny if suspended
	e.AddPolicy(Policy{
		Name:   "suspended-deny",
		Effect: EffectDeny,
		Condition: func(s, o, env Attributes) bool {
			suspended, ok := s["suspended"]
			return ok && suspended == true
		},
	})

	// Admin and not suspended should be allowed
	result := e.Evaluate(
		Attributes{"role": "admin", "suspended": false},
		nil,
		nil,
	)
	assert.True(t, result)

	// Admin but suspended should be denied (DENY overrides)
	result = e.Evaluate(
		Attributes{"role": "admin", "suspended": true},
		nil,
		nil,
	)
	assert.False(t, result)

	// User and not suspended should be denied (no matching ALLOW)
	result = e.Evaluate(
		Attributes{"role": "user", "suspended": false},
		nil,
		nil,
	)
	assert.False(t, result)
}

// TestEvaluateRoleBasedAccess tests role-based access control
func TestEvaluateRoleBasedAccess(t *testing.T) {
	e := NewEvaluator()

	// Admins can do anything
	e.AddPolicy(Policy{
		Name:   "admin-all-access",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			role, ok := s["role"]
			return ok && role == "admin"
		},
	})

	// Moderators can moderate
	e.AddPolicy(Policy{
		Name:   "moderator-moderate",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			role, subjectOk := s["role"]
			action, objectOk := o["action"]
			return subjectOk && role == "moderator" && objectOk && action == "moderate"
		},
	})

	// Users can only read
	e.AddPolicy(Policy{
		Name:   "user-read",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			role, subjectOk := s["role"]
			action, objectOk := o["action"]
			return subjectOk && role == "user" && objectOk && action == "read"
		},
	})

	// Admin can do anything
	result := e.Evaluate(
		Attributes{"role": "admin"},
		Attributes{"action": "delete"},
		nil,
	)
	assert.True(t, result)

	// Moderator can moderate
	result = e.Evaluate(
		Attributes{"role": "moderator"},
		Attributes{"action": "moderate"},
		nil,
	)
	assert.True(t, result)

	// Moderator cannot delete
	result = e.Evaluate(
		Attributes{"role": "moderator"},
		Attributes{"action": "delete"},
		nil,
	)
	assert.False(t, result)

	// User can read
	result = e.Evaluate(
		Attributes{"role": "user"},
		Attributes{"action": "read"},
		nil,
	)
	assert.True(t, result)

	// User cannot moderate
	result = e.Evaluate(
		Attributes{"role": "user"},
		Attributes{"action": "moderate"},
		nil,
	)
	assert.False(t, result)
}

// TestEvaluateDataOwnershipPolicy tests data ownership policy
func TestEvaluateDataOwnershipPolicy(t *testing.T) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "owner-access",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			userID, userOk := s["user_id"]
			ownerID, ownerOk := o["owner_id"]
			return userOk && ownerOk && userID == ownerID
		},
	})

	e.AddPolicy(Policy{
		Name:   "admin-bypass",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			role, ok := s["role"]
			return ok && role == "admin"
		},
	})

	// Owner accessing own data
	result := e.Evaluate(
		Attributes{"user_id": "user-1"},
		Attributes{"owner_id": "user-1"},
		nil,
	)
	assert.True(t, result)

	// Non-owner accessing data
	result = e.Evaluate(
		Attributes{"user_id": "user-1"},
		Attributes{"owner_id": "user-2"},
		nil,
	)
	assert.False(t, result)

	// Admin accessing any data
	result = e.Evaluate(
		Attributes{"user_id": "admin-1", "role": "admin"},
		Attributes{"owner_id": "user-2"},
		nil,
	)
	assert.True(t, result)
}

// BenchmarkEvaluate benchmarks policy evaluation
func BenchmarkEvaluate(b *testing.B) {
	e := NewEvaluator()

	e.AddPolicy(Policy{
		Name:   "test",
		Effect: EffectAllow,
		Condition: func(s, o, env Attributes) bool {
			role, ok := s["role"]
			return ok && role == "admin"
		},
	})

	subject := Attributes{"role": "admin"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Evaluate(subject, nil, nil)
	}
}

// BenchmarkComplexEvaluation benchmarks complex policy evaluation
func BenchmarkComplexEvaluation(b *testing.B) {
	e := NewEvaluator()

	// Multiple policies
	for i := 0; i < 10; i++ {
		e.AddPolicy(Policy{
			Name:   "policy",
			Effect: EffectAllow,
			Condition: func(s, o, env Attributes) bool {
				role, ok := s["role"]
				resource, resOk := o["resource"]
				return ok && resOk && role == "admin" && resource != ""
			},
		})
	}

	subject := Attributes{"role": "admin", "permissions": []string{"read", "write"}}
	object := Attributes{"resource": "users", "owner_id": "user-1"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e.Evaluate(subject, object, nil)
	}
}
