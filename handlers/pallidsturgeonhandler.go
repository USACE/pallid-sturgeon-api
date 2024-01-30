package handlers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
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
	id := c.QueryParam("id")

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	projects, err := sd.Store.GetProjects(userInfo.OfficeCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, projects)
}

func (sd *PallidSturgeonHandler) GetProjectsFilter(c echo.Context) error {
	project := c.QueryParam("project")

	projects, err := sd.Store.GetProjectsFilter(project)
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
	showAll := c.QueryParam("showAll")
	fieldOffices, err := sd.Store.GetFieldOffices(showAll)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, fieldOffices)
}

func (sd *PallidSturgeonHandler) GetSeasons(c echo.Context) error {
	project := c.QueryParam("project")
	seasons, err := sd.Store.GetSeasons(project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, seasons)
}

func (sd *PallidSturgeonHandler) GetSampleUnitTypes(c echo.Context) error {
	sampleUnitTypes, err := sd.Store.GetSampleUnitTypes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, sampleUnitTypes)
}

func (sd *PallidSturgeonHandler) GetSegments(c echo.Context) error {
	office, project := c.QueryParam("office"), c.QueryParam("project")
	segments, err := sd.Store.GetSegments(office, project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, segments)
}

func (sd *PallidSturgeonHandler) GetSampleUnit(c echo.Context) error {
	sampleUnitType, segment := c.QueryParam("sampleUnitType"), c.QueryParam("segment")
	bends, err := sd.Store.GetSampleUnit(sampleUnitType, segment)
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

func (sd *PallidSturgeonHandler) GetSetSite1(c echo.Context) error {
	microstructure := c.QueryParam("microstructure")
	setSiteItems, err := sd.Store.GetSetSite1(microstructure)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, setSiteItems)
}

func (sd *PallidSturgeonHandler) GetSetSite2(c echo.Context) error {
	setsite1 := c.QueryParam("setsite1")
	setSiteItems, err := sd.Store.GetSetSite2(setsite1)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, setSiteItems)
}

func (sd *PallidSturgeonHandler) GetYears(c echo.Context) error {
	year, err := sd.Store.GetYears()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, year)
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
	id, year, segmentCode, seasonCode, bendrn, siteId := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("segmentCode"), c.QueryParam("seasonCode"), c.QueryParam("bendrn"), c.QueryParam("siteId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	siteDataEntries, err := sd.Store.GetSiteDataEntries(siteId, year, userInfo.OfficeCode, userInfo.ProjectCode, segmentCode, seasonCode, bendrn, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, siteDataEntries)
}

