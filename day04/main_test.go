package main

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	tt := []struct {
		pass  *password
		valid bool
	}{
		{newPassword(123456), false},
		{newPassword(123446), true},
		{newPassword(111111), true},
		{newPassword(223450), false},
		{newPassword(123789), false},
	}

	for _, tc := range tt {
		t.Run(tc.pass.String(), func(t *testing.T) {
			if valid := tc.pass.isValid(); tc.valid != valid {
				t.Fatalf("expected password validity to be %v; got %v", tc.valid, valid)
			}
		})
	}
}

func TestNext(t *testing.T) {
	last := newPassword(999_999)

	count := 0
	for p := newPassword(0); true; p.next() {
		count++
		if p.equals(last) {
			break
		}
	}
	if count != 1_000_000 {
		t.Fatalf("expected to count 1,000,000 passwords, counted %d", count)
	}
}
