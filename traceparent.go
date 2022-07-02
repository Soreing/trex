package trex

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

// Generates N*2 random hexadecimal digits as a string
// N is the number of bytes of randomness used
func GenerateRadomHexString(n int) (string, error) {
	buff := make([]byte, n)
	if _, err := rand.Read(buff); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(buff), nil
	}
}

// Decodes a traceparent header into sections
// Returns version, transaction id, parent id and flag values.
func DecodeTraceparent(traceparent string) (string, string, string, string, error) {
	// Fast fail for common case of empty string
	if traceparent == "" {
		return "", "", "", "", fmt.Errorf("traceparent is empty string")
	}

	hexfmt, err := regexp.Compile("^[0-9A-Fa-f]*$")
	vals := strings.Split(traceparent, "-")

	if len(vals) == 4 {
		ver, tid, pid, flg := vals[0], vals[1], vals[2], vals[3]
		if !hexfmt.MatchString(ver) || len(ver) != 2 {
			err = fmt.Errorf("invalid traceparent version")
		} else if !hexfmt.MatchString(pid) || len(pid) != 16 {
			err = fmt.Errorf("invalid traceparent parent id")
		} else if !hexfmt.MatchString(flg) || len(flg) != 2 {
			err = fmt.Errorf("invalid traceparent flag")
		} else if !hexfmt.MatchString(tid) || len(tid) != 32 {
			err = fmt.Errorf("invalid traceparent trace id")
		} else if tid == "00000000000000000000000000000000" {
			err = fmt.Errorf("traceparent trace id value is zero")
		} else {
			return ver, tid, pid, flg, nil
		}
	} else {
		err = fmt.Errorf("invalid traceparent trace id")
	}

	return "", "", "", "", err
}
