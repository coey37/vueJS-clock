package main

import (
	"log"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/handler"
	"github.com/VolticFroogo/Froogo/middleware/myJWT"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Printf("Error initializing 