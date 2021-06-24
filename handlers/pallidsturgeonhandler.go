package handlers

import (
	"net/http"
	"strconv"
	"time"

	"di2e.net/cwbi/pallid_sturgeon_api/server/models"
	"di2e.net/cwbi/pallid_sturgeon_api/server/stores"
	"github.com/labstack/echo/v4"
)

type PallidSturgeonHandler struct {
	Store *stores.PallidSturgeonStore
}

func (ps *PallidSturgeonHandler) Version(c echo.Context) error {
	return c.String(http.StatusOK, "Pallid Sturgeon API v0.01")
}

func (sd *PallidSturgeonHandler) GetSeasons(c echo.Context) error {
	seasons, err := sd.Store.GetSeasons()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, seasons)
}

func (sd *PallidSturgeonHandler) GetUploadSessionId(c echo.Context) error {
	seasons, err := sd.Store.GetUploadSessionId()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, seasons)
}

func (sd *PallidSturgeonHandler) SiteUpload(c echo.Context) error {
	var err error
	uploadSites := []models.UploadSite{}
	if err := c.Bind(&uploadSites); err != nil {
		return err
	}
	for _, uploadSite := range uploadSites {
		uploadSite.LastUpdated = time.Now()
		uploadSite.UploadedBy = "DeeLiang"
		err = sd.Store.SaveSiteUpload(uploadSite)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) FishUpload(c echo.Context) error {
	var err error
	uploadFishs := []models.UploadFish{}
	if err := c.Bind(&uploadFishs); err != nil {
		return err
	}
	for _, uploadFish := range uploadFishs {
		uploadFish.LastUpdated = time.Now()
		uploadFish.UploadedBy = "DeeLiang"
		err = sd.Store.SaveFishUpload(uploadFish)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) SearchUpload(c echo.Context) error {
	var err error
	uploadSearches := []models.UploadSearch{}
	if err := c.Bind(&uploadSearches); err != nil {
		return err
	}
	for _, uploadSearch := range uploadSearches {
		uploadSearch.LastUpdated = time.Now()
		uploadSearch.UploadedBy = "DeeLiang"
		err = sd.Store.SaveSearchUpload(uploadSearch)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) SupplementalUpload(c echo.Context) error {
	var err error
	uploadSupplementals := []models.UploadSupplemental{}
	if err := c.Bind(&uploadSupplementals); err != nil {
		return err
	}
	for _, uploadSupplemental := range uploadSupplementals {
		uploadSupplemental.LastUpdated = time.Now()
		uploadSupplemental.UploadedBy = "DeeLiang"
		err = sd.Store.SaveSupplementalUpload(uploadSupplemental)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) ProcedureUpload(c echo.Context) error {
	var err error
	uploadProcedures := []models.UploadProcedure{}
	if err := c.Bind(&uploadProcedures); err != nil {
		return err
	}
	for _, uploadProcedure := range uploadProcedures {
		uploadProcedure.LastUpdated = time.Now()
		uploadProcedure.UploadedBy = "DeeLiang"
		err = sd.Store.SaveProcedureUpload(uploadProcedure)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) MoriverUpload(c echo.Context) error {
	var err error
	UploadMorivers := []models.UploadMoriver{}
	if err := c.Bind(&UploadMorivers); err != nil {
		return err
	}
	for _, uploadMoriver := range UploadMorivers {
		uploadMoriver.LastUpdated = time.Now()
		uploadMoriver.UploadedBy = "DeeLiang"
		err = sd.Store.SaveMoriverUpload(uploadMoriver)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) TelemetryUpload(c echo.Context) error {
	var err error
	uploadTelemetrys := []models.UploadTelemetry{}
	if err := c.Bind(&uploadTelemetrys); err != nil {
		return err
	}
	for _, uploadTelemetry := range uploadTelemetrys {
		uploadTelemetry.LastUpdated = time.Now()
		uploadTelemetry.UploadedBy = "DeeLiang"
		err = sd.Store.SaveTelemetryUpload(uploadTelemetry)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) CallStoreProcedures(c echo.Context) error {
	var err error
	uploadSessionId := c.Param("uploadSessionId")
	id, err := strconv.Atoi(uploadSessionId)
	if err != nil {
		return err
	}

	procedureOut, err := sd.Store.CallStoreProcedures("DeeLiang", id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, procedureOut)
}
