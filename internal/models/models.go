package models

import "database/sql"

type Signup struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Login struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

type DBUserDataRow struct {
	Password string `json:"password"`
	UserId   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DBUserRow struct {
	Password string `json:"password"`
	UserId   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type DBResetPasswordData struct {
	DBUserRow
	ResetPasswordToken       sql.NullString `json:"reset_password_token"`
	ResetPasswordTokenExpiry sql.NullTime   `json:"reset_password_token_expiry"`
}

type UpdatePasswordReq struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type ResetPasswordReq struct {
	NewPassword string `json:"newPassword"`
}

type ForgotPasswordReq struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
}
