package utils

import (
	"fmt"
	"net/http"
)

func WriteIfError(w http.ResponseWriter, err error, reponseMsg string, status int) bool {
	if err != nil {
		fmt.Println(reponseMsg)
		fmt.Println(err.Error())
		http.Error(w, reponseMsg, status)
		return true
	}
	return false
}
