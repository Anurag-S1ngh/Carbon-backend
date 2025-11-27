package otp

import (
	"fmt"
	"math/rand"
)

func GenerateOTP() string {
	return fmt.Sprintf("%06d", 100000+rand.Intn(900000))
}
