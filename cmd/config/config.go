/*
Package config создание конфиг файла для сервера
*/
package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
)

// ServerConfig структура конфига для сервера.
type ServerConfig struct {
	Addr         string `json:"address" env:"ADDRESS"`
	DatabaseDSN  string `json:"database_dsn" env:"DATABASE_DSN"`
	Limit        int    `json:"limit" env:"LIMIT"`
	AddrAge      string `json:"addr_age" env:"ADDR_AGE"`
	AddrGender   string `json:"addr_gender" env:"ADDR_GENDER"`
	AddrNational string `json:"addr_national" env:"ADDR_NATIONAL"`
	Kafka        Kafka
}

type Kafka struct {
	Addr          string `json:"kafka_addr" env:"KAFKA_ADDR"`
	TopicConsumer string `json:"consumer" env:"CONSUMER"`
	TopicProducer string `json:"producer" env:"PRODUCER"`
}

// ServerConfigInit инициализация конфига.
func ServerConfigInit() *ServerConfig {
	cfgSrv := &ServerConfig{}
	flag.StringVar(&cfgSrv.Addr, "a", "localhost:8080", "ADDRESS")
	flag.StringVar(&cfgSrv.DatabaseDSN, "d", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "DATABASE_DSN")
	flag.StringVar(&cfgSrv.AddrAge, "age", "https://api.agify.io/?name=", "API age")
	flag.StringVar(&cfgSrv.AddrGender, "gender", "https://api.genderize.io/?name=", "API gender")
	flag.StringVar(&cfgSrv.AddrNational, "national", "https://api.nationalize.io/?name=", "API national")
	flag.IntVar(&cfgSrv.Limit, "limit", 10, "limit for filter")
	flag.StringVar(&cfgSrv.Kafka.Addr, "kafka", "localhost:9092", "kafka addr")
	flag.StringVar(&cfgSrv.Kafka.TopicConsumer, "consumer", "FIO", "kafka consumer topic")
	flag.StringVar(&cfgSrv.Kafka.TopicProducer, "producer", "FIO_FAILED", "kafka producer topic")
	flag.Parse()
	if err := env.Parse(cfgSrv); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return cfgSrv
}
