// Package hash provides implementations for dto.Hasher
package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/tony-spark/metrico/internal/dto"
	"github.com/tony-spark/metrico/internal/model"
)

// Sha256Hmac provides dto.Hasher implementation based on SHA-256
type Sha256Hmac struct {
	key string
}

// NewSha256Hmac creates a new Sha256Hmac with a given key
func NewSha256Hmac(key string) *Sha256Hmac {
	return &Sha256Hmac{
		key: key,
	}
}

func (s Sha256Hmac) Hash(m dto.Metric) (string, error) {
	hash, err := hashBin(m, s.key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash), nil
}

func (s Sha256Hmac) Check(m dto.Metric) (bool, error) {
	calc, err := hashBin(m, s.key)
	if err != nil {
		return false, err
	}
	orig, err := hex.DecodeString(m.Hash)
	if err != nil {
		return false, fmt.Errorf("could not check hash: %w", err)
	}
	return bytes.Equal(calc, orig), nil
}

func hashBin(m dto.Metric, key string) ([]byte, error) {
	var repr string
	switch m.MType {
	case model.COUNTER:
		repr = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	case model.GAUGE:
		repr = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	default:
		return nil, fmt.Errorf("coulnd not calculate hash for unknown metric type: %s", m.MType)
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(repr))
	return h.Sum(nil), nil
}
