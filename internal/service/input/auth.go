package input

// Register register
type Register struct {
	Account         string `form:"account,omitempty" binding:"min=3,max=50" json:"account,omitempty"`
	Password        string `form:"password,omitempty" binding:"min=8,max=50" json:"password,omitempty"`
	ConfirmPassword string `form:"confirm_password,omitempty" binding:"min=8,max=50" json:"confirm_password,omitempty"`
	DisplayName     string `form:"display_name,omitempty" binding:"min=3,max=20" json:"display_name,omitempty"`
}

// Login login user
type Login struct {
	Account  string `json:"account,omitempty" form:"account" binding:"min=3,max=50"`
	Password string `json:"password,omitempty" form:"password" binding:"min=8,max=50"`
}
