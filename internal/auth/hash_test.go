package auth

import "testing"

func TestHashPassword(t *testing.T) {
	cases := []struct {
		input    string
		expected error
	}{
		{
			input:    "coucou",
			expected: nil,
		},
	}

	for _, c := range cases {
		hash, err := HashPassword(c.input)
		if err != nil {
			t.Errorf("Error :%s", err)
			continue
		}
		err = CheckPasswordHash(c.input, hash)
		if err != nil {
			t.Errorf("Password don't match :%s", err)
			continue
		}
	}
}
