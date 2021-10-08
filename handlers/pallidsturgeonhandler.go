package handlers

import (
	"net/http"
	"os"
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

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	projects, err := sd.Store.GetProjects("MO")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, projects)
}

func (sd *PallidSturgeonHandler) GetRoles(c echo.Context) error {
	roles, err := sd.Store.GetRoles()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, roles)
}

func (sd *PallidSturgeonHandler) GetFieldOffices(c echo.Context) error {
	fieldOffices, err := sd.Store.GetFieldOffices()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, fieldOffices)
}

func (sd *PallidSturgeonHandler) GetSeasons(c echo.Context) error {
	seasons, err := sd.Store.GetSeasons()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, seasons)
}

func (sd *PallidSturgeonHandler) GetSampleMethods(c echo.Context) error {
	sampleMethods, err := sd.Store.GetSampleMethods()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sampleMethods)
}

func (sd *PallidSturgeonHandler) GetSampleUnitTypes(c echo.Context) error {
	sampleUnitTypes, err := sd.Store.GetSampleUnitTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sampleUnitTypes)
}

func (sd *PallidSturgeonHandler) GetSegments(c echo.Context) error {

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	segments, err := sd.Store.GetSegments("MO")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, segments)
}

func (sd *PallidSturgeonHandler) GetBends(c echo.Context) error {
	bends, err := sd.Store.GetBends()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, bends)
}