func (sd *PallidSturgeonHandler) SaveSiteDataEntry(c echo.Context) error {
	code, sampleUnitType, segment := c.QueryParam("code"), c.QueryParam("sampleUnitType"), c.QueryParam("segment")
	siteData := models.Sites{}
	if err := c.Bind(&siteData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	siteData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	siteData.UploadedBy = user.FirstName + " " + user.LastName
	id, err := sd.Store.SaveSiteDataEntry(code, sampleUnitType, segment, siteData)
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
	id, tableId, fieldId, mrId := c.QueryParam("id"), c.QueryParam("tableId"), c.QueryParam("fieldId"), c.QueryParam("mrId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
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

func (sd *PallidSturgeonHandler) DeleteFishDataEntry(c echo.Context) error {
	id := c.Param("id")

	err := sd.Store.DeleteFishDataEntry(id)
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
	moriverData.SetDate = processStringTime(DerefString(moriverData.SetDate), "app")
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
	moriverData.SetDate = processStringTime(DerefString(moriverData.SetDate), "app")
	err := sd.Store.UpdateMoriverDataEntry(moriverData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetSupplementalDataEntries(c echo.Context) error {
	id, tableId, fieldId, geneticsVial, pitTag, mrId, fId := c.QueryParam("id"), c.QueryParam("tableId"), c.QueryParam("fieldId"), c.QueryParam("geneticsVial"), c.QueryParam("pitTag"), c.QueryParam("mrId"), c.QueryParam("fId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
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
	user := c.Get("PSUSER").(models.User)
	supplementalData.UploadedBy = user.FirstName + " " + user.LastName
	err := sd.Store.UpdateSupplementalDataEntry(supplementalData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) DeleteSupplementalDataEntry(c echo.Context) error {
	id := c.Param("id")

	err := sd.Store.DeleteSupplementalDataEntry(id)
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
	searchData.SearchDate = processStringTime(DerefString(searchData.SearchDate), "app")
	id, err := sd.Store.SaveSearchDataEntry(searchData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateSearchDataEntry(c echo.Context) error {
	searchData := models.UploadSearch{}
	if err := c.Bind(&searchData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	searchData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	searchData.UploadedBy = user.FirstName + " " + user.LastName
	searchData.SearchDate = processStringTime(DerefString(searchData.SearchDate), "app")
	err := sd.Store.UpdateSearchDataEntry(searchData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetProcedureDataEntries(c echo.Context) error {
	id, tableId, fId, mrId := c.QueryParam("id"), c.QueryParam("tableId"), c.QueryParam("fId"), c.QueryParam("mrId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	dataSummary, err := sd.Store.GetProcedureDataEntries(tableId, fId, mrId, userInfo.OfficeCode, queryParams)
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
	procedureData.ProcedureDate = processStringTime(DerefString(procedureData.ProcedureDate), "app")
	procedureData.DstStartDate = processStringTime(DerefString(procedureData.DstStartDate), "app")
	id, err := sd.Store.SaveProcedureDataEntry(procedureData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, id)
}

func (sd *PallidSturgeonHandler) UpdateProcedureDataEntry(c echo.Context) error {
	procedureData := models.UploadProcedure{}
	if err := c.Bind(&procedureData); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	procedureData.LastUpdated = time.Now()
	user := c.Get("PSUSER").(models.User)
	procedureData.UploadedBy = user.FirstName + " " + user.LastName
	procedureData.ProcedureDate = processStringTime(DerefString(procedureData.ProcedureDate), "app")
	procedureData.DstStartDate = processStringTime(DerefString(procedureData.DstStartDate), "app")
	err := sd.Store.UpdateProcedureDataEntry(procedureData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) DeleteProcedureDataEntry(c echo.Context) error {
	id := c.Param("id")

	err := sd.Store.DeleteProcedureDataEntry(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "successfully deleted procedure data entry id "+id)
}

func (sd *PallidSturgeonHandler) GetTelemetryDataEntries(c echo.Context) error {
	id, tableId, seId := c.QueryParam("id"), c.QueryParam("tableId"), c.QueryParam("seId")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
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
	id, err := sd.Store.SaveTelemetryDataEntry(telemetryData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(200, id)
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

func (sd *PallidSturgeonHandler) DeleteTelemetryDataEntry(c echo.Context) error {
	id := c.Param("id")

	err := sd.Store.DeleteTelemetryDataEntry(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, `{"result":"success"}`)
}

func (sd *PallidSturgeonHandler) GetFullFishDataSummary(c echo.Context) error {
	id, year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	fileName, err := sd.Store.GetFullFishDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetFishDataSummary(c echo.Context) error {
	id, year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	dataSummary, err := sd.Store.GetFishDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullSuppDataSummary(c echo.Context) error {
	id, year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	fileName, err := sd.Store.GetFullSuppDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetSuppDataSummary(c echo.Context) error {
	id, year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	dataSummary, err := sd.Store.GetSuppDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullMissouriDataSummary(c echo.Context) error {
	id, project, year, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("project"), c.QueryParam("year"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	fileName, err := sd.Store.GetFullMissouriDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetMissouriDataSummary(c echo.Context) error {
	id, project, year, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("project"), c.QueryParam("year"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	dataSummary, err := sd.Store.GetMissouriDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullGeneticDataSummary(c echo.Context) error {
	id, year, project, fromDate, toDate, broodstock, hatchwild, speciesId, archive := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("fromDate"), c.QueryParam("toDate"), c.QueryParam("broodstock"), c.QueryParam("hatchwild"), c.QueryParam("speciesId"), c.QueryParam("archive")

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	fileName, err := sd.Store.GetFullGeneticDataSummary(year, userInfo.OfficeCode, projectVal, fromDate, toDate, broodstock, hatchwild, speciesId, archive)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetGeneticDataSummary(c echo.Context) error {
	id, year, project, fromDate, toDate, broodstock, hatchwild, speciesId, archive := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("fromDate"), c.QueryParam("toDate"), c.QueryParam("broodstock"), c.QueryParam("hatchwild"), c.QueryParam("speciesId"), c.QueryParam("archive")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	dataSummary, err := sd.Store.GetGeneticDataSummary(year, userInfo.OfficeCode, projectVal, fromDate, toDate, broodstock, hatchwild, speciesId, archive, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullSearchDataSummary(c echo.Context) error {
	id, year, project, approved, season, segment, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("segment"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	fileName, err := sd.Store.GetFullSearchDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, segment, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetSearchDataSummary(c echo.Context) error {
	id, year, project, approved, season, segment, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("segment"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	dataSummary, err := sd.Store.GetSearchDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, segment, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetTelemetryDataSummary(c echo.Context) error {
	id, year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	dataSummary, err := sd.Store.GetTelemetryDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullTelemetryDataSummary(c echo.Context) error {
	id, year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	fileName, err := sd.Store.GetFullTelemetryDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetProcedureDataSummary(c echo.Context) error {
	id, year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	dataSummary, err := sd.Store.GetProcedureDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, dataSummary)
}

func (sd *PallidSturgeonHandler) GetFullProcedureDataSummary(c echo.Context) error {
	id, year, project, approved, season, spice, month, fromDate, toDate := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("project"), c.QueryParam("approved"), c.QueryParam("season"), c.QueryParam("spice"), c.QueryParam("month"), c.QueryParam("fromDate"), c.QueryParam("toDate")

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	fileName, err := sd.Store.GetFullProcedureDataSummary(year, userInfo.OfficeCode, projectVal, approved, season, spice, month, fromDate, toDate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer os.Remove(fileName)
	return c.Inline(fileName, fileName)
}

func (sd *PallidSturgeonHandler) GetMissouriDatasheetById(c echo.Context) error {
	id, siteId, project, segment, season, bend := c.QueryParam("id"), c.QueryParam("siteId"), c.QueryParam("project"), c.QueryParam("segment"), c.QueryParam("season"), c.QueryParam("bend")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// set project
	projectVal := ""
	if userInfo.ProjectCode == "2" {
		projectVal = "2"
	} else {
		projectVal = project
	}

	missouriData, err := sd.Store.GetMissouriDatasheetById(siteId, userInfo.OfficeCode, projectVal, segment, season, bend, queryParams)
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
	// Retrieve single uploaded file from the request.
	file, err := c.FormFile("files")
	log.Printf("file: %v", file)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "Failed! No files were submitted",
			Status:  "Failed",
			Data:    nil,
		})
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "Failed! Unable to open the file",
			Status:  "Failed",
			Data:    nil,
		})
	}
	defer fileContent.Close()

	// Create the CSV reader
	csvReader := csv.NewReader(fileContent)
	csvReader.FieldsPerRecord = -1

	csvData, err := csvReader.ReadAll()
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "Failed! Unable to read CSV file",
			Status:  "Failed",
			Data:    nil,
		})
	}

	var moriverUpload models.UploadMoriver
	var moriverUploads []models.UploadMoriver

	for _, each := range csvData {
		moriverUpload.SiteID, _ = strconv.Atoi(each[0])
		moriverUpload.SiteFid = each[1]
		moriverUpload.MrFid = each[2]
		moriverUpload.SeFieldID = each[3]
		moriverUpload.Season = each[4]
		moriverUploads = append(moriverUploads, moriverUpload)
	}

	// Convert to JSON
	jsonData, err := json.Marshal(moriverUploads)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Message: "Failed! Unable to convert to JSON",
			Status:  "Failed",
			Data:    nil,
		})
	}

	fmt.Println(string(jsonData))

	return c.JSON(http.StatusBadRequest, &models.Response{
		Message: "Upload is finished :)",
		Status:  "Failed",
		Data:    nil,
	})
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
	id := c.QueryParam("id")

	userInfo, err := sd.Store.GetUserRoleById(id)
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
	id := c.QueryParam("id")

	userInfo, err := sd.Store.GetUserRoleById(id)
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
	id := c.QueryParam("id")

	userInfo, err := sd.Store.GetUserRoleById(id)
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
	id := c.QueryParam("id")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	unapprovedDataSheets, err := sd.Store.GetUnapprovedDataSheets(userInfo.ProjectCode, userInfo.OfficeCode, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, unapprovedDataSheets)
}

func (sd *PallidSturgeonHandler) GetBafiDataSheets(c echo.Context) error {
	id := c.QueryParam("id")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
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
	id := c.QueryParam("id")
	queryParams, err := marshalQuery(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userInfo, err := sd.Store.GetUserRoleById(id)
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

func (sd *PallidSturgeonHandler) GetUploadSessionLogs(c echo.Context) error {
	uploadSessionId := c.QueryParam("uploadSessionId")

	user := c.Get("PSUSER").(models.User)
	userInfo, err := sd.Store.GetUser(user.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	bends, err := sd.Store.GetUploadSessionLogs(userInfo.FirstName+" "+userInfo.LastName, uploadSessionId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, bends)
}

func (sd *PallidSturgeonHandler) GetSitesExport(c echo.Context) error {
	id, year, segmentCode, seasonCode, bendrn := c.QueryParam("id"), c.QueryParam("year"), c.QueryParam("segmentCode"), c.QueryParam("seasonCode"), c.QueryParam("bendrn")

	userInfo, err := sd.Store.GetUserRoleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	exportData, err := sd.Store.GetSitesExport(year, userInfo.OfficeCode, userInfo.ProjectCode, segmentCode, seasonCode, bendrn)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, exportData)
}
