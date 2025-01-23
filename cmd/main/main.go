package main

import "effective-mobile-task/internal/config"

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	_ = cfg
}
