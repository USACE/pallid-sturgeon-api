package main

import (
	//"log"

	. "github.com/USACE/pallid_sturgeon_api/server/auth"

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
	auth := Auth{}

	err := auth.LoadVerificationKey(appconfig.IPPK)
	if err != nil {
		log.Fatalf("Unable to load a verification key:%s.\nShutting down.", err)
	}
	pallidSturgeonStore, err := stores.InitStores(appconfig)
	if err != nil {
		log.Printf("Unable to connect to the Main Pallid Sturgeon database: %s", err)
	}

	authStore, err := stores.InitAuthStore(appconfig)
	if err != nil {
		log.Printf("Unable to connect to the Auth database: %s", err)
	}

	auth.Store = authStore

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	PallidSturgeonH := handlers.PallidSturgeonHandler{
		Store: pallidSturgeonStore,
	}

	userH := handlers.UserHandler{
		Store: authStore,
	}

	e.GET(urlContext+"/version", PallidSturgeonH.Version)

	e.GET(urlContext+"/projects", auth.Authorize(PallidSturgeonH.GetProjects, PUBLIC))
	e.GET(urlContext+"/roles", auth.Authorize(PallidSturgeonH.GetRoles, PUBLIC))
	e.GET(urlContext+"/fieldOffices", auth.Authorize(PallidSturgeonH.GetFieldOffices, PUBLIC))
	e.GET(urlContext+"/seasons", auth.Authorize(PallidSturgeonH.GetSeasons, PUBLIC))
	e.GET(urlContext+"/segments", auth.Authorize(PallidSturgeonH.GetSegments, PUBLIC))
	e.GET(urlContext+"/fieldOffices", auth.Authorize(PallidSturgeonH.GetFieldOffices, PUBLIC))
	e.GET(urlContext+"/sampleMethods", auth.Authorize(PallidSturgeonH.GetSampleMethods, PUBLIC))
	e.GET(urlContext+"/sampleUnitTypes", auth.Authorize(PallidSturgeonH.GetSampleUnitTypes, PUBLIC))
	e.GET(urlContext+"/bends", auth.Authorize(PallidSturgeonH.GetBends, PUBLIC))
	e.GET(urlContext+"/siteDataEntry", auth.Authorize(PallidSturgeonH.GetSiteDataEntries, PUBLIC))
	e.GET(urlContext+"/siteDataEntryById", auth.Authorize(PallidSturgeonH.GetSiteDataEntryById, PUBLIC))
	e.POST(urlContext+"/siteDataEntry", auth.Authorize(PallidSturgeonH.SaveSiteDataEntry, PUBLIC))
	e.PUT(urlContext+"/siteDataEntry", auth.Authorize(PallidSturgeonH.UpdateSiteDataEntry, PUBLIC))
	e.GET(urlContext+"/fishDataEntry", auth.Authorize(PallidSturgeonH.GetFishDataEntries, PUBLIC))
	e.POST(urlContext+"/fishDataEntry", auth.Authorize(PallidSturgeonH.SaveFishDataEntry, PUBLIC))
	e.PUT(urlContext+"/fishDataEntry", auth.Authorize(PallidSturgeonH.UpdateFishDataEntry, PUBLIC))
	e.GET(urlContext+"/moriverDataEntry", auth.Authorize(PallidSturgeonH.GetMoriverDataEntries, PUBLIC))
	e.POST(urlContext+"/moriverDataEntry", auth.Authorize(PallidSturgeonH.SaveMoriverDataEntry, PUBLIC))
	e.PUT(urlContext+"/moriverDataEntry", auth.Authorize(PallidSturgeonH.UpdateMoriverDataEntry, PUBLIC))
	e.GET(urlContext+"/supplementalDataEntry", auth.Authorize(PallidSturgeonH.GetSupplementalDataEntries, PUBLIC))
	e.POST(urlContext+"/supplementalDataEntry", auth.Authorize(PallidSturgeonH.SaveSupplementalDataEntry, PUBLIC))
	e.PUT(urlContext+"/supplementalDataEntry", auth.Authorize(PallidSturgeonH.UpdateSupplementalDataEntry, PUBLIC))
	e.GET(urlContext+"/fishFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullFishDataSummary, PUBLIC))
	e.GET(urlContext+"/fishDataSummary", auth.Authorize(PallidSturgeonH.GetFishDataSummary, PUBLIC))
	e.GET(urlContext+"/suppFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullSuppDataSummary, PUBLIC))
	e.GET(urlContext+"/suppDataSummary", auth.Authorize(PallidSturgeonH.GetSuppDataSummary, PUBLIC))
	e.GET(urlContext+"/missouriFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullMissouriDataSummary, PUBLIC))
	e.GET(urlContext+"/missouriDataSummary", auth.Authorize(PallidSturgeonH.GetMissouriDataSummary, PUBLIC))
	e.GET(urlContext+"/geneticFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullGeneticDataSummary, PUBLIC))
	e.GET(urlContext+"/geneticDataSummary", auth.Authorize(PallidSturgeonH.GetGeneticDataSummary, PUBLIC))
	e.GET(urlContext+"/searchFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullSearchDataSummary, PUBLIC))
	e.GET(urlContext+"/searchDataSummary", auth.Authorize(PallidSturgeonH.GetSearchDataSummary, PUBLIC))
	e.GET(urlContext+"/telemetryFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullTelemetryDataSummary, PUBLIC))
	e.GET(urlContext+"/telemetryDataSummary", auth.Authorize(PallidSturgeonH.GetTelemetryDataSummary, PUBLIC))
	e.GET(urlContext+"/procedureFullDataSummary", auth.Authorize(PallidSturgeonH.GetFullProcedureDataSummary, PUBLIC))
	e.GET(urlContext+"/procedureDataSummary", auth.Authorize(PallidSturgeonH.GetProcedureDataSummary, PUBLIC))
	e.GET(urlContext+"/uploadSessionId", auth.Authorize(PallidSturgeonH.GetUploadSessionId, PUBLIC))
	e.POST(urlContext+"/upload", auth.Authorize(PallidSturgeonH.Upload, PUBLIC))
	// e.POST(urlContext+"/siteUpload", auth.Authorize(PallidSturgeonH.SiteUpload, PUBLIC))
	// e.POST(urlContext+"/fishUpload", auth.Authorize(PallidSturgeonH.FishUpload, PUBLIC))
	// e.POST(urlContext+"/searchUpload", auth.Authorize(PallidSturgeonH.SearchUpload, PUBLIC))
	// e.POST(urlContext+"/supplementalUpload", auth.Authorize(PallidSturgeonH.SupplementalUpload, PUBLIC))
	// e.POST(urlContext+"/procedureUpload", auth.Authorize(PallidSturgeonH.ProcedureUpload, PUBLIC))
	// e.POST(urlContext+"/moriverUpload", auth.Authorize(PallidSturgeonH.MoriverUpload, PUBLIC))
	// e.POST(urlContext+"/telemetryUpload", auth.Authorize(PallidSturgeonH.TelemetryUpload, PUBLIC))
	e.POST(urlContext+"/storeProcedure/:uploadSessionId", auth.Authorize(PallidSturgeonH.CallStoreProcedures, PUBLIC))
	e.GET(urlContext+"/errorCount", auth.Authorize(PallidSturgeonH.GetErrorCount, PUBLIC))
	e.GET(urlContext+"/officeErrorLog", auth.Authorize(PallidSturgeonH.GetOfficeErrorLogs, PUBLIC))
	e.GET(urlContext+"/usgNoVialNumbers", auth.Authorize(PallidSturgeonH.GetUsgNoVialNumbers, PUBLIC))
	e.GET(urlContext+"/unapprovedDataSheets", auth.Authorize(PallidSturgeonH.GetUnapprovedDataSheets, PUBLIC))
	e.GET(urlContext+"/uncheckedDataSheets", auth.Authorize(PallidSturgeonH.GetUncheckedDataSheets, PUBLIC))
	e.POST(urlContext+"/uploadDownloadZip", auth.Authorize(PallidSturgeonH.UploadDownloadZip, PUBLIC))
	e.GET(urlContext+"/downloadInfo", auth.Authorize(PallidSturgeonH.GetDownloadInfo, PUBLIC))
	e.GET(urlContext+"/downloadZip", auth.Authorize(PallidSturgeonH.GetDownloadZip, PUBLIC))

	e.GET(urlContext+"/userRoleOffice/:email", auth.Authorize(userH.GetUserRoleOffice, PUBLIC))
	e.GET(urlContext+"/userAccessRequests", auth.Authorize(userH.GetUserAccessRequests, ADMIN))
	e.GET(urlContext+"/users", auth.Authorize(userH.GetUsers, ADMIN))
	e.POST(urlContext+"/userRoleOffice", auth.Authorize(userH.AddUserRoleOffice, ADMIN))
	e.PUT(urlContext+"/userRoleOffice", auth.Authorize(userH.UpdateUserRoleOffice, ADMIN))

	// e.Logger.Fatal(e.Start(":8080"))
	// force update
	e.Logger.Debug(e.Start(":8080"))
}
