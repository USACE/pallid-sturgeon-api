package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/USACE/pallid_sturgeon_api/server/models"
	"github.com/USACE/pallid_sturgeon_api/server/stores"
	"github.com/labstack/echo/v4"
)

type PallidSturgeonHandler struct {
	Store *stores.PallidSturgeonStore
}

func (ps *PallidSturgeonHandler) Version(c echo.Context) error {
	return c.String(http.StatusOK, "Pallid Sturgeon API v0.01")
}

func (sd *PallidSturgeonHandler) GetProjects(c echo.Context) error {
	projects, err := sd.Store.GetProjects()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, projects)
}

func (sd *PallidSturgeonHandler) GetSeasons(c echo.Context) error {
	seasons, err := sd.Store.GetSeasons()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, seasons)
}

func (sd *PallidSturgeonHandler) GetSegments(c echo.Context) error {
	segments, err := sd.Store.GetSegments()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, segments)
}

func (sd *PallidSturgeonHandler) GetBends(c echo.Context) error {
	bends, err := sd.Store.GetBends()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, bends)
}

func (sd *PallidSturgeonHandler) GetFishDataEntries(c echo.Context) error {
	tableId, fieldId := c.QueryParam("tableId"), c.QueryParam("fieldId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return err
	}
	dataSummary, err := sd.Store.GetFishDataEntries(tableId, fieldId, queryParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveFishDataEntry(c echo.Context) error {
	fishData := models.UploadFish{}
	if err := c.Bind(&fishData); err != nil {
		return err
	}
	id, err := sd.Store.SaveFishDataEntry(fishData)
	if err != nil {
		return err
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateFishDataEntry(c echo.Context) error {

	fishData := models.UploadFish{}
	if err := c.Bind(&fishData); err != nil {
		return err
	}
	_, err := sd.Store.UpdateFishDataEntry(fishData)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetMoriverDataEntries(c echo.Context) error {
	tableId, fieldId := c.QueryParam("tableId"), c.QueryParam("fieldId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return err
	}
	dataSummary, err := sd.Store.GetMoriverDataEntries(tableId, fieldId, queryParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveMoriverDataEntry(c echo.Context) error {
	moriverData := models.UploadMoriver{}
	if err := c.Bind(&moriverData); err != nil {
		return err
	}
	id, err := sd.Store.SaveMoriverDataEntry(moriverData)
	if err != nil {
		return err
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateMoriverDataEntry(c echo.Context) error {

	moriverData := models.UploadMoriver{}
	if err := c.Bind(&moriverData); err != nil {
		return err
	}
	_, err := sd.Store.UpdateMoriverDataEntry(moriverData)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetSupplementalDataEntries(c echo.Context) error {
	tableId, fieldId, geneticsVial, pitTag := c.QueryParam("tableId"), c.QueryParam("fieldId"), c.QueryParam("geneticsVial"), c.QueryParam("pitTag")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return err
	}
	dataSummary, err := sd.Store.GetSupplementalDataEntries(tableId, fieldId, geneticsVial, pitTag, queryParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveSupplementalDataEntry(c echo.Context) error {
	supplementalData := models.UploadSupplemental{}
	if err := c.Bind(&supplementalData); err != nil {
		return err
	}
	id, err := sd.Store.SaveSupplementalDataEntry(supplementalData)
	if err != nil {
		return err
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateSupplementalDataEntry(c echo.Context) error {

	supplementalData := models.UploadSupplemental{}
	if err := c.Bind(&supplementalData); err != nil {
		return err
	}
	_, err := sd.Store.UpdateSupplementalDataEntry(supplementalData)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetFishDataSummary(c echo.Context) error {
	year, officeCode, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("officeCode"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return err
	}
	dataSummary, err := sd.Store.GetFishDataSummary(year, officeCode, project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetSuppDataSummary(c echo.Context) error {
	year, officeCode, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("officeCode"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return err
	}
	dataSummary, err := sd.Store.GetSuppDataSummary(year, officeCode, project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetMissouriDataSummary(c echo.Context) error {
	year, officeCode, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("officeCode"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return err
	}
	dataSummary, err := sd.Store.GetMissouriDataSummary(year, officeCode, project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetGeneticDataSummary(c echo.Context) error {
	year, officeCode, project, fromDate, toDate, broodstock, hatchwild, speciesId, archive := c.QueryParam("year"), c.QueryParam("officeCode"), c.QueryParam("project"), c.QueryParam("fromDate"), c.QueryParam("toDate"), c.QueryParam("broodstock"), c.QueryParam("hatchwild"), c.QueryParam("speciesId"), c.QueryParam("archive")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return err
	}
	dataSummary, err := sd.Store.GetGeneticDataSummary(year, officeCode, project, fromDate, toDate, broodstock, hatchwild, speciesId, archive, queryParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetSearchDataSummary(c echo.Context) error {
	queryParams, err := marshalQuery(c)
	if err != nil {
		return err
	}
	dataSummary, err := sd.Store.GetSearchDataSummary(queryParams)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetUploadSessionId(c echo.Context) error {
	sessionId, err := sd.Store.GetUploadSessionId()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, sessionId)
}

func (sd *PallidSturgeonHandler) Upload(c echo.Context) error {
	var err error
	uploads := models.Upload{}
	if err := c.Bind(&uploads); err != nil {
		return err
	}

	sessionId, err := sd.Store.GetUploadSessionId()
	if err != nil {
		return err
	}

	for _, uploadSite := range uploads.SiteUpload {
		uploadSite.LastUpdated = time.Now()
		uploadSite.UploadedBy = "DeeLiang"
		uploadSite.UploadSessionId = sessionId
		err = sd.Store.SaveSiteUpload(uploadSite)
		if err != nil {
			return err
		}
	}

	for _, uploadFish := range uploads.FishUpload {
		uploadFish.LastUpdated = time.Now()
		uploadFish.UploadedBy = "DeeLiang"
		uploadFish.UploadSessionId = sessionId
		err = sd.Store.SaveFishUpload(uploadFish)
		if err != nil {
			return err
		}
	}

	for _, uploadSearch := range uploads.SearchUpload {
		uploadSearch.LastUpdated = time.Now()
		uploadSearch.UploadedBy = "DeeLiang"
		uploadSearch.UploadSessionId = sessionId
		err = sd.Store.SaveSearchUpload(uploadSearch)
		if err != nil {
			return err
		}
	}

	for _, uploadSupplemental := range uploads.UploadSupplemental {
		uploadSupplemental.LastUpdated = time.Now()
		uploadSupplemental.UploadedBy = "DeeLiang"
		uploadSupplemental.UploadSessionId = sessionId
		err = sd.Store.SaveSupplementalUpload(uploadSupplemental)
		if err != nil {
			return err
		}
	}
	for _, uploadProcedure := range uploads.ProcedureUpload {
		uploadProcedure.LastUpdated = time.Now()
		uploadProcedure.UploadedBy = "DeeLiang"
		uploadProcedure.UploadSessionId = sessionId
		err = sd.Store.SaveProcedureUpload(uploadProcedure)
		if err != nil {
			return err
		}
	}

	for _, uploadMoriver := range uploads.MoriverUpload {
		uploadMoriver.LastUpdated = time.Now()
		uploadMoriver.UploadedBy = "DeeLiang"
		uploadMoriver.UploadSessionId = sessionId
		err = sd.Store.SaveMoriverUpload(uploadMoriver)
		if err != nil {
			return err
		}
	}

	for _, uploadTelemetry := range uploads.TelemetryUpload {
		uploadTelemetry.LastUpdated = time.Now()
		uploadTelemetry.UploadedBy = "DeeLiang"
		uploadTelemetry.UploadSessionId = sessionId
		err = sd.Store.SaveTelemetryUpload(uploadTelemetry)
		if err != nil {
			return err
		}
	}

	procedureOut, err := sd.Store.CallStoreProcedures("DeeLiang", sessionId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, procedureOut)
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
