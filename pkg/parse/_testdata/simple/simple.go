package simple

// +Foo=true
// +Bar=123
type User struct {
	// +something:id=true
	ID          string  `json:"id"`
	DisplayName *string `json:"display_name"`
	Email       string  `json:"email"`
	Age         int     `json:"age"`
}
