package hash

import (
	"crypto/sha256"
	"fmt"
)

const hashFormat = "%s:%s" // salt:password
const hashRepeats = 100

func CreateHash(salt, password string) string {
	sum := []byte(fmt.Sprintf(hashFormat, salt, password))
	var crutch [32]byte
	for i := 0; i < hashRepeats; i++ {
		crutch = sha256.Sum256(sum)
		sum = crutch[:]
	}

	return string(sum)
}
