// Package main is the entry point for the Ollanta CLI scanner.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/user/ollanta/ollantascanner/scan"
)

func main() {
	opts, err := scan.ParseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	ctx := context.Background()
	r, err := scan.Run(ctx, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scan error:", err)
		os.Exit(1)
	}

	scan.PrintSummary(r)

	switch opts.Format {
	case "json":
		if path, err := r.SaveJSON(opts.ProjectDir); err != nil {
			fmt.Fprintln(os.Stderr, "json error:", err)
		} else {
			fmt.Println("Report saved to", path)
		}
	case "sarif":
		if path, err := r.SaveSARIF(opts.ProjectDir); err != nil {
			fmt.Fprintln(os.Stderr, "sarif error:", err)
		} else {
			fmt.Println("SARIF saved to", path)
		}
	case "all":
		if path, err := r.SaveJSON(opts.ProjectDir); err != nil {
			fmt.Fprintln(os.Stderr, "json error:", err)
		} else {
			fmt.Println("Report saved to", path)
		}
		if path, err := r.SaveSARIF(opts.ProjectDir); err != nil {
			fmt.Fprintln(os.Stderr, "sarif error:", err)
		} else {
			fmt.Println("SARIF saved to", path)
		}
	}
}
