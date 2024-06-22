package main

import (
	"fmt"

	"github.com/iamsorryprincess/project-layout/cmd/service-a/config"
	"github.com/iamsorryprincess/project-layout/internal/pkg/cfg"
)

func main() {
	fmt.Println("test")

	var configuration config.Config
	if err := cfg.Parse("cmd/service-a/config.json", &configuration); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(configuration)
}
