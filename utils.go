package main

import (
	"fmt"
	"os"
)

func stopIfErrf(format string, err error, a ...any) {
	if err != nil {
		a = append(a, err)
		fmt.Printf(format, a)
		os.Exit(1)
	}
}
