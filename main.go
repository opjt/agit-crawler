package main

import (
	"agit-crawler/app/core/bootstrap"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	bootstrap.RunServer(bootstrap.CommonModules)
}
