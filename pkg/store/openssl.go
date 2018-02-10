package store

import (
	"crypto/x509"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adamdecaf/cert-manage/pkg/certutil"
	"github.com/adamdecaf/cert-manage/pkg/file"
	"github.com/adamdecaf/cert-manage/pkg/whitelist"
)

var (
	openSSLCertPaths = []string{
		"/etc/ssl/certs/ca-certificates.crt", // Ubuntu
		"/usr/local/etc/openssl/certs", // Darwin/OSX
	}

	openSSLRehashPaths = []string{
		"/usr/bin/c_rehash", // Ubuntu
		"/usr/local/opt/openssl/bin/c_rehash", // Darwin/OSX
	}
)

type opensslStore struct{}

func OpenSSLStore() Store {
	return opensslStore{}
}

func (s opensslStore) Add(certs []*x509.Certificate) error {
	dir, err := s.findCertPath()
	if err != nil {
		return err
	}

	for i := range certs {
		fp := certutil.GetHexSHA256Fingerprint(*certs[i])
		path := filepath.Join(dir, fmt.Sprintf("%s.crt", fp))
		err = certutil.ToFile(path, certs[i:i+1])
		if err != nil {
			return err
		}
	}
	if len(certs) > 0 {
		return s.rehash()
	}
	return nil
}

func (s opensslStore) Backup() error {
	return nil
}

func (s opensslStore) GetInfo() *Info {
	out, err := exec.Command("openssl", "version").CombinedOutput()
	if err != nil {
		return &Info{ // just return something non-nil
			Name: "OpenSSL",
		}
	}

	// 'LibreSSL 2.2.7' or 'OpenSSL 1.0.2g  1 Mar 2016'
	parts := strings.Split(string(out), " ")

	return &Info{
		Name:    strings.TrimSpace(parts[0]),
		Version: strings.TrimSpace(parts[1]),
	}
}

func (s opensslStore) List() ([]*x509.Certificate, error) {
	return nil, nil
}

func (s opensslStore) Remove(wh whitelist.Whitelist) error {
	return nil
}

func (s opensslStore) Restore(where string) error {
	return nil
}

// TODO(adam): What do we do if multiple exist
func (s opensslStore) findCertPath() (string, error) {
	for i := range openSSLCertPaths {
		if file.Exists(openSSLCertPaths[i]) {
			return openSSLCertPaths[i], nil
		}
	}
	return "", errors.New("unable to find openssl cert directory")
}

func (s opensslStore) rehash() error {
	var bin string
	for i := range openSSLRehashPaths {
		if file.Exists(openSSLRehashPaths[i]) {
			bin = openSSLRehashPaths[i]
			break
		}
	}
	if bin == "" {
		return errors.New("unable to find openssl c_rehash binary")
	}

	// run c_rehash
	cmd := exec.Command(bin)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if debug {
			fmt.Printf("store/openssl: Command was: %s\n", strings.Join(cmd.Args, " "))
			fmt.Printf("store/openssl: Output: %q\n", string(out))
		}
		return err
	}
	return nil
}
