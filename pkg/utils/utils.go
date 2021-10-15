package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func GetOffsetFromHeader(h http.Header) (int64, error) {
	byteRange := h.Get("range")
	if byteRange == "" {
		return 0, nil
	}
	if len(byteRange) < 7 {
		return 0, fmt.Errorf("get range header error, range=%s", byteRange)
	}
	if byteRange[:6] != "bytes=" {
		return 0, fmt.Errorf("get range header error, range=%s", byteRange)
	}
	bytePos := strings.Split(byteRange[6:], "-")
	return strconv.ParseInt(bytePos[0], 0, 64)
}

func GetHashFromHeader(h http.Header) string {
	digest := h.Get("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}

func GetHashFromHeaderValue(value string) string {
	if len(value) < 9 {
		return ""
	}
	if value[:8] != "SHA-256=" {
		return ""
	}
	return value[8:]
}

func GetSizeFromHeader(h http.Header) (int64, error) {
	return strconv.ParseInt(h.Get("content-length"), 0, 64)
}

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	io.Copy(h, r)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
