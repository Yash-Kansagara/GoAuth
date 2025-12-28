package models

type Signup struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Login struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type DBUserRow struct {
	Password string `json:"password"`
	UserId   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdatePasswordReq struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
