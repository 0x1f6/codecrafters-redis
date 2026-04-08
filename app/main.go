package main

import (
	"fmt"
	"os"

	"github.com/0x1f6/codecrafters-redis/internal/redis"
	"github.com/0x1f6/codecrafters-redis/internal/tcp"
)

func main() {
	fmt.Println("=== Program starting... ===")

	r := redis.New()
	t := tcp.NewServer("0.0.0.0:6379", r)
	if err := t.ServeForever(); err != nil {
		fmt.Println("Server error:", err)
		os.Exit(1)
	}
}
