package main

import (
	"task-manager/delivery/routers"
)
func main(){
	router := routers.NewRouter()
	router.InitRoutes()
	router.Router.Run(":8080")
}