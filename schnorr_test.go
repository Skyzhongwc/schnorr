package schnorr

import (
	"math/big"
	"strings"
	"testing"

	"encoding/hex"
)

func TestSign(t *testing.T) {
	for _, test := range testCases {
		if test.d == "" {
			continue
		}

		// given
		d, ok := new(big.Int).SetString(test.d, 16)
		if !ok {
			t.Fatalf("Unexpected error from new(big.Int).SetString(%s, 16)", test.d)
		}

		var m [32]byte

		message, err := hex.DecodeString(test.m)
		if err != nil {
			t.Fatalf("Unexpected error from hex.DecodeString(%s): %v", test.m, err)
		}
		copy(m[:], message)

		// when
		result, err := Sign(d, m)
		if err != nil {
			t.Fatalf("Unexpected error from Sign(%s, %s): %v", test.d, test.m, err)
		}

		observed := hex.EncodeToString(result[:])
		expected := strings.ToLower(test.sig)

		// then
		if observed != expected {
			t.Fatalf("Sign(%s, %s) = %s, want %s", test.d, test.m, observed, expected)
		}
	}
}

func TestVerify(t *testing.T) {
	for _, test := range testCases {
		// given
		var (
			pk  [33]byte
			m   [32]byte
			sig [64]byte
		)

		pubKey, err := hex.DecodeString(test.pk)
		if err != nil {
			t.Fatalf("Unexpected error from hex.DecodeString(%s): %v", test.pk, err)
		}
		copy(pk[:], pubKey)

		message, err := hex.DecodeString(test.m)
		if err != nil {
			t.Fatalf("Unexpected error from hex.DecodeString(%s): %v", test.m, err)
		}
		copy(m[:], message)

		signature, err := hex.DecodeString(test.sig)
		if err != nil {
			t.Fatalf("Unexpected error from hex.DecodeString(%s): %v", test.sig, err)
		}
		copy(sig[:], signature)

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("Unexpected panic from Verify(%s, %s, %s): %v ", test.pk, test.m, test.sig, r)
			}
		}()

		// when
		observed, err := Verify(pk, m, sig)
		if err != nil && test.err == nil {
			t.Fatalf("Unexpected error from Verify(%s, %s, %s): %v", test.pk, test.m, test.sig, err)
		} else if err != nil && err.Error() != test.err.Error() {
			t.Fatalf("Unexpected error from Verify(%s, %s, %s): %v", test.pk, test.m, test.sig, err)
		}

		// then
		if observed != test.result {
			t.Fatalf("Verify(%s, %s, %s) = %v, want %v", test.pk, test.m, test.sig, observed, test.result)
		}
	}
}

func BenchmarkSign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testCases {
			if test.d == "" {
				continue
			}
			var m [32]byte
			d, _ := new(big.Int).SetString(test.d, 16)
			message, _ := hex.DecodeString(test.m)
			copy(m[:], message)
			Sign(d, m)
		}
	}
}

func BenchmarkVerify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, test := range testCases {
			var (
				pk  [33]byte
				m   [32]byte
				sig [64]byte
			)

			pubKey, _ := hex.DecodeString(test.pk)
			message, _ := hex.DecodeString(test.m)
			signature, _ := hex.DecodeString(test.sig)
			copy(pk[:], pubKey)
			copy(m[:], message)
			copy(sig[:], signature)

			Verify(pk, m, sig)
		}
	}
}
