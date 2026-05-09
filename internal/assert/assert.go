package assert

import "testing"

func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got %q; expected %q", actual, expected)
	}

}
