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

		err = sd.Store.SiteUploadSP(uploadSite.UploadedBy, uploadSite.UploadSessionId)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}
