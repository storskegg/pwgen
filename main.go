package main

import (
	"fmt"
	"github.com/awnumar/memguard"
	"log"
)

const pwLen = 84

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

	str := Encode()

	fmt.Println("Password: ", str, len(str))
}
