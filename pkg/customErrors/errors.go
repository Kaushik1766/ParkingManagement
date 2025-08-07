package customerrors

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

type UserNotFound struct{}

func (e UserNotFound) Error() string {
	return "user not found"
}

type Unathorized struct{}

func (e Unathorized) Error() string {
	return "user unauthorized"
}

func DisplayError(msg string) {
	log.Printf("Error: %s\n", msg)
	color.Red(msg)
	color.Green("Press Enter to continue...")
	fmt.Scanln()
}
