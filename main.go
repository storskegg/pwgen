package main

import "fmt"

func main() {
	memguard.CatchInterrupt(func() {
		fmt.Println("Interrupt signal received. Exiting...")
	})
	// Make sure to destroy all LockedBuffers when returning.
	defer memguard.DestroyAll()

	getToIt()
}

func getToIt() {}