package handler

import (
	"bufio"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/tommaso-merlini/dream-interpreter/view/components"
)

func AddEmail(c echo.Context) error {
	email := c.FormValue("email")
	file, err := os.OpenFile("emails.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a writer for the file
	writer := bufio.NewWriter(file)

	// Write the new line to the file
	_, err = writer.WriteString(email + " " + time.Now().String() + "\n")
	if err != nil {
		return err
	}

	// Flush the writer to ensure all buffered operations are written to the file
	err = writer.Flush()
	if err != nil {
		return err
	}

	return render(c, components.EmailAdded())
}
