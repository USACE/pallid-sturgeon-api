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
	return c.String(http.StatusOK, "Pallid Sturgeon API v0.02")
}

func (sd *PallidSturgeonHandler) GetProjects(c echo.Context) error {

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	projects, err := sd.Store.GetProjects(userInfo.OfficeCode)
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

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	segments, err := sd.Store.GetSegments(userInfo.OfficeCode)
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

func (sd *PallidSturgeonHandler) GetBendRn(c echo.Context) error {
	bends, err := sd.Store.GetBendRn()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, bends)
}

func (sd *PallidSturgeonHandler) GetMeso(c echo.Context) error {
	macro := c.QueryParam("macro")
	mesoItems, err := sd.Store.GetMeso(macro)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, mesoItems)
}

func (sd *PallidSturgeonHandler) GetStructureFlow(c echo.Context) error {
	microStructure := c.QueryParam("microStructure")
	structureFlowItems, err := sd.Store.GetStructureFlow(microStructure)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, structureFlowItems)
}

func (sd *PallidSturgeonHandler) GetStructureMod(c echo.Context) error {
	structureFlow := c.QueryParam("structureFlow")
	structureModItems, err := sd.Store.GetStructureMod(structureFlow)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, structureModItems)
}

func (sd *PallidSturgeonHandler) GetSpecies(c echo.Context) error {
	species, err := sd.Store.GetSpecies()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, species)
}

func (sd *PallidSturgeonHandler) GetFtPrefixes(c echo.Context) error {
	ftPrefixes, err := sd.Store.GetFtPrefixes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, ftPrefixes)
}

func (sd *PallidSturgeonHandler) GetMr(c echo.Context) error {
	mr, err := sd.Store.GetMr()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, mr)
}

func (sd *PallidSturgeonHandler) GetOtolith(c echo.Context) error {
	otolith, err := sd.Store.GetOtolith()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, otolith)
}

func (sd *PallidSturgeonHandler) GetHeaderData(c echo.Context) error {
	siteId := c.QueryParam("siteId")
	headerDataItems, err := sd.Store.GetHeaderData(siteId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, headerDataItems)
}

