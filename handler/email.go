package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func EmailDump(c echo.Context) error {
	filePath := "./emails.txt"

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Set the response header
	c.Response().Header().Set(echo.HeaderContentType, http.DetectContentType([]byte(filePath)))
	c.Response().
		Header().
		Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=\"%s\"", filePath))

	// Copy the file to the response writer
	_, err = io.Copy(c.Response().Writer, file)
	if err != nil {
		return err
	}

	return nil
}
