package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/chill-institute/chill-contracts/v2/internal/workflowpolicy"
)

func main() {
	against := flag.String("against", "", "baseline Git ref, or release for the preceding stable release")
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := workflowpolicy.CheckCompatibility(ctx, ".", *against, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
