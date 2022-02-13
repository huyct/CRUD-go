package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/huyct/CRUD-go/middlewares"
	"github.com/huyct/CRUD-go/routes"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	router := httprouter.New()

}
