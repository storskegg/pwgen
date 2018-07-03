package main

import (
	"fmt"
	"github.com/awnumar/memguard"
	"os"
	"strconv"
)

func main() {
	memguard.DisableUnixCoreDumps()

	memguard.CatchInterrupt(func() {
		fmt.Println("Interrupt signal received. Exiting...")
		memguard.SafeExit(1)
	})
	// Make sure to destroy all LockedBuffers when returning.
	defer memguard.DestroyAll()

	var pwLenArg string
	if len(os.Args) > 1 {
		pwLenArg = os.Args[1]
	} else {
		pwLenArg = "84"
	}

	var pwLen int64
	var err error

	pwLen, err = strconv.ParseInt(pwLenArg, 10, 64)
	if err != nil {
		pwLen = 84
	}

	fmt.Print(Encode(pwLen))
}
