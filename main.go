// SPDX-FileCopyrightText: 2020 Kent Gibson <warthog618@gmail.com>
//
// SPDX-License-Identifier: MIT

// A simple example that toggles an output pin.
package main

import (
	"fmt"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"os"
	"os/signal"
	"syscall"
)

// This example drives GPIO 4, which is pin J8-7 on a Raspberry Pi.
// The pin is toggled high and low at 1Hz with a 50% duty cycle.
// Do not run this on a device which has this pin externally driven.
func main() {
	c, err := gpiod.NewChip("gpiochip0")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	values := map[int]string{0: "inactive", 1: "active"}
	ledPin := rpi.GPIO17
	buttonPin := rpi.GPIO18
	v := 0
	l, err := c.RequestLine(ledPin, gpiod.AsOutput(v))
	if err != nil {
		panic(err)
	}
	defer func() {
		l.Reconfigure(gpiod.AsInput)
		l.Close()
	}()

	fmt.Printf("Set pin %d %s\n", ledPin, values[v])

	b, err := c.RequestLine(buttonPin, gpiod.AsInput)
	if err != nil {
		panic(err)
	}

	defer func() {
		b.Reconfigure(gpiod.AsInput)
		b.Close()
	}()

	// capture exit signals to ensure pin is reverted to input on exit.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)

	for {
		select {
		case <-quit:
			return
		default:
			r, err := b.Value()

			if err != nil {
				panic(err)
			}

			if r == 1 {
				v = 0
			} else {
				v = 1
			}
			l.SetValue(v)
			fmt.Printf("Set pin %d %s\n", ledPin, values[v])
		}
	}
}
