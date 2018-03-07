package main

import (
	"fmt"
	"github.com/awnumar/memguard"
	"log"
	"unsafe"
)

const pwLen = 64

func main() {
	memguard.DisableUnixCoreDumps()

	memguard.CatchInterrupt(func() {
		fmt.Println("Interrupt signal received. Exiting...")
		memguard.SafeExit(1)
	})
	// Make sure to destroy all LockedBuffers when returning.
	defer memguard.DestroyAll()

	generatePassword()
}

// generatePassword is the meat of the application
func generatePassword() {
	log.Println("Generating password of length", pwLen)

	b, err := memguard.NewImmutableRandom(pwLen)
	if err != nil {
		fmt.Println(err)
		memguard.SafeExit(2)
	}

	defer b.Destroy()

	bPtr := (*[pwLen]byte)(unsafe.Pointer(&b.Buffer()[0]))

	str := Encode(bPtr)

	fmt.Println("Password: ", str, len(str))
}
