package globals

const (
	// +Foo=true
	// +Bar=123
	myConstant            = "test"
	myStringType MyString = "test"
)

type MyString string

// +Foo=true
// +Bar=123
var myFunc = func(
	arg string,
) error {
	return nil
}

// +Foo=true
// +Bar=123
var myFunc2 = []func(string) error{
	func(arg string) error {
		return nil
	},
}
