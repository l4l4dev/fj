package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	applicationauth "github.com/l4l4dev/fj/internal/application/auth"
	infrastructureauth "github.com/l4l4dev/fj/internal/infrastructure/auth"
	infrastructureconfig "github.com/l4l4dev/fj/internal/infrastructure/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newAuthCommand() *cobra.Command {
	command := &cobra.Command{Use: "auth", Short: "Manage authentication"}
	command.AddCommand(newAuthLoginCommand(nil, nil))
	return command
}

func newAuthLoginCommand(verifier applicationauth.TokenVerifier, store applicationauth.CredentialStore) *cobra.Command {
	var instanceName string
	var tokenStdin bool
	command := &cobra.Command{Use: "login", Short: "Store an access token for a configured instance", Args: cobra.NoArgs, RunE: func(command *cobra.Command, _ []string) error {
		configuration, err := infrastructureconfig.Load()
		if err != nil {
			return newCommandError(categoryValidation, "load configuration", err)
		}
		instance, err := configuration.SelectInstance(instanceName)
		if err != nil {
			return newCommandError(categoryValidation, "select instance", err)
		}

		token, err := readLoginToken(command, tokenStdin, string(instance.Endpoint))
		if err != nil {
			return err
		}

		if verifier == nil {
			verifier = infrastructureauth.NewForgejoVerifier(versionFromContext(command.Context()))
		}
		if store == nil {
			store = infrastructureauth.NewFileStore()
		}
		result, err := applicationauth.NewLoginUseCase(verifier, store).Execute(command.Context(), applicationauth.LoginRequest{Instance: instance, Token: token})
		if err != nil {
			return mapApplicationError(err, "auth login")
		}
		_, err = fmt.Fprintf(command.OutOrStdout(), "Logged in to %s as %s\nToken stored in %s\n", instance.Endpoint, result.Username, result.Path)
		return err
	}}
	command.Flags().StringVar(&instanceName, "instance", "", "configured Forgejo instance profile")
	command.Flags().BoolVar(&tokenStdin, "token-stdin", false, "read the token from the first line of standard input")
	return command
}

func readLoginToken(command *cobra.Command, tokenStdin bool, endpoint string) (string, error) {
	if tokenStdin {
		// A read failure and an empty line are both handled as a missing token
		// by the use case, so the error value is deliberately ignored here.
		line, _ := bufio.NewReader(command.InOrStdin()).ReadString('\n')
		return strings.TrimSpace(line), nil
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", newCommandError(categoryValidation, "auth login", fmt.Errorf("stdin is not a terminal; use --token-stdin"))
	}
	command.PrintErrf("Token for %s: ", endpoint)
	entered, err := term.ReadPassword(int(os.Stdin.Fd()))
	command.PrintErrln()
	if err != nil {
		return "", newCommandError(categoryValidation, "auth login", fmt.Errorf("could not read the token from the terminal"))
	}
	return strings.TrimSpace(string(entered)), nil
}
