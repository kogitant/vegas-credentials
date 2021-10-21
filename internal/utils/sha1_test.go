package utils

import "testing"

func TestGenerateSHA1(t *testing.T) {
	input := "foobar"

	// want generated with https://passwordsgenerator.net/sha1-hash-generator/
	want := "8843d7f92416211de9ebb963ff4ce28125932878"

	output := GenerateSHA1(input)
	if output != want {
		t.Fatalf(`generateSha1Hash("%s") = %q, want match for %#q`, input, output, want)
	}
}