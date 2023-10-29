package main

import (
	"fmt"
	"os"
)

func stopIfErrf(format string, err error) {
	if err != nil {
		fmt.Printf(format, err)
		os.Exit(1)
	}
}
