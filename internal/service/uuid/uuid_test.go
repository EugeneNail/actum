package uuid

import "testing"

func TestVersion(t *testing.T) {
	for i := 0; i < 10000; i++ {
		uuid := New()
		if uuid[14] != '4' {
			t.Errorf("Expected the 13th character of uuid %s to be 4", uuid)
		}
	}
}

func TestVariant(t *testing.T) {
	for i := 0; i < 10000; i++ {
		uuid := New()
		variant := uuid[19]

		hasInvalidVariant := !(variant == '8' || variant == '9' || variant == 'a' || variant == 'b' || variant == 'c' || variant == 'd')
		if hasInvalidVariant {
			t.Errorf("Expected the 17th character of uuid %s to be 8, 9, a, b, c or d", uuid)
		}
	}
}

func BenchmarkUuid(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}
