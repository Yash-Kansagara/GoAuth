package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	utils "github.com/Yash-Kansagara/GoAuth/internal/Utils"
	"github.com/Yash-Kansagara/GoAuth/internal/db"
	"github.com/Yash-Kansagara/GoAuth/internal/models"
)

func RegisterSignupHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /signup", PostSignupHandler)
	mux.HandleFunc("POST /login", PostLoginHandler)
	mux.HandleFunc("POST /logout", PostLogoutHandler)
	mux.HandleFunc("POST /updatePassword", PostUpdatePasswordHandler)
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
	rowResp := db.QueryRow("SELECT * FROM users WHERE username = ? OR email = ?", usernameOrEmail, usernameOrEmail)
	row := models.DBUserRow{}
	err := rowResp.Scan(&row.UserId, &row.Username, &row.Email, &row.Password)
	return row, err
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
