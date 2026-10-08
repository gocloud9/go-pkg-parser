package embedded

type Parent struct {
}

type ParentInterface interface {
}

type ParentStruct func()

// +Foo=true
// +Bar=123
type Child struct {
	// +Foo=true
	// +Bar=123
	Parent `yaml:",inline"`
	// +Foo=true
	// +Bar=123
	ParentInterface `yaml:",inline"`
}

// +Foo=true
// +Bar=123
type ChildInterface interface {
	// +Foo=true
	// +Bar=123
	Parent
	// +Foo=true
	// +Bar=123
	ParentInterface
}
