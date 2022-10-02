package main

import (
	"context"
	"nglapi/database"
	"nglapi/global"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	global.ContextConsume = context.Background()
	database.InitDatabase()
	StartWeb()
}
