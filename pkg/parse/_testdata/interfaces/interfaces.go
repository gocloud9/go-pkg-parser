package interfaces

// +Foo=true
// +Bar=123
type MyInterface interface {
	// +Foo=true
	// +Bar=123
	DoSomething(input *string, f func([]TestInterface) map[string]TestStruct) (output string, err error)
}

type TestStruct struct{}

type TestInterface interface{}
