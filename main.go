package main

import (
	"fmt"
	"github.com/awnumar/memguard"
	"gopkg.in/urfave/cli.v2"
	"os"
)

func main() {
	memguard.DisableUnixCoreDumps()

	memguard.CatchInterrupt(func() {
		fmt.Println("Interrupt signal received. Exiting...")
		memguard.SafeExit(1)
	})
	// Make sure to destroy all LockedBuffers when returning.
	defer memguard.DestroyAll()

	app := &cli.App{
		Name:  "greet",
		Usage: "say a greeting",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "hex",
				Usage: "Hexadecimal Output",
				Value: false,
			},
			&cli.IntFlag{
				Name:  "n",
				Usage: "Output Length",
				Value: 84,
			},
		},
		Action: func(c *cli.Context) error {
			pwLen := c.Int("n")

			fmt.Println(Encode(pwLen, c.Bool("hex")))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		memguard.DestroyAll()
		panic(err)
	}
}
