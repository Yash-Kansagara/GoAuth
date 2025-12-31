package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	utils "github.com/Yash-Kansagara/GoAuth/internal/Utils"
	"github.com/Yash-Kansagara/GoAuth/internal/db"
	"github.com/Yash-Kansagara/GoAuth/internal/models"
	"gopkg.in/gomail.v2"
)

func RegisterSignupHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /signup", PostSignupHandler)
	mux.HandleFunc("POST /login", PostLoginHandler)
	mux.HandleFunc("POST /logout", PostLogoutHandler)
	mux.HandleFunc("POST /updatePassword", PostUpdatePasswordHandler)
	mux.HandleFunc("POST /forgotPassword", PostForgotPasswordHandler)
}

func PostSignupHandler(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if utils.WriteIfError(w, err, "Error reading request", http.StatusInternalServerError) {
		return
	}

	signupReq := models.Signup{}
	json.Unmarshal(bodyBytes, &signupReq)
	fmt.Println(signupReq)

	signupReq.Password = utils.GetHash(signupReq.Password)

	db := db.GetDB()
	stmt, err := db.Prepare("INSERT INTO users (username, email, password) VALUES (?,?,?)")
	if utils.WriteIfError(w, err, "Error signing up", http.StatusInternalServerError) {
		return
	}
	_, err = stmt.Exec(signupReq.Username, signupReq.Email, signupReq.Password)
	if utils.WriteIfError(w, err, "Error signing up", http.StatusInternalServerError) {
		return
	}

	w.Header().Set("Content-type", "Application/text")
	w.Write([]byte("Signup successful"))

}

func PostLoginHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if utils.WriteIfError(w, err, "Error reading request", http.StatusInternalServerError) {
		return
	}

	loginReq := models.Login{}
	json.Unmarshal(bodyBytes, &loginReq)

	row, err := getCurrentUserFromDB(loginReq.UsernameOrEmail)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "Application/json")
		w.Write([]byte("{\"status\":\"failed\",\"error\":\"user not found\"}"))
		return
	}
	if utils.WriteIfError(w, err, "Error fetching user data", http.StatusInternalServerError) {
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	if isPasswordValid(row.Password, loginReq.Password) {
		// generate JWT token to send with login success

		token, err := utils.SignToken(row.Username, row.UserId, row.Email)
		if err == nil {
			setJWTCookie(w, token)
		}
		w.Write([]byte("{\"status\":\"success\",\"error\":null}"))
	} else {
		setJWTCookie(w, "")
		w.Write([]byte("{\"status\":\"failed\",\"error\":\"Incorrect passord\"}"))
	}
}

func PostLogoutHandler(w http.ResponseWriter, r *http.Request) {
	setJWTCookie(w, "")
	w.Header().Set("Content-Type", "Application/json")
	w.Write([]byte("{\"status\":\"success\",\"error\":null}"))
}

func PostUpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if utils.WriteIfError(w, err, "Error reading request", http.StatusInternalServerError) {
		return
	}

	updatePasswordReq := models.UpdatePasswordReq{}
	err = json.Unmarshal(bodyBytes, &updatePasswordReq)
	if utils.WriteIfError(w, err, "Error reading request", http.StatusInternalServerError) {
		return
	}

	row, err := getCurrentUserFromDB(updatePasswordReq.UsernameOrEmail)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "Application/json")
		w.Write([]byte("{\"status\":\"failed\",\"error\":\"user not found\"}"))
		return
	}
	if utils.WriteIfError(w, err, "Error fetching user data", http.StatusInternalServerError) {
		return
	}

	w.Header().Set("Content-Type", "Application/json")
	if isPasswordValid(row.Password, updatePasswordReq.CurrentPassword) {
		err = updatePasswordInDB(row.UserId, updatePasswordReq.NewPassword)
		if utils.WriteIfError(w, err, "Error Updating password", http.StatusInternalServerError) {
			return
		}
		w.Write([]byte("{\"status\":\"success\",\"error\":null}"))
	} else {
		w.Write([]byte("{\"status\":\"failed\",\"error\":\"Incorrect passowrd\"}"))
	}

}

func isPasswordValid(hash string, password string) bool {
	parts := strings.Split(hash, ":")
	if utils.GetHashWithSalt(parts[0], password) == parts[1] {
		return true
	}

	return false
}

func getCurrentUserFromDB(usernameOrEmail string) (models.DBUserRow, error) {
	db := db.GetDB()
	rowResp := db.QueryRow("SELECT userid, username, email, password FROM users WHERE username = ? OR email = ?", usernameOrEmail, usernameOrEmail)
	row := models.DBUserRow{}
	err := rowResp.Scan(&row.UserId, &row.Username, &row.Email, &row.Password)
	return row, err
}

