package hash

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/tony-spark/metrico/internal"
	"github.com/tony-spark/metrico/internal/dto"
)

type Sha256Keyed struct {
	key string
}

func NewSha256Keyed(key string) *Sha256Keyed {
	return &Sha256Keyed{
		key: key,
	}
}

func (s Sha256Keyed) Hash(m dto.Metric) (string, error) {
	hash, err := hashBin(m, s.key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash), nil
}

func (s Sha256Keyed) Check(m dto.Metric) (bool, error) {
	calc, err := hashBin(m, s.key)
	if err != nil {
		return false, err
	}
	orig, err := hex.DecodeString(m.Hash)
	if err != nil {
		return false, err
	}
	return bytes.Equal(calc, orig), nil
}

func hashBin(m dto.Metric, key string) ([]byte, error) {
	var repr string
	switch m.MType {
	case internal.COUNTER:
		repr = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	case internal.GAUGE:
		repr = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
	default:
		return nil, fmt.Errorf("coulnd not calculate hash for unknown metric type: %s", m.MType)
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(repr))
	return h.Sum(nil), nil
}
