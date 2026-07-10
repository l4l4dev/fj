package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type errorCategory uint8

const (
	categoryInternal errorCategory = iota
	categoryValidation
	categoryAuthentication
	categoryRemote
)

type commandError struct {
	category  errorCategory
	operation string
	cause     error
	message   string
}

func newCommandError(category errorCategory, operation string, cause error) error {
	return commandError{
		category:  category,
		operation: operation,
		cause:     cause,
	}
}

func newCommandErrorWithMessage(category errorCategory, operation, message string, cause error) error {
	return commandError{category: category, operation: operation, message: message, cause: cause}
}

func (err commandError) Error() string {
	if err.message != "" {
		return fmt.Sprintf("%s: %s", err.operation, err.message)
	}
	return fmt.Sprintf("%s: %s", err.operation, err.category.message())
}

func (err commandError) Unwrap() error {
	return err.cause
}

func (category errorCategory) message() string {
	switch category {
	case categoryValidation:
		return "invalid input"
	case categoryAuthentication:
		return "authentication failed"
	case categoryRemote:
		return "remote operation failed"
	default:
		return "internal error"
	}
}

func (category errorCategory) exitCode() int {
	switch category {
	case categoryValidation:
		return 2
	case categoryAuthentication:
		return 3
	case categoryRemote:
		return 4
	default:
		return 1
	}
}

func Execute(command *cobra.Command, args []string) int {
	command.SetArgs(args)

	if _, _, err := command.Find(args); err != nil {
		return presentError(command, newCommandError(categoryValidation, "execute command", err))
	}
	if err := command.Execute(); err != nil {
		return presentError(command, err)
	}
	return 0
}

func presentError(command *cobra.Command, err error) int {
	var classified commandError
	if !errors.As(err, &classified) {
		category := categoryInternal
		if strings.HasPrefix(err.Error(), "unknown command ") {
			category = categoryValidation
		}
		classified = commandError{
			category:  category,
			operation: "execute command",
			cause:     err,
		}
	}

	command.PrintErrln("Error:", classified.Error())
	return classified.category.exitCode()
}
