package webapp

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generatePromoCode() string {
	n := rand.Intn(100000) // 0..99999
	return fmt.Sprintf("%05d", n)
}
