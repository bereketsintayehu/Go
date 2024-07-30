package main

import (
    "library/services"
    "library/controllers"
)

func main() {
    // Initialize the library service and controller
    libraryService := services.NewLibraryService()
    libraryController := controllers.NewLibraryController(libraryService)

    libraryController.Menu()
}