func getUserResetPasswordData(usernameOrEmail string) (models.DBResetPasswordData, error) {
	db := db.GetDB()
	rowResp := db.QueryRow("SELECT * FROM users WHERE username = ? OR email = ?", usernameOrEmail, usernameOrEmail)
	row := models.DBResetPasswordData{}
	err := rowResp.Scan(&row.UserId, &row.Username, &row.Email, &row.Password, &row.ResetPasswordToken, &row.ResetPasswordTokenExpiry)
	return row, err
}

func setUserResetPasswordData(userid string, token string, expiray time.Time) (*sql.Tx, error) {
	sqldb := db.GetDB()
	tx, err := sqldb.Begin()
	if err != nil {
		return nil, err
	}
	rowResp, err := tx.Exec("UPDATE users SET reset_password_token = ?,reset_password_token_expiry = ? WHERE userid = ?", token, expiray, userid)
	if err != nil {
		return nil, err
	}
	affected, err := rowResp.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if affected != 1 {
		tx.Rollback()
		err = fmt.Errorf("Incorrect rows affected: ", affected)
		return nil, err
	}
	return tx, nil
}

func updatePasswordInDB(userid string, password string) error {
	passwordHash := utils.GetHash(password)
	database := db.GetDB()

	query, err := database.Prepare("UPDATE users SET password = ? WHERE userid = ?")
	if err != nil {
		return err
	}

	res, err := query.Exec(passwordHash, userid)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rows != 1 {
		return fmt.Errorf("Incorrect Password Update, affected rows = ", rows)
	}

	return nil
}

func setJWTCookie(w http.ResponseWriter, value string) {

	var exp time.Time
	if len(value) == 0 {
		exp = time.Unix(0, 0)
	} else {
		exp = time.Now().Add(24 * time.Hour)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  exp,
	})
}

func PostForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if utils.WriteIfError(w, err, "Error reading request", http.StatusInternalServerError) {
		return
	}

	forgotPasswordReq := models.ForgotPasswordReq{}
	json.Unmarshal(bodyBytes, &forgotPasswordReq)

	row, err := getUserResetPasswordData(forgotPasswordReq.UsernameOrEmail)
	if utils.WriteIfError(w, err, "Server error", http.StatusInternalServerError) {
		return
	}

	if row.ResetPasswordTokenExpiry.Valid && row.ResetPasswordTokenExpiry.Time.After(time.Now()) {
		// not expired
		// resend same mail / error
		fmt.Println("row:", row)
		fmt.Println(row.ResetPasswordTokenExpiry.Time)
		w.Header().Set("Content-Type", "Application/json")
		w.Write([]byte("{\"status\":\"fail\",\"error\":\"mail already sent\"}"))
	} else {
		// expired or do not exist
		// create new and send mail
		randomHash := utils.GetRandomHash()
		resetExpiryEnv := os.Getenv("FORGET_PASS_EXPIRY_DURATION")
		resetExpDuration, err := time.ParseDuration(resetExpiryEnv)
		if err != nil {
			resetExpDuration = time.Duration(10 * time.Minute)
		}

		tx, err := setUserResetPasswordData(row.UserId, randomHash, time.Now().Add(resetExpDuration))
		defer tx.Commit() // commit only if we are able to send the mail
		if utils.WriteIfError(w, err, "Error Updating user data", http.StatusInternalServerError) {
			tx.Rollback()
			return
		}

		randomHash = url.QueryEscape(randomHash)
		resetUrl := utils.GetHostUrl(fmt.Sprintf("/resetPassword?id=%s&token=%s", row.UserId, randomHash))

		fmt.Println("reset url:", resetUrl)
		err = sendPasswordResetMail(row.Email, resetUrl, resetExpDuration)
		if utils.WriteIfError(w, err, "Error sending password reset email", http.StatusInternalServerError) {
			tx.Rollback()
			return
		}

		w.Header().Set("Content-Type", "Application/json")
		w.Write([]byte("{\"status\":\"success\",\"error\":null}"))
	}

}

func sendPasswordResetMail(emailId string, resetUrl string, duration time.Duration) error {
	email := gomail.NewMessage()
	email.SetHeader("From", "goauth@goauth.com")
	email.SetHeader("To", emailId)
	email.SetHeader("Subject", "Password reset")
	email.SetBody("text/plain", fmt.Sprintf("Click this link to reset your password: %s.\n Link valid till %s", resetUrl, duration))

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		fmt.Println("Invalid smtp port, default port 1025 will be used, Error:", err)
		smtpPort = 1025
	}
	dialer := gomail.NewDialer(smtpHost, smtpPort, "", "")
	return dialer.DialAndSend(email)
}
