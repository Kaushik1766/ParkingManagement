package main

import (
	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/app"
)

func main() {
	gormDb, _ := db.InitDB()

	app := app.NewApp(gormDb)
	// logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.SetOutput(logFile)
	app.Run()
}
