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
	// auth := Auth{}

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

	// auth.Store = authStore

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	PallidSturgeonH := handlers.PallidSturgeonHandler{
		Store: pallidSturgeonStore,
	}

	// userH := handlers.UserHandler{
	// 	Store: authStore,
	// }

	e.GET(urlContext+"/version", PallidSturgeonH.Version)

	e.GET(urlContext+"/projects", PallidSturgeonH.GetProjects)
	e.GET(urlContext+"/roles", PallidSturgeonH.GetRoles)
	e.GET(urlContext+"/fieldOffices", PallidSturgeonH.GetFieldOffices)
	e.GET(urlContext+"/seasons", PallidSturgeonH.GetSeasons)
	e.GET(urlContext+"/segments", PallidSturgeonH.GetSegments)
	e.GET(urlContext+"/fieldOffices", PallidSturgeonH.GetFieldOffices)
	e.GET(urlContext+"/sampleMethods", PallidSturgeonH.GetSampleMethods)
	e.GET(urlContext+"/sampleUnitTypes", PallidSturgeonH.GetSampleUnitTypes)
	e.GET(urlContext+"/bends", PallidSturgeonH.GetBends)
	e.GET(urlContext+"/siteDataEntry", PallidSturgeonH.GetSiteDataEntries)
	e.POST(urlContext+"/siteDataEntry", PallidSturgeonH.SaveSiteDataEntry)
	e.PUT(urlContext+"/siteDataEntry", PallidSturgeonH.UpdateSiteDataEntry)
	e.GET(urlContext+"/fishDataEntry", PallidSturgeonH.GetFishDataEntries)
	e.POST(urlContext+"/fishDataEntry", PallidSturgeonH.SaveFishDataEntry)
	e.PUT(urlContext+"/fishDataEntry", PallidSturgeonH.UpdateFishDataEntry)
	e.GET(urlContext+"/moriverDataEntry", PallidSturgeonH.GetMoriverDataEntries)
	e.POST(urlContext+"/moriverDataEntry", PallidSturgeonH.SaveMoriverDataEntry)
	e.PUT(urlContext+"/moriverDataEntry", PallidSturgeonH.UpdateMoriverDataEntry)
	e.GET(urlContext+"/supplementalDataEntry", PallidSturgeonH.GetSupplementalDataEntries)
	e.POST(urlContext+"/supplementalDataEntry", PallidSturgeonH.SaveSupplementalDataEntry)
	e.PUT(urlContext+"/supplementalDataEntry", PallidSturgeonH.UpdateSupplementalDataEntry)
	e.GET(urlContext+"/fishFullDataSummary", PallidSturgeonH.GetFullFishDataSummary)
	e.GET(urlContext+"/fishDataSummary", PallidSturgeonH.GetFishDataSummary)
	e.GET(urlContext+"/suppFullDataSummary", PallidSturgeonH.GetFullSuppDataSummary)
	e.GET(urlContext+"/suppDataSummary", PallidSturgeonH.GetSuppDataSummary)
	e.GET(urlContext+"/missouriFullDataSummary", PallidSturgeonH.GetFullMissouriDataSummary)
	e.GET(urlContext+"/missouriDataSummary", PallidSturgeonH.GetMissouriDataSummary)
	e.GET(urlContext+"/geneticFullDataSummary", PallidSturgeonH.GetFullGeneticDataSummary)
	e.GET(urlContext+"/geneticDataSummary", PallidSturgeonH.GetGeneticDataSummary)
	e.GET(urlContext+"/searchFullDataSummary", PallidSturgeonH.GetFullSearchDataSummary)
	e.GET(urlContext+"/searchDataSummary", PallidSturgeonH.GetSearchDataSummary)
	e.GET(urlContext+"/telemetryFullDataSummary", PallidSturgeonH.GetFullTelemetryDataSummary)
	e.GET(urlContext+"/telemetryDataSummary", PallidSturgeonH.GetTelemetryDataSummary)
	e.GET(urlContext+"/procedureFullDataSummary", PallidSturgeonH.GetFullProcedureDataSummary)
	e.GET(urlContext+"/procedureDataSummary", PallidSturgeonH.GetProcedureDataSummary)
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
	e.GET(urlContext+"/errorCount", PallidSturgeonH.GetErrorCount)
	e.GET(urlContext+"/usgNoVialNumbers", PallidSturgeonH.GetUsgNoVialNumbers)
	e.GET(urlContext+"/unapprovedDataSheets", PallidSturgeonH.GetUnapprovedDataSheets)
	e.GET(urlContext+"/uncheckedDataSheets", PallidSturgeonH.GetUncheckedDataSheets)
	e.POST(urlContext+"/uploadDownloadZip", PallidSturgeonH.UploadDownloadZip)
	e.GET(urlContext+"/downloadInfo", PallidSturgeonH.GetDownloadInfo)
	e.GET(urlContext+"/downloadZip", PallidSturgeonH.GetDownloadZip)

	// e.GET(urlContext+"/userRoleOffice/:email", auth.Authorize(userH.GetUserRoleOffice, PUBLIC))
	// e.GET(urlContext+"/userAccessRequests", auth.Authorize(userH.GetUserAccessRequests, PUBLIC))
	// e.POST(urlContext+"/userRoleOffice", auth.Authorize(userH.AddUserRoleOffice, PUBLIC))

	// e.GET(urlContext+"/projects", auth.Authorize(PallidSturgeonH.GetProjects, PUBLIC))
	// e.GET(urlContext+"/roles", auth.Authorize(PallidSturgeonH.GetRoles, PUBLIC))
	// e.GET(urlContext+"/fieldOffices", auth.Authorize(PallidSturgeonH.GetFieldOffices, PUBLIC))
	// e.GET(urlContext+"/seasons", auth.Authorize(PallidSturgeonH.GetSeasons, PUBLIC))
	// e.GET(urlContext+"/segments", auth.Authorize(PallidSturgeonH.GetSegments, PUBLIC))
	// e.GET(urlContext+"/fieldOffices", auth.Authorize(PallidSturgeonH.GetFieldOffices, PUBLIC))
	// e.GET(urlContext+"/sampleMethods", auth.Authorize(PallidSturgeonH.GetSampleMethods, PUBLIC))
	// e.GET(urlContext+"/sampleUnitTypes", auth.Authorize(PallidSturgeonH.GetSampleUnitTypes, PUBLIC))
	// e.GET(urlContext+"/bends", auth.Authorize(PallidSturgeonH.GetBends, PUBLIC))
	// e.GET(urlContext+"/siteDataEntry", auth.Authorize(PallidSturgeonH.GetSiteDataEntries, PUBLIC))
	// e.POST(urlContext+"/siteDataEntry", auth.Authorize(PallidSturgeonH.SaveSiteDataEntry, PUBLIC))
	// e.PUT(urlContext+"/siteDataEntry", auth.Authorize(PallidSturgeonH.UpdateSiteDataEntry, PUBLIC))
	// e.GET(urlContext+"/fishDataEntry", auth.Authorize(PallidSturgeonH.GetFishDataEntries, PUBLIC))
	// e.POST(urlContext+"/fishDataEntry", auth.Authorize(PallidSturgeonH.SaveFishDataEntry, PUBLIC))
	// e.PUT(urlContext+"/fishDataEntry", auth.Authorize(PallidSturgeonH.UpdateFishDataEntry, PUBLIC))
	// e.GET(urlContext+"/moriverDataEntry", auth.Authorize(PallidSturgeonH.GetMoriverDataEntries, PUBLIC))
	// e.POST(urlContext+"/moriverDataEntry", auth.Authorize(PallidSturgeonH.SaveMoriverDataEntry, PUBLIC))
	// e.PUT(urlContext+"/moriverDataEntry", auth.Authorize(PallidSturgeonH.UpdateMoriverDataEntry, PUBLIC))
	// e.GET(urlContext+"/supplementalDataEntry", auth.Authorize(PallidSturgeonH.GetSupplementalDataEntries, PUBLIC))
	// e.POST(urlContext+"/supplementalDataEntry", auth.Authorize(PallidSturgeonH.SaveSupplementalDataEntry, PUBLIC))
	// e.PUT(urlContext+"/supplementalDataEntry", auth.Authorize(PallidSturgeonH.UpdateSupplementalDataEntry, PUBLIC))
	// e.GET(urlContext+"/fishFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullFishDataSummary, PUBLIC))
	// e.GET(urlContext+"/fishDataSummary", auth.Authorize(PallidSturgeonH.GetFishDataSummary, PUBLIC))
	// e.GET(urlContext+"/suppFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullSuppDataSummary, PUBLIC))
	// e.GET(urlContext+"/suppDataSummary", auth.Authorize(PallidSturgeonH.GetSuppDataSummary, PUBLIC))
	// e.GET(urlContext+"/missouriFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullMissouriDataSummary, PUBLIC))
	// e.GET(urlContext+"/missouriDataSummary", auth.Authorize(PallidSturgeonH.GetMissouriDataSummary, PUBLIC))
	// e.GET(urlContext+"/geneticFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullGeneticDataSummary, PUBLIC))
	// e.GET(urlContext+"/geneticDataSummary", auth.Authorize(PallidSturgeonH.GetGeneticDataSummary, PUBLIC))
	// e.GET(urlContext+"/searchFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullSearchDataSummary, PUBLIC))
	// e.GET(urlContext+"/searchDataSummary", auth.Authorize(PallidSturgeonH.GetSearchDataSummary, PUBLIC))
	// e.GET(urlContext+"/telemetryFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullTelemetryDataSummary, PUBLIC))
	// e.GET(urlContext+"/telemetryDataSummary", auth.Authorize(PallidSturgeonH.GetTelemetryDataSummary, PUBLIC))
	// e.GET(urlContext+"/procedureFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullProcedureDataSummary, PUBLIC))
	// e.GET(urlContext+"/procedureDataSummary", auth.Authorize(PallidSturgeonH.GetProcedureDataSummary, PUBLIC))
	// e.GET(urlContext+"/uploadSessionId", auth.Authorize(PallidSturgeonH.GetUploadSessionId, PUBLIC))
	// e.POST(urlContext+"/upload", auth.Authorize(PallidSturgeonH.Upload, PUBLIC))
	// e.POST(urlContext+"/siteUpload", auth.Authorize(PallidSturgeonH.SiteUpload, PUBLIC))
	// e.POST(urlContext+"/fishUpload", auth.Authorize(PallidSturgeonH.FishUpload, PUBLIC))
	// e.POST(urlContext+"/searchUpload", auth.Authorize(PallidSturgeonH.SearchUpload, PUBLIC))
	// e.POST(urlContext+"/supplementalUpload", auth.Authorize(PallidSturgeonH.SupplementalUpload, PUBLIC))
	// e.POST(urlContext+"/procedureUpload", auth.Authorize(PallidSturgeonH.ProcedureUpload, PUBLIC))
	// e.POST(urlContext+"/moriverUpload", auth.Authorize(PallidSturgeonH.MoriverUpload, PUBLIC))
	// e.POST(urlContext+"/telemetryUpload", auth.Authorize(PallidSturgeonH.TelemetryUpload, PUBLIC))
	// e.POST(urlContext+"/storeProcedure/:uploadSessionId", auth.Authorize(PallidSturgeonH.CallStoreProcedures, PUBLIC))
	// e.POST(urlContext+"/errorCount", auth.Authorize(PallidSturgeonH.GetErrorCount, PUBLIC))
	// e.POST(urlContext+"/usgNoVialNumbers", auth.Authorize(PallidSturgeonH.GetUsgNoVialNumbers, PUBLIC))
	// e.POST(urlContext+"/unapprovedDataSheets", auth.Authorize(PallidSturgeonH.GetUnapprovedDataSheets, PUBLIC))
	// e.POST(urlContext+"/uncheckedDataSheets", auth.Authorize(PallidSturgeonH.GetUncheckedDataSheets, PUBLIC))
	// e.GET(urlContext+"/downloadInfo", auth.Authorize(PallidSturgeonH.GetDownloadInfo))
	// e.GET(urlContext+"/downloadZip", auth.Authorize(PallidSturgeonH.GetDownloadZip))

	// e.Logger.Fatal(e.Start(":8080"))
	e.Logger.Debug(e.Start(":8080"))
}
