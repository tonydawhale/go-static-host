package main

import (
	"log"

	"go-static-host/s3utils"
	"go-static-host/server"
	"go-static-host/mongoutils"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting Server...")
	loadEnv()
	s3utils.Init()
	mongoutils.Init()
	server.Init()
}

func loadEnv() {
	err := godotenv.Load()

    if err != nil {
        log.Fatal("Error loading .env file")
    }
	log.Println("Environment Variables Loaded")
}