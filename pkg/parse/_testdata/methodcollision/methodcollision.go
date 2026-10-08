package methodcollision

type A struct{}

func (*A) Validate() {
}

type B struct{}

func (*B) Validate() {
}

func Validate() {
}
