package builtinsregistry

import "github.com/dborchard/tiny_crdb/pkg/f_sql/e_sem/tree"

var registry = map[string]definition{}

// Subscription is a hook to be called once on all registered builtins.
type Subscription func(name string, props *tree.FunctionProperties, overloads []tree.Overload)

var subscriptions []Subscription

// Register registers a builtin. Intending to be called at init time, it panics
// if a function of the same name has already been registered.
func Register(name string, props *tree.FunctionProperties, overloads []tree.Overload) {
	if _, exists := registry[name]; exists {
		panic("duplicate builtin: " + name)
	}
	registry[name] = definition{
		props:     props,
		overloads: overloads,
	}
}

// GetBuiltinProperties provides low-level access to a built-in function's properties.
// For a better, semantic-rich interface consider using tree.FunctionDefinition
// instead, and resolve function names via ResolvableFunctionReference.Resolve().
func GetBuiltinProperties(name string) (*tree.FunctionProperties, []tree.Overload) {
	def, ok := registry[name]
	if !ok {
		return nil, nil
	}
	return def.props, def.overloads
}

type definition struct {
	props     *tree.FunctionProperties
	overloads []tree.Overload
}

// AddSubscription defines a hook to be called once on all registered builtins.
// Subscriptions should be added at init() time, but are load-order independent:
// if you add a subscription after a function is registered, it will immediately
// be called on that function, while functions that are registered afterwards will
// also trigger the hook.
func AddSubscription(s Subscription) {
	for name, def := range registry {
		s(name, def.props, def.overloads)
	}
	subscriptions = append(subscriptions, s)
}
