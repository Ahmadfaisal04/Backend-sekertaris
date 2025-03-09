package main

import (
	"Sekertaris/config"
	"fmt"
)

func main() {
	config.ConnectDB()
	fmt.Print("Connected to database")
	
}