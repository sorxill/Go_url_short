package main

import (
	"fmt"
	"go_url_short/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
