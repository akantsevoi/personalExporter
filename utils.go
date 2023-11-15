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

func ptr[T any](v T) *T {
	return &v
}

type pair[T1 any, T2 any] struct {
	First  T1
	Second T2
}
