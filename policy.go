package gostack

import (
	"fmt"
	"reflect"
)

type policyRegistry struct {
	policies map[reflect.Type]any
}

func newPolicyRegistry() *policyRegistry {
	return &policyRegistry{policies: make(map[reflect.Type]any)}
}

// Policy registers a policy for a model type (pointer to struct).
func (a *App) Policy(modelExample any, policy any) {
	t := reflect.TypeOf(modelExample)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	a.policyReg.policies[t] = policy
}

func (p *policyRegistry) authorize(c *Context, action string, resource any) error {
	if resource == nil {
		return fmt.Errorf("forbidden")
	}
	t := reflect.TypeOf(resource)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pol, ok := p.policies[t]
	if !ok {
		return nil
	}
	user := c.User()
	if user == nil {
		return fmt.Errorf("unauthorized")
	}
	methodName := stringsTitle(action)
	mv := reflect.ValueOf(pol)
	m := mv.MethodByName(methodName)
	if !m.IsValid() {
		return nil
	}
	out := m.Call([]reflect.Value{reflect.ValueOf(user), reflect.ValueOf(resource)})
	if len(out) > 0 && out[0].Type().Kind() == reflect.Bool && !out[0].Bool() {
		return fmt.Errorf("forbidden")
	}
	return nil
}

func stringsTitle(action string) string {
	if action == "" {
		return ""
	}
	a := []rune(action)
	if a[0] >= 'a' && a[0] <= 'z' {
		a[0] -= 'a' - 'A'
	}
	return string(a)
}