func (sd *PallidSturgeonHandler) GetSiteDataEntries(c echo.Context) error {
	year, projectCode, segmentCode, seasonCode, bendrn, siteId := c.QueryParam("year"), c.QueryParam("projectCode"), c.QueryParam("segmentCode"), c.QueryParam("seasonCode"), c.QueryParam("bendrn"), c.QueryParam("siteId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	siteDataEntries, err := sd.Store.GetSiteDataEntries(siteId, year, userInfo.OfficeCode, projectCode, segmentCode, seasonCode, bendrn, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, siteDataEntries)
}

func (sd *PallidSturgeonHandler) SaveSiteDataEntry(c echo.Context) error {
	siteData := models.Sites{}
	if err := c.Bind(&siteData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	siteData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	siteData.UploadedBy = user.FirstName + " " + user.LastName
	id, err := sd.Store.SaveSiteDataEntry(siteData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateSiteDataEntry(c echo.Context) error {

	siteData := models.Sites{}
	if err := c.Bind(&siteData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	siteData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	siteData.UploadedBy = user.FirstName + " " + user.LastName
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

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetFishDataEntries(tableId, fieldId, mrId, userInfo.OfficeCode, queryParams)
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
	user := c.Get("PSUSER").(models.User)
	fishData.UploadedBy = user.FirstName + " " + user.LastName
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
	user := c.Get("PSUSER").(models.User)
	fishData.UploadedBy = user.FirstName + " " + user.LastName
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

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetMoriverDataEntries(tableId, fieldId, userInfo.OfficeCode, queryParams)
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
	user := c.Get("PSUSER").(models.User)
	moriverData.UploadedBy = user.FirstName + " " + user.LastName
	// @TODO: remove hardcoded value and fix date formatting
	moriverData.SetDate = "12-Dec-2022"
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
	user := c.Get("PSUSER").(models.User)
	moriverData.UploadedBy = user.FirstName + " " + user.LastName
	// @TODO: remove hardcoded value and fix date formatting
	moriverData.SetDate = "12-Dec-2022"
	err := sd.Store.UpdateMoriverDataEntry(moriverData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetSupplementalDataEntries(c echo.Context) error {
	tableId, fieldId, geneticsVial, pitTag, mrId, fId := c.QueryParam("tableId"), c.QueryParam("fieldId"), c.QueryParam("geneticsVial"), c.QueryParam("pitTag"), c.QueryParam("mrId"), c.QueryParam("fId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetSupplementalDataEntries(tableId, fieldId, geneticsVial, pitTag, mrId, fId, userInfo.OfficeCode, queryParams)
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
	user := c.Get("PSUSER").(models.User)
	supplementalData.UploadedBy = user.FirstName + " " + user.LastName
	err := sd.Store.SaveSupplementalDataEntry(supplementalData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, err)
}

func (sd *PallidSturgeonHandler) UpdateSupplementalDataEntry(c echo.Context) error {

	supplementalData := models.UploadSupplemental{}
	if err := c.Bind(&supplementalData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	supplementalData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	supplementalData.UploadedBy = user.FirstName + " " + user.LastName
	err := sd.Store.UpdateSupplementalDataEntry(supplementalData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetSearchDataEntries(c echo.Context) error {
	tableId, siteId := c.QueryParam("tableId"), c.QueryParam("siteId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetSearchDataEntries(tableId, siteId, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveSearchDataEntry(c echo.Context) error {
	searchData := models.UploadSearch{}
	if err := c.Bind(&searchData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	searchData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	searchData.UploadedBy = user.FirstName + " " + user.LastName
	// @TODO: remove hardcoded value and fix date formatting
	searchData.SearchDate = "12-Dec-2022"
	err := sd.Store.SaveSearchDataEntry(searchData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, err)
}

func (sd *PallidSturgeonHandler) UpdateSearchDataEntry(c echo.Context) error {
	searchData := models.UploadSearch{}
	if err := c.Bind(&searchData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	searchData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	searchData.UploadedBy = user.FirstName + " " + user.LastName
	// @TODO: remove hardcoded value and fix date formatting
	searchData.SearchDate = "12-Dec-2022"
	err := sd.Store.UpdateSearchDataEntry(searchData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetProcedureDataEntries(c echo.Context) error {
	tableId, fId := c.QueryParam("tableId"), c.QueryParam("fId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)
	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetProcedureDataEntries(tableId, fId, userInfo.OfficeCode, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveProcedureDataEntry(c echo.Context) error {
	procedureData := models.UploadProcedure{}
	if err := c.Bind(&procedureData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	procedureData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	procedureData.UploadedBy = user.FirstName + " " + user.LastName
	// @TODO: remove hardcoded value and fix date formatting
	procedureData.ProcedureDate = "12-Dec-2022"
	procedureData.DstStartDate = "12-Dec-2022"
	err := sd.Store.SaveProcedureDataEntry(procedureData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, err)
}

func (sd *PallidSturgeonHandler) UpdateProcedureDataEntry(c echo.Context) error {
	procedureData := models.UploadProcedure{}
	if err := c.Bind(&procedureData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	procedureData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	procedureData.UploadedBy = user.FirstName + " " + user.LastName
	// @TODO: remove hardcoded value and fix date formatting
	procedureData.ProcedureDate = "12-Dec-2022"
	procedureData.DstStartDate = "12-Dec-2022"
	err := sd.Store.UpdateProcedureDataEntry(procedureData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetTelemetryDataEntries(c echo.Context) error {
	tableId, seId := c.QueryParam("tableId"), c.QueryParam("seId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)
	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetTelemetryDataEntries(tableId, seId, userInfo.OfficeCode, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) SaveTelemetryDataEntry(c echo.Context) error {
	telemetryData := models.UploadTelemetry{}
	if err := c.Bind(&telemetryData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	telemetryData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	telemetryData.UploadedBy = user.FirstName + " " + user.LastName
	err := sd.Store.SaveTelemetryDataEntry(telemetryData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, err)
}

func (sd *PallidSturgeonHandler) UpdateTelemetryDataEntry(c echo.Context) error {
	telemetryData := models.UploadTelemetry{}
	if err := c.Bind(&telemetryData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	telemetryData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	telemetryData.UploadedBy = user.FirstName + " " + user.LastName
	err := sd.Store.UpdateTelemetryDataEntry(telemetryData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetFullFishDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	fileName, err := sd.Store.GetFullFishDataSummary(year, userInfo.OfficeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetFishDataSummary(c echo.Context) error {
	year, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetFishDataSummary(year, userInfo.OfficeCode, userInfo.ProjectCode, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullSuppDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	fileName, err := sd.Store.GetFullSuppDataSummary(year, userInfo.OfficeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetSuppDataSummary(c echo.Context) error {
	year, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetSuppDataSummary(year, userInfo.OfficeCode, userInfo.ProjectCode, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullMissouriDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	fileName, err := sd.Store.GetFullMissouriDataSummary(year, userInfo.OfficeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetMissouriDataSummary(c echo.Context) error {
	year, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetMissouriDataSummary(year, userInfo.OfficeCode, userInfo.ProjectCode, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullGeneticDataSummary(c echo.Context) error {
	year, project, fromDate, toDate, broodstock, hatchwild, speciesId, archive := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("fromDate"), c.QueryParam("toDate"), c.QueryParam("broodstock"), c.QueryParam("hatchwild"), c.QueryParam("speciesId"), c.QueryParam("archive")

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	fileName, err := sd.Store.GetFullGeneticDataSummary(year, userInfo.OfficeCode, project, fromDate, toDate, broodstock, hatchwild, speciesId, archive)
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

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetGeneticDataSummary(year, userInfo.OfficeCode, project, fromDate, toDate, broodstock, hatchwild, speciesId, archive, queryParams)
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

	dataSummary, err := sd.Store.GetTelemetryDataSummary(year, "ZZ", project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullTelemetryDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	fileName, err := sd.Store.GetFullTelemetryDataSummary(year, userInfo.OfficeCode, project, approved, season, spice, month, fromDate, toDate)
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

	dataSummary, err := sd.Store.GetProcedureDataSummary(year, "ZZ", project, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullProcedureDataSummary(c echo.Context) error {
	year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	fileName, err := sd.Store.GetFullProcedureDataSummary(year, userInfo.OfficeCode, project, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetMissouriDatasheetById(c echo.Context) error {
	siteId, project, segment, season, bend := c.QueryParam("siteId"), c.QueryParam("project"), c.QueryParam("segment"), c.QueryParam("season"), c.QueryParam("bend")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	missouriData, err := sd.Store.GetMissouriDatasheetById(siteId, userInfo.OfficeCode, project, segment, season, bend, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, missouriData)
}

func (sd *PallidSturgeonHandler) GetSearchDatasheetById(c echo.Context) error {
	siteId := c.QueryParam("siteId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	searchData, err := sd.Store.GetSearchDatasheetById(siteId, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, searchData)
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

	user := c.Get("PSUSER").(models.User)

	for _, uploadSite := range uploads.SiteUpload.Items {
		uploadSite.LastUpdated = time.Now()
		uploadSite.UploadedBy = user.FirstName + " " + user.LastName
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
		uploadFish.UploadedBy = user.FirstName + " " + user.LastName
		uploadFish.UploadSessionId = sessionId
		uploadFish.EditInitials = uploads.EditInitials
		uploadFish.UploadFilename = uploads.FishUpload.UploadFilename
		err = sd.Store.SaveFishUpload(uploadFish)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadSearch := range uploads.SearchUpload.Items {
		uploadSearch.SearchDateTime = processTimeString(uploadSearch.SearchDate)
		uploadSearch.LastUpdated = time.Now()
		uploadSearch.UploadedBy = user.FirstName + " " + user.LastName
		uploadSearch.UploadSessionId = sessionId
		uploadSearch.EditInitials = uploads.EditInitials
		uploadSearch.UploadFilename = uploads.SearchUpload.UploadFilename
		err = sd.Store.SaveSearchUpload(uploadSearch)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadSupplemental := range uploads.SupplementalUpload.Items {
		uploadSupplemental.LastUpdated = time.Now()
		uploadSupplemental.UploadedBy = user.FirstName + " " + user.LastName
		uploadSupplemental.UploadSessionId = sessionId
		uploadSupplemental.EditInitials = uploads.EditInitials
		uploadSupplemental.UploadFilename = uploads.SupplementalUpload.UploadFilename
		err = sd.Store.SaveSupplementalUpload(uploadSupplemental)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadProcedure := range uploads.ProcedureUpload.Items {
		uploadProcedure.ProcedureDateTime = processTimeString(uploadProcedure.ProcedureDate)
		uploadProcedure.DstStartDateTime = processTimeString(uploadProcedure.DstStartDate)
		uploadProcedure.LastUpdated = time.Now()
		uploadProcedure.UploadedBy = user.FirstName + " " + user.LastName
		uploadProcedure.UploadSessionId = sessionId
		uploadProcedure.EditInitials = uploads.EditInitials
		uploadProcedure.UploadFilename = uploads.ProcedureUpload.UploadFilename
		err = sd.Store.SaveProcedureUpload(uploadProcedure)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadMoriver := range uploads.MoriverUpload.Items {
		uploadMoriver.SetDateTime = processTimeString(uploadMoriver.SetDate)
		uploadMoriver.LastUpdated = time.Now()
		uploadMoriver.UploadedBy = user.FirstName + " " + user.LastName
		uploadMoriver.UploadSessionId = sessionId
		uploadMoriver.EditInitials = uploads.EditInitials
		uploadMoriver.UploadFilename = uploads.MoriverUpload.UploadFilename
		err = sd.Store.SaveMoriverUpload(uploadMoriver)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	for _, uploadTelemetry := range uploads.TelemetryUpload.Items {
		uploadTelemetry.LastUpdated = time.Now()
		uploadTelemetry.UploadedBy = user.FirstName + " " + user.LastName
		uploadTelemetry.UploadSessionId = sessionId
		uploadTelemetry.EditInitials = uploads.EditInitials
		uploadTelemetry.UploadFilename = uploads.TelemetryUpload.UploadFilename
		err = sd.Store.SaveTelemetryUpload(uploadTelemetry)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	procedureOut, err := sd.Store.CallStoreProcedures(user.FirstName+" "+user.LastName, sessionId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, procedureOut)
}

func (sd *PallidSturgeonHandler) CallStoreProcedures(c echo.Context) error {
	var err error
	uploadSessionId := c.Param("uploadSessionId")
	id, err := strconv.Atoi(uploadSessionId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)
	procedureOut, err := sd.Store.CallStoreProcedures(user.FirstName+" "+user.LastName, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, procedureOut)
}

func (sd *PallidSturgeonHandler) GetErrorCount(c echo.Context) error {

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	errorCounts, err := sd.Store.GetErrorCount(userInfo.OfficeCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, errorCounts)
}

func (sd *PallidSturgeonHandler) GetOfficeErrorLogs(c echo.Context) error {

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	officeErrorLogs, err := sd.Store.GetOfficeErrorLogs(userInfo.OfficeCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, officeErrorLogs)
}

func (sd *PallidSturgeonHandler) GetUsgNoVialNumbers(c echo.Context) error {

	user := c.Get("PSUSER").(models.User)

	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	usgNoVialNumbers, err := sd.Store.GetUsgNoVialNumbers(userInfo.OfficeCode, userInfo.ProjectCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, usgNoVialNumbers)
}

func (sd *PallidSturgeonHandler) GetUnapprovedDataSheets(c echo.Context) error {
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)
	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	unapprovedDataSheets, err := sd.Store.GetUnapprovedDataSheets(userInfo.ProjectCode, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, unapprovedDataSheets)
}

func (sd *PallidSturgeonHandler) GetBafiDataSheets(c echo.Context) error {
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)
	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	bafiDataSheets, err := sd.Store.GetBafiDataSheets(userInfo.OfficeCode, userInfo.ProjectCode, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, bafiDataSheets)
}

func (sd *PallidSturgeonHandler) GetUncheckedDataSheets(c echo.Context) error {
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	user := c.Get("PSUSER").(models.User)
	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	uncheckedDataSheets, err := sd.Store.GetUncheckedDataSheets(userInfo.OfficeCode, userInfo.ProjectCode, queryParams)
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
