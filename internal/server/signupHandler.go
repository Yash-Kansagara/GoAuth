package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	utils "github.com/Yash-Kansagara/GoAuth/internal/Utils"
	"github.com/Yash-Kansagara/GoAuth/internal/db"
	"github.com/Yash-Kansagara/GoAuth/internal/models"
)

func RegisterSignupHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /signup", PostSignupHandler)
	mux.HandleFunc("POST /login", PostLoginHandler)
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

	db := db.GetDB()
	rowResp := db.QueryRow("SELECT (password) FROM users WHERE username = ? OR email = ?", loginReq.UsernameOrEmail, loginReq.UsernameOrEmail)
	row := models.DBPasswordRow{}
	err = rowResp.Scan(&row.Password)
	if err == sql.ErrNoRows {
		w.Header().Set("Content-Type", "Application/json")
		w.Write([]byte("{\"status\":\"failed\",\"error\":\"user not found\"}"))
		return
	}
	if utils.WriteIfError(w, err, "Error fetching user data", http.StatusInternalServerError) {
		return
	}
	parts := strings.Split(row.Password, ":")
	w.Header().Set("Content-Type", "Application/json")
	if utils.GetHashWithSalt(parts[0], loginReq.Password) == parts[1] {
		w.Write([]byte("{\"status\":\"success\",\"error\":null}"))
	} else {
		w.Write([]byte("{\"status\":\"failed\",\"error\":\"wrong passord\"}"))
	}
}
