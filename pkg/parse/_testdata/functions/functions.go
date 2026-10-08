package functions

// +Foo=true
// +Bar=123
var myFunc = func() {}

var (
	// +Foo=true
	// +Bar=123
	myGroupedFunc = func() {}
)

// +Foo=true
// +Bar=123
func Test1() {
}

// +Foo=true
// +Bar=123
func Test2() error {
	return nil
}

// +Foo=true
// +Bar=123
func Test3(

	// +ArgFoo=true
	// +ArgBar=123
	arg string,

) {
}

// +Foo=true
// +Bar=123
type Field struct{}

// +Foo=true
// +Bar=123
func Test4(arg Field) {
}

// +Foo=true
// +Bar=123
func (Field) Test5(arg Field) {

}

// +Foo=true
// +Bar=123
type Reference struct{}

// +Foo=true
// +Bar=123
func (*Reference) Test6(arg Field) {

}

func Variadic(vArg ...Field) {
}
