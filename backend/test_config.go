package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/powerx-plugins/scrum/internal/config"
)

func main() {
	fmt.Println("Loading configuration...")
	
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 序列化配置为 JSON 进行查看
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	fmt.Println("Configuration loaded successfully:")
	fmt.Println(string(data))
}