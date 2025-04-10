// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfexec

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func sleepMock(rawDuration string) {
	signal.Ignore(os.Interrupt)

	d, err := time.ParseDuration(rawDuration)
	if err != nil {
		log.Fatalf("invalid duration format: %s", err)
	}

	fmt.Printf("sleeping for %s\n", d)

	time.Sleep(d)
}