func (sd *PallidSturgeonHandler) GetSiteDataEntries(c echo.Context) error {
	year, projectCode, segmentCode, seasonCode, bendrn := c.QueryParam("year"), c.QueryParam("projectCode"), c.QueryParam("segmentCode"), c.QueryParam("seasonCode"), c.QueryParam("bendrn")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	dataSummary, err := sd.Store.GetSiteDataEntries(year, "MO", projectCode, segmentCode, seasonCode, bendrn, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveSiteDataEntry(c echo.Context) error {
	siteData := models.UploadSite{}
	if err := c.Bind(&siteData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	siteData.LastUpdated = time.Now()
	// user := c.Get("PSUSER").(models.User)
	// siteData.UploadedBy = "Dee Liang"

	siteData.UploadedBy = "Dee Liang"
	id, err := sd.Store.SaveSiteDataEntry(siteData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateSiteDataEntry(c echo.Context) error {

	siteData := models.UploadSite{}
	if err := c.Bind(&siteData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	siteData.LastUpdated = time.Now()
	// user := c.Get("PSUSER").(models.User)
	// siteData.UploadedBy = "Dee Liang"
	siteData.UploadedBy = "Dee Liang"
	err := sd.Store.UpdateSiteDataEntry(siteData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetFishDataEntries(c echo.Context) error {
	tableId, fieldId, mrId := c.QueryParam("tableId"), c.QueryParam("fieldId"), c.QueryParam("mrId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	dataSummary, err := sd.Store.GetFishDataEntries(tableId, fieldId, mrId, "MO", queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveFishDataEntry(c echo.Context) error {
	fishData := models.UploadFish{}
	if err := c.Bind(&fishData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fishData.LastUpdated = time.Now()
	// user := c.Get("PSUSER").(models.User)
	// fishData.UploadedBy = "Dee Liang"
	fishData.UploadedBy = "Dee Liang"
	id, err := sd.Store.SaveFishDataEntry(fishData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateFishDataEntry(c echo.Context) error {

	fishData := models.UploadFish{}
	if err := c.Bind(&fishData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fishData.LastUpdated = time.Now()
	//user := c.Get("PSUSER").(models.User)
	fishData.UploadedBy = "Dee Liang"
	err := sd.Store.UpdateFishDataEntry(fishData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetMoriverDataEntries(c echo.Context) error {
	tableId, fieldId := c.QueryParam("tableId"), c.QueryParam("fieldId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	dataSummary, err := sd.Store.GetMoriverDataEntries(tableId, fieldId, "MO", queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveMoriverDataEntry(c echo.Context) error {
	moriverData := models.UploadMoriver{}
	if err := c.Bind(&moriverData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	moriverData.LastUpdated = time.Now()
	// user := c.Get("PSUSER").(models.User)
	// moriverData.UploadedBy = "Dee Liang"
	moriverData.UploadedBy = "Dee Liang"
	id, err := sd.Store.SaveMoriverDataEntry(moriverData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateMoriverDataEntry(c echo.Context) error {

	moriverData := models.UploadMoriver{}
	if err := c.Bind(&moriverData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	moriverData.LastUpdated = time.Now()
	//user := c.Get("PSUSER").(models.User)
	moriverData.UploadedBy = "Dee Liang"
	err := sd.Store.UpdateMoriverDataEntry(moriverData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetSupplementalDataEntries(c echo.Context) error {
	tableId, fieldId, geneticsVial, pitTag, mrId := c.QueryParam("tableId"), c.QueryParam("fieldId"), c.QueryParam("geneticsVial"), c.QueryParam("pitTag"), c.QueryParam("mrId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetSupplementalDataEntries(tableId, fieldId, geneticsVial, pitTag, mrId, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveSupplementalDataEntry(c echo.Context) error {
	supplementalData := models.UploadSupplemental{}
	if err := c.Bind(&supplementalData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	supplementalData.LastUpdated = time.Now()
	// user := c.Get("PSUSER").(models.User)
	// supplementalData.UploadedBy = "Dee Liang"
	supplementalData.UploadedBy = "Dee Liang"
	id, err := sd.Store.SaveSupplementalDataEntry(supplementalData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateSupplementalDataEntry(c echo.Context) error {

	supplementalData := models.UploadSupplemental{}
	if err := c.Bind(&supplementalData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	supplementalData.LastUpdated = time.Now()
	// user := c.Get("PSUSER").(models.User)
	// supplementalData.UploadedBy = "Dee Liang"
	supplementalData.UploadedBy = "Dee Liang"
	err := sd.Store.UpdateSupplementalDataEntry(supplementalData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetFullFishDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	fileName, err := sd.Store.GetFullFishDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetFishDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	dataSummary, err := sd.Store.GetFishDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullSuppDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	fileName, err := sd.Store.GetFullSuppDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetSuppDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	dataSummary, err := sd.Store.GetSuppDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullMissouriDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	fileName, err := sd.Store.GetFullMissouriDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetMissouriDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	dataSummary, err := sd.Store.GetMissouriDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullGeneticDataSummary(c echo.Context) error {
	year, project, fromDate, toDate, broodstock, hatchwild, speciesId, archive := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("fromDate"), c.QueryParam("toDate"), c.QueryParam("broodstock"), c.QueryParam("hatchwild"), c.QueryParam("speciesId"), c.QueryParam("archive")

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	fileName, err := sd.Store.GetFullGeneticDataSummary(year, "MO", project, fromDate, toDate, broodstock, hatchwild, speciesId, archive)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetGeneticDataSummary(c echo.Context) error {
	year, project, fromDate, toDate, broodstock, hatchwild, speciesId, archive := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("fromDate"), c.QueryParam("toDate"), c.QueryParam("broodstock"), c.QueryParam("hatchwild"), c.QueryParam("speciesId"), c.QueryParam("archive")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	dataSummary, err := sd.Store.GetGeneticDataSummary(year, "MO", project, fromDate, toDate, broodstock, hatchwild, speciesId, archive, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullSearchDataSummary(c echo.Context) error {

	fileName, err := sd.Store.GetFullSearchDataSummary()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetSearchDataSummary(c echo.Context) error {
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	dataSummary, err := sd.Store.GetSearchDataSummary(queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetTelemetryDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	dataSummary, err := sd.Store.GetTelemetryDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullTelemetryDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	fileName, err := sd.Store.GetFullTelemetryDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetProcedureDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	dataSummary, err := sd.Store.GetProcedureDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullProcedureDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }

	fileName, err := sd.Store.GetFullProcedureDataSummary(year, "MO", project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetUploadSessionId(c echo.Context) error {
	sessionId, err := sd.Store.GetUploadSessionId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sessionId)
}

func (sd *PallidSturgeonHandler) Upload(c echo.Context) error {
	var err error
	uploads := models.Upload{}
	if err := c.Bind(&uploads); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	sessionId, err := sd.Store.GetUploadSessionId()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//user := c.Get("PSUSER").(models.User)

	for _, uploadSite := range uploads.SiteUpload.Items {
		uploadSite.LastUpdated = time.Now()
		//uploadSite.UploadedBy = "Dee Liang"
		uploadSite.UploadedBy = "Dee Liang"
		uploadSite.UploadSessionId = sessionId
		uploadSite.EditInitials = uploads.EditInitials
		uploadSite.UploadFilename = uploads.SiteUpload.UploadFilename
		err = sd.Store.SaveSiteUpload(uploadSite)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadFish := range uploads.FishUpload.Items {
		uploadFish.LastUpdated = time.Now()
		uploadFish.UploadedBy = "Dee Liang"
		uploadFish.UploadSessionId = sessionId
		uploadFish.EditInitials = uploads.EditInitials
		uploadFish.UploadFilename = uploads.SiteUpload.UploadFilename
		err = sd.Store.SaveFishUpload(uploadFish)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadSearch := range uploads.SearchUpload.Items {
		uploadSearch.LastUpdated = time.Now()
		uploadSearch.UploadedBy = "Dee Liang"
		uploadSearch.UploadSessionId = sessionId
		uploadSearch.EditInitials = uploads.EditInitials
		uploadSearch.UploadFilename = uploads.SiteUpload.UploadFilename
		err = sd.Store.SaveSearchUpload(uploadSearch)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadSupplemental := range uploads.UploadSupplemental.Items {
		uploadSupplemental.LastUpdated = time.Now()
		uploadSupplemental.UploadedBy = "Dee Liang"
		uploadSupplemental.UploadSessionId = sessionId
		uploadSupplemental.EditInitials = uploads.EditInitials
		uploadSupplemental.UploadFilename = uploads.SiteUpload.UploadFilename
		err = sd.Store.SaveSupplementalUpload(uploadSupplemental)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	for _, uploadProcedure := range uploads.ProcedureUpload.Items {
		uploadProcedure.LastUpdated = time.Now()
		uploadProcedure.UploadedBy = "Dee Liang"
		uploadProcedure.UploadSessionId = sessionId
		uploadProcedure.EditInitials = uploads.EditInitials
		uploadProcedure.UploadFilename = uploads.SiteUpload.UploadFilename
		err = sd.Store.SaveProcedureUpload(uploadProcedure)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadMoriver := range uploads.MoriverUpload.Items {
		uploadMoriver.LastUpdated = time.Now()
		uploadMoriver.UploadedBy = "Dee Liang"
		uploadMoriver.UploadSessionId = sessionId
		uploadMoriver.EditInitials = uploads.EditInitials
		uploadMoriver.UploadFilename = uploads.SiteUpload.UploadFilename
		err = sd.Store.SaveMoriverUpload(uploadMoriver)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadTelemetry := range uploads.TelemetryUpload.Items {
		uploadTelemetry.LastUpdated = time.Now()
		uploadTelemetry.UploadedBy = "Dee Liang"
		uploadTelemetry.UploadSessionId = sessionId
		uploadTelemetry.EditInitials = uploads.EditInitials
		uploadTelemetry.UploadFilename = uploads.SiteUpload.UploadFilename
		err = sd.Store.SaveTelemetryUpload(uploadTelemetry)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	procedureOut, err := sd.Store.CallStoreProcedures("Dee Liang", sessionId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, procedureOut)
}

func (sd *PallidSturgeonHandler) SiteUpload(c echo.Context) error {
	var err error
	uploadSites := []models.UploadSite{}
	if err := c.Bind(&uploadSites); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, uploadSite := range uploadSites {
		uploadSite.LastUpdated = time.Now()
		// user := c.Get("PSUSER").(models.User)
		// uploadSite.UploadedBy = "Dee Liang"
		uploadSite.UploadedBy = "Dee Liang"
		err = sd.Store.SaveSiteUpload(uploadSite)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) FishUpload(c echo.Context) error {
	var err error
	uploadFishs := []models.UploadFish{}
	if err := c.Bind(&uploadFishs); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, uploadFish := range uploadFishs {
		uploadFish.LastUpdated = time.Now()
		// user := c.Get("PSUSER").(models.User)
		// uploadFish.UploadedBy = "Dee Liang"
		uploadFish.UploadedBy = "Dee Liang"
		err = sd.Store.SaveFishUpload(uploadFish)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) SearchUpload(c echo.Context) error {
	var err error
	uploadSearches := []models.UploadSearch{}
	if err := c.Bind(&uploadSearches); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, uploadSearch := range uploadSearches {
		uploadSearch.LastUpdated = time.Now()
		// user := c.Get("PSUSER").(models.User)
		// uploadSearch.UploadedBy = "Dee Liang"
		uploadSearch.UploadedBy = "Dee Liang"
		err = sd.Store.SaveSearchUpload(uploadSearch)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) SupplementalUpload(c echo.Context) error {
	var err error
	uploadSupplementals := []models.UploadSupplemental{}
	if err := c.Bind(&uploadSupplementals); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, uploadSupplemental := range uploadSupplementals {
		uploadSupplemental.LastUpdated = time.Now()
		// user := c.Get("PSUSER").(models.User)
		// uploadSupplemental.UploadedBy = "Dee Liang"
		uploadSupplemental.UploadedBy = "Dee Liang"
		err = sd.Store.SaveSupplementalUpload(uploadSupplemental)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) ProcedureUpload(c echo.Context) error {
	var err error
	uploadProcedures := []models.UploadProcedure{}
	if err := c.Bind(&uploadProcedures); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, uploadProcedure := range uploadProcedures {
		uploadProcedure.LastUpdated = time.Now()
		// user := c.Get("PSUSER").(models.User)
		// uploadProcedure.UploadedBy = "Dee Liang"
		uploadProcedure.UploadedBy = "Dee Liang"
		err = sd.Store.SaveProcedureUpload(uploadProcedure)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) MoriverUpload(c echo.Context) error {
	var err error
	UploadMorivers := []models.UploadMoriver{}
	if err := c.Bind(&UploadMorivers); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, uploadMoriver := range UploadMorivers {
		uploadMoriver.LastUpdated = time.Now()
		// user := c.Get("PSUSER").(models.User)
		// uploadMoriver.UploadedBy = "Dee Liang"
		uploadMoriver.UploadedBy = "Dee Liang"
		err = sd.Store.SaveMoriverUpload(uploadMoriver)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) TelemetryUpload(c echo.Context) error {
	var err error
	uploadTelemetrys := []models.UploadTelemetry{}
	if err := c.Bind(&uploadTelemetrys); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	for _, uploadTelemetry := range uploadTelemetrys {
		uploadTelemetry.LastUpdated = time.Now()
		// user := c.Get("PSUSER").(models.User)
		// uploadTelemetry.UploadedBy = "Dee Liang"
		uploadTelemetry.UploadedBy = "Dee Liang"
		err = sd.Store.SaveTelemetryUpload(uploadTelemetry)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) CallStoreProcedures(c echo.Context) error {
	var err error
	uploadSessionId := c.Param("uploadSessionId")
	id, err := strconv.Atoi(uploadSessionId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	//user := c.Get("PSUSER").(models.User)
	//procedureOut, err := sd.Store.CallStoreProcedures(user.FirstName+" "+user.LastName, id)
	procedureOut, err := sd.Store.CallStoreProcedures("Dee Liang", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, procedureOut)
}

func (sd *PallidSturgeonHandler) GetErrorCount(c echo.Context) error {

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	errorCounts, err := sd.Store.GetErrorCount("MO")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, errorCounts)
}

func (sd *PallidSturgeonHandler) GetUsgNoVialNumbers(c echo.Context) error {

	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	usgNoVialNumbers, err := sd.Store.GetUsgNoVialNumbers("MO")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, usgNoVialNumbers)
}

func (sd *PallidSturgeonHandler) GetUnapprovedDataSheets(c echo.Context) error {

	usgNoVialNumbers, err := sd.Store.GetUnapprovedDataSheets()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, usgNoVialNumbers)
}

func (sd *PallidSturgeonHandler) GetUncheckedDataSheets(c echo.Context) error {

	queryParams, err := marshalQuery(c)
	// user := c.Get("PSUSER").(models.User)

	// userInfo, err := sd.Store.GetUser(user.Email)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	uncheckedDataSheets, err := sd.Store.GetUncheckedDataSheets("MO", queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, uncheckedDataSheets)
}

func (sd *PallidSturgeonHandler) GetDownloadInfo(c echo.Context) error {
	downloadInfo, err := sd.Store.GetDownloadInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, downloadInfo)
}

func (sd *PallidSturgeonHandler) UploadDownloadZip(c echo.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files := form.File["file"]

	downloadInfo, err := sd.Store.UploadDownloadZip(files[0])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, downloadInfo)
}

func (sd *PallidSturgeonHandler) GetDownloadZip(c echo.Context) error {

	downloadZipName, err := sd.Store.GetDownloadZip()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(downloadZipName)
	return c.Inline(downloadZipName, downloadZipName)
}
