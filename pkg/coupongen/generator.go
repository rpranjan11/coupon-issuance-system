// pkg/coupongen/generator.go
package coupongen

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

var (
	numbers        = "0123456789"
	koreanChars    = []rune("가나다라마바사아자차카타파하")
	generatorMutex sync.Mutex
	usedCodes      = make(map[string]struct{})
)

// GenerateCode generates a unique coupon code with Korean characters and numbers
func GenerateCode(length int) string {
	if length <= 0 {
		length = 10 // Default length
	}

	generatorMutex.Lock()
	defer generatorMutex.Unlock()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		var sb strings.Builder

		// Ensure at least one Korean character
		sb.WriteRune(koreanChars[r.Intn(len(koreanChars))])

		// Fill the rest with a mix of numbers and Korean characters
		for i := 1; i < length; i++ {
			if r.Intn(2) == 0 {
				sb.WriteByte(numbers[r.Intn(len(numbers))])
			} else {
				sb.WriteRune(koreanChars[r.Intn(len(koreanChars))])
			}
		}

		code := sb.String()

		// Check if code already exists
		if _, exists := usedCodes[code]; !exists {
			usedCodes[code] = struct{}{}
			return code
		}
	}
}
