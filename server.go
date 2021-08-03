package main

import (
	//"log"

	//. "github.com/USACE/pallid_sturgeon_api/server/auth"

	"log"

	"github.com/USACE/pallid_sturgeon_api/server/config"
	"github.com/USACE/pallid_sturgeon_api/server/handlers"
	"github.com/USACE/pallid_sturgeon_api/server/stores"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var urlContext string = "/psapi"

func main() {
	appconfig := config.GetEnv()
	//auth := Auth{}

	//err := LoadVerificationKeys(appconfig.IPPK)
	// err := auth.LoadVerificationKey(appconfig.IPPK)
	// if err != nil {
	// 	log.Fatalf("Unable to load a verification key:%s.\nShutting down.", err)
	// }
	pallidSturgeonStore, err := stores.InitStores(appconfig)
	if err != nil {
		log.Printf("Unable to connect to the Main Pallid Sturgeon database: %s", err)
	}

	// authStore, err := stores.InitAuthStore(appconfig)
	// if err != nil {
	// 	log.Printf("Unable to connect to the Auth database: %s", err)
	// }

	//auth.Store = authStore

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	PallidSturgeonH := handlers.PallidSturgeonHandler{
		Store: pallidSturgeonStore,
	}

	// userH := handlers.UserHandler{
	// 	Config: appconfig,
	// }

	//r := mux.NewRouter()

	e.GET(urlContext+"/version", PallidSturgeonH.Version)
	e.GET(urlContext+"/projects", PallidSturgeonH.GetProjects)
	e.GET(urlContext+"/seasons", PallidSturgeonH.GetSeasons)
	e.GET(urlContext+"/segments", PallidSturgeonH.GetSegments)
	e.GET(urlContext+"/bends", PallidSturgeonH.GetBends)
	e.GET(urlContext+"/fishDataSummary", PallidSturgeonH.GetFishDataSummary)
	e.GET(urlContext+"/uploadSessionId", PallidSturgeonH.GetUploadSessionId)
	e.POST(urlContext+"/upload", PallidSturgeonH.Upload)
	e.POST(urlContext+"/siteUpload", PallidSturgeonH.SiteUpload)
	e.POST(urlContext+"/fishUpload", PallidSturgeonH.FishUpload)
	e.POST(urlContext+"/searchUpload", PallidSturgeonH.SearchUpload)
	e.POST(urlContext+"/supplementalUpload", PallidSturgeonH.SupplementalUpload)
	e.POST(urlContext+"/procedureUpload", PallidSturgeonH.ProcedureUpload)
	e.POST(urlContext+"/moriverUpload", PallidSturgeonH.MoriverUpload)
	e.POST(urlContext+"/telemetryUpload", PallidSturgeonH.TelemetryUpload)
	e.POST(urlContext+"/storeProcedure/:uploadSessionId", PallidSturgeonH.CallStoreProcedures)
	// e.Logger.Fatal(e.Start(":8080"))
	e.Logger.Debug(e.Start(":8080"))
}
