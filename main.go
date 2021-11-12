package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"syscall"

	"git.iamthefij.com/iamthefij/slog"
	"github.com/yawn/ykoath"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

var (
	oath                *ykoath.OATH
	serviceName         = "com.iamthefij.yk-cli"
	version             = "dev"
	errFailedValidation = errors.New("failed validation, password may be incorrect")
)

func setPassword(s *ykoath.Select) error {
	fmt.Print("Enter Password: ")

	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return fmt.Errorf("failed reading password from input: %w", err)
	}

	password := string(bytePassword)

	// get key
	key := s.DeriveKey(password)

	// verify password is correct with a validate call
	ok, err := oath.Validate(s, key)
	if err != nil {
		return fmt.Errorf("error in validate: %w", err)
	}

	if !ok {
		return errFailedValidation
	}

	err = keyring.Set(
		serviceName,
		s.DeviceID(),
		password,
	)
	if err != nil {
		return fmt.Errorf("error saving password in keyring: %w", err)
	}

	return nil
}

func getPasskey(s *ykoath.Select) ([]byte, error) {
	password, err := keyring.Get(serviceName, s.DeviceID())
	if err != nil {
		return nil, fmt.Errorf("error retrieving key from keyring: %w", err)
	}

	return s.DeriveKey(password), nil
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n  %s [target site]\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "Prints TOTP code for provided target site. If no site is provided, %s will list all sites.\n\n", os.Args[0])
	fmt.Fprint(flag.CommandLine.Output(), "If a touch is required for your code, the command will not return until the your key is touched.\n\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.BoolVar(&slog.DebugLevel, "debug", false, "enable debug logging")
	showVersion := flag.Bool("version", false, "print version and exit")
	shouldSetPassword := flag.Bool("set-password", false, "prompt for key password and store in system keychain")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version %s\n", version)
		os.Exit(0)
	}

	var err error

	oath, err = ykoath.New()
	slog.FatalOnErr(err, "failed to initialize new oath")

	defer oath.Close()

	if slog.DebugLevel {
		oath.Debug = slog.Debug
	}

	// Select oath to begin
	s, err := oath.Select()
	slog.FatalOnErr(err, "failed to select oath")

	// Check to see if we are trying to set a password
	if *shouldSetPassword {
		err = setPassword(s)
		slog.FatalOnErr(err, "failed to save password")

		return
	}

	// If required, authenticate with password from keychain
	if s.Challenge != nil {
		passKey, err := getPasskey(s)
		slog.FatalOnErr(err, "failed retrieving password key")

		ok, err := oath.Validate(s, passKey)
		slog.FatalOnErr(err, "validation failed")

		if !ok {
			slog.Fatal("failed validation, password is incorrect")
		}
	} else {
		slog.Debug("no challenge required")
	}

	if flag.NArg() == 0 || flag.Arg(0) == "list" {
		// List names only
		names, err := oath.List()
		slog.FatalOnErr(err, "failed to list names")

		for _, name := range names {
			fmt.Println(name.Name)
		}
	} else {
		name := flag.Arg(0)
		code, err := oath.Calculate(name, func(string) error {
			fmt.Fprintf(os.Stderr, "Touch key to generate code for %s...", name)
			return nil
		})
		slog.FatalOnErr(err, "failed to retrieve credential")

		fmt.Println(code)
	}
}
