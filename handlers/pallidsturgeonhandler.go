package handlers

import (
	"net/http"
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

		err = sd.Store.UploadSiteDatasheetCheck(uploadSite.UploadedBy, uploadSite.UploadSessionId)
		if err != nil {
			return err
		}

		err = sd.Store.UploadSiteDatasheet(uploadSite.UploadedBy)
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

		err = sd.Store.UploadFishDatasheetCheck(uploadFish.UploadedBy, uploadFish.UploadSessionId)
		if err != nil {
			return err
		}

		err = sd.Store.UploadFishDatasheet(uploadFish.UploadedBy)
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

		// err = sd.Store.UploadSearchDatasheetCheck(uploadSearch.UploadedBy, uploadSearch.UploadSessionId)
		// if err != nil {
		// 	return err
		// }

		// err = sd.Store.UploadSearchDatasheet(uploadSearch.UploadedBy)
		// if err != nil {
		// 	return err
		// }
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

		// err = sd.Store.UploadSearchDatasheetCheck(uploadSupplemental.UploadedBy, uploadSupplemental.UploadSessionId)
		// if err != nil {
		// 	return err
		// }

		// err = sd.Store.UploadSearchDatasheet(uploadSupplemental.UploadedBy)
		// if err != nil {
		// 	return err
		// }
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

		// err = sd.Store.UploadSearchDatasheetCheck(uploadProcedure.UploadedBy, uploadProcedure.UploadSessionId)
		// if err != nil {
		// 	return err
		// }

		// err = sd.Store.UploadSearchDatasheet(uploadProcedure.UploadedBy)
		// if err != nil {
		// 	return err
		// }
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) MrUpload(c echo.Context) error {
	var err error
	uploadMrs := []models.UploadMr{}
	if err := c.Bind(&uploadMrs); err != nil {
		return err
	}
	for _, uploadMr := range uploadMrs {
		uploadMr.LastUpdated = time.Now()
		uploadMr.UploadedBy = "DeeLiang"
		err = sd.Store.SaveMrUpload(uploadMr)
		if err != nil {
			return err
		}

		// err = sd.Store.UploadSearchDatasheetCheck(uploadMr.UploadedBy, uploadMr.UploadSessionId)
		// if err != nil {
		// 	return err
		// }

		// err = sd.Store.UploadSearchDatasheet(uploadMr.UploadedBy)
		// if err != nil {
		// 	return err
		// }
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) TelemetryFishUpload(c echo.Context) error {
	var err error
	uploadTelemetryFishes := []models.UploadTelemetryFish{}
	if err := c.Bind(&uploadTelemetryFishes); err != nil {
		return err
	}
	for _, uploadTelemetryFish := range uploadTelemetryFishes {
		uploadTelemetryFish.LastUpdated = time.Now()
		uploadTelemetryFish.UploadedBy = "DeeLiang"
		err = sd.Store.SaveTelemetryFishUpload(uploadTelemetryFish)
		if err != nil {
			return err
		}

		// err = sd.Store.UploadSearchDatasheetCheck(uploadTelemetryFish.UploadedBy, uploadTelemetryFish.UploadSessionId)
		// if err != nil {
		// 	return err
		// }

		// err = sd.Store.UploadSearchDatasheet(uploadTelemetryFish.UploadedBy)
		// if err != nil {
		// 	return err
		// }
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}
