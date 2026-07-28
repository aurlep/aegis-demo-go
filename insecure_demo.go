// INTENTIONALLY INSECURE — a target for gosec / Semgrep.
// Not wired into the running server; it exists so the generated pipeline has
// real findings. Every import is stdlib, so the build stays clean. Do not copy.
package main

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"math/rand"
	"os/exec"
)

// Hardcoded credentials — secret scanners (Gitleaks, TruffleHog) should flag these.
const (
	apiToken   = "ghp_012345678901234567890123456789abcdef"
	dbPassword = "SuperSecret123!"
)

// G401/G501: MD5 for a password digest.
func insecureHash(pw string) string {
	sum := md5.Sum([]byte(pw))
	return fmt.Sprintf("%x", sum)
}

// G404: weak randomness for a security token.
func insecureToken() int {
	return rand.Intn(1_000_000)
}

// G204: command injection from untrusted input.
func insecureExec(userInput string) ([]byte, error) {
	return exec.Command("sh", "-c", "ping "+userInput).Output()
}

// G402: TLS certificate verification disabled.
func insecureTLS() *tls.Config {
	return &tls.Config{InsecureSkipVerify: true}
}
