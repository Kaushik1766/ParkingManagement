package main

import (
	"context"

	buildingrepository "github.com/Kaushik1766/ParkingManagement/internal/repository/building_repository"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type User struct {
	PK    string
	SK    string
	Name  string
	Email string
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("aws config not found")
	}

	client := dynamodb.NewFromConfig(cfg)

	buildingRepo := buildingrepository.NewNOSQLBuidlingRepository(client)

	// fmt.Println(buildingRepo.DeleteBuildingByID("fadfasfa"))

	// fmt.Print(buildingRepo.AddBuilding("advant"))

	// buildings, err := buildingRepo.GetAllBuildingSummary()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(buildings)

	// building, err := buildingRepo.GetBuildingByID(uuid.MustParse("88898cb2-6534-4602-82a0-6c93a6efc095"))
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(building)

	err = buildingRepo.DeleteBuildingByID("88898cb2-6534-4602-82a0-6c93a6efc095")
	if err != nil {
		panic(err)
	}

	// gormDb, _ := db.InitDB()

	// 	gormDb,
	// 	models.Building{},
	// 	models.Floor{},
	// 	models.Slot{},
	// 	models.Office{},
	// 	models.User{},
	// 	models.Vehicle{},
	// 	models.ParkingHistory{},
	// )
	// if err != nil {
	// 	panic("error migrating models" + err.Error())
	// }

	// app := app.NewApp(gormDb)
	// logFile, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.SetOutput(logFile)
	// app.Run()
}
