package main

import (
	"fmt"
	"log"
	"os"

	"shortURL/config"
	"shortURL/dao"
	"shortURL/router"
	"shortURL/service"
)

func main() {
	cfgPath := "config.yaml"
	if v := os.Getenv("CONFIG_PATH"); v != "" {
		cfgPath = v
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Printf("config file not found, using defaults: %v", err)
		log.Println("set CONFIG_PATH env to override; refer to config.yaml.example")
		cfg = config.Default
	}

	dao.InitMySQL(cfg.MySQL.DSN)
	dao.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	service.InitBloom(cfg.Bloom.Capacity, cfg.Bloom.FalsePositive)
	if err := service.LoadBloomFromDB(); err != nil {
		log.Printf("warn: failed to load bloom from db: %v", err)
	}

	r := router.Setup(cfg.Server.Domain)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("shortURL server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
