package utils

import (
	"fmt"
	"os"
	"strings"
)

func GetHostUrl(path string) string {
	host := os.Getenv("API_HOST")
	port := os.Getenv("API_PORT")
	path = strings.TrimPrefix(path, "/")
	resetUrl := fmt.Sprintf("%s:%s/%s", host, port, path)

	return resetUrl
}
