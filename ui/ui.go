package ui

import (
	"crypto/x509"
	"errors"
	"fmt"
	"strings"
)

type uiface func(certs []*x509.Certificate, cfg *Config) error

var (
	uiOptions = map[string]uiface{
		"cli": showCertsOnCli,
		"web": showCertsOnWeb,
	}
)

// UI - what technology to display results on
func DefaultUI() string {
	return "cli"
}
func GetUIs() []string {
	return []string{DefaultFormat(), "web"}
}

// Formats - how the data is displayed on the UI
func DefaultFormat() string {
	return "table"
}
func GetFormats() []string {
	out := make([]string, 0)
	for k, _ := range outputFormats {
		out = append(out, k)
	}
	return out
}

type Config struct {
	// If we should only show the certificate count, rather than each one
	Count bool

	// What format to print certificates in, formats are defined in ../main.go and
	// checked in print.go
	Format string

	// Which user interface to show users, e.g. cli or web
	// Default (and possible) value(s) can be found in the ui package
	UI string
}

func ListCertificates(certs []*x509.Certificate, cfg *Config) error {
	if cfg.Count { // ignore any cfg.UI setting
		fmt.Printf("%d\n", len(certs))
		return nil
	}

	// Show something meaningful if there's nothing otherwise
	if len(certs) == 0 {
		return errors.New("No certififcates to display")
	}

	fn, ok := uiOptions[strings.ToLower(cfg.UI)]
	if !ok {
		return fmt.Errorf("Unknown ui '%s'", cfg.UI)
	}
	return fn(certs, cfg)
}
