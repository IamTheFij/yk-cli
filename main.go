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
		slog.Error("failed reading password from input")

		return err
	}

	// get password
	key := s.DeriveKey(string(bytePassword))

	ok, err := oath.Validate(s, key)
	if err != nil {
		return err
	}

	if !ok {
		return errFailedValidation
	}

	return keyring.Set(
		serviceName,
		s.DeviceID(),
		string(key),
	)
}

func getPassword(s *ykoath.Select) (string, error) {
	return keyring.Get(serviceName, s.DeviceID())
}

func main() {
	flag.BoolVar(&slog.DebugLevel, "debug", false, "enable debug logging")
	showVersion := flag.Bool("version", false, "print version and exit")
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
	if flag.Arg(0) == "set-password" {
		err = setPassword(s)
		slog.FatalOnErr(err, "failed to save password")

		return
	}

	// If required, authenticate with password from keychain
	if s.Challenge != nil {
		passKey, err := getPassword(s)
		slog.FatalOnErr(err, "failed retrieving password key")

		ok, err := oath.Validate(s, []byte(passKey))
		slog.FatalOnErr(err, "validation failed")

		if !ok {
			slog.Fatal("failed validation, password is incorrect")
		}
	}

	if flag.Arg(0) == "list" {
		// List names only
		names, err := oath.List()
		slog.FatalOnErr(err, "failed to list names")

		for _, name := range names {
			fmt.Println(name.Name)
		}
	} else {
		name := flag.Arg(0)
		code, err := oath.CalculateOne(name)
		slog.FatalOnErr(err, "failed to retrieve credential")

		fmt.Println(code)
	}
}
