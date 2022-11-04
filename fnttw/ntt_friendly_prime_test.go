package fnttw

import (
	"testing"
)

func TestNttFriendlyPrimes(t *testing.T) {
	NttFriendlyPrimes(1<<16, 50, 63)
}
