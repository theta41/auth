package auth

import "testing"

func TestHashGeneration(t *testing.T) {
	exp := "65e84be33532fb784c48129675f9eff3a682b27168c0ea744b2cf58ee02337c5"
	got := GetHash("qwerty")

	if exp != got {
		t.Errorf("Error while generating hash.\nExpected: %s\nGot: %s", exp, got)
	}
}
