# Instrumentation API

Table of Contents
- [Pallid Sturgeon API](#pallidsturgeon-api)
    - [projects](#projects)
    - [seasons](#seasons)
    - [segments](#segments)
    - [bends](#bends)
    - [siteDataEntry](#siteDataEntry)
    - [fishDataEntry](#fishDataEntry)
    - [moriverDataEntry](#moriverDataEntry)
    - [supplementalDataEntry](#supplementalDataEntry)
    - [upload](#upload)
    - [fishDataSummary](#fishDataSummary)
    - [suppDataSummary](#suppDataSummary)
    - [missouriDataSummary](#missouriDataSummary)
    - [geneticDataSummary](#geneticDataSummary)
    - [searchDataSummary](#searchDataSummary)
---
### Projects
- List Projects \
  [http://localhost:8080/psapi/projects](http://localhost:8080/psapi/projects)

---
### Seasons
- List Seasons \
  [http://localhost:8080/psapi/seasons](http://localhost:8080/psapi/seasons)

---
### Segments
- List SeasSegmentsons \
  [http://localhost:8080/psapi/segments](http://localhost:8080/psapi/segments)

---
### Bends
- List Bends \
  [http://localhost:8080/psapi/bends](http://localhost:8080/psapi/bends)

---
### SiteDataEntry
- List siteDataEntry \
  [http://localhost:8080/psapi/siteDataEntry?fieldId=20120827-031-01-01-01&orderby=f_id%20desc&page=0&size=5](http://localhost:8080/psapi/fishDataEntry?fieldId=20120827-031-01-01-01&orderby=f_id%20desc&page=0&size=5)
- Create siteDataEntry \
  `http://localhost:8080/psapi/siteDataEntry`
    - Example `POST` body
        ```
      {
        "siteFid": "F-1",
        "siteYear": 2013,
        "fieldOffice": "SD",
        "project": "1",
        "segment": "7",
        "season": "ST",
        "sampleUnitTypeCode": "B",
        "bendrn": "R",
        "editInitials": "DG",
        "comments": "changed year"
      }
- Update siteDataEntry \
  `http://localhost:8080/psapi/siteDataEntry`
    - Example `PUT` body
        ```
      {
        "siteId": 10122,
        "siteFid": "F-1",
        "siteYear": 2013,
        "fieldOffice": "SD",
        "project": "1",
        "segment": "7",
        "season": "ST",
        "sampleUnitTypeCode": "B",
        "bendrn": "R",
        "editInitials": "DG",
        "comments": "changed year2"
      }
---
### FishDataEntry
- List fishDataEntry \
  [http://localhost:8080/psapi/fishDataEntry?fieldId=20120827-031-01-01-01&orderby=f_id%20desc&page=0&size=5](http://localhost:8080/psapi/fishDataEntry?fieldId=20120827-031-01-01-01&orderby=f_id%20desc&page=0&size=5)
- Create fishDataEntry \
  `http://localhost:8080/psapi/fishDataEntry`
    - Example `POST` body
        ```
        {
          "id": null,
          "siteId": 0,
          "fieldOffice": "",
          "project": null,
          "segment": null,
          "uniqueID": null,
          "panelhook": "",
          "ffid": "20120827-031-01-01-01",
          "mrId": 950,
          "bait": "",
          "species": "WTBS",
          "length": 101,
          "weight": null,
          "fishcount": 1,
          "finCurl": "",
          "otolith": "",
          "rayspine": "",
          "scale": "",
          "ftprefix": "",
          "ftnum": "",
          "ftmr": "",
          "editInitials": "DG"
        }
- Update fishDataEntry \
  `http://localhost:8080/psapi/fishDataEntry`
    - Example `PUT` body
        ```
        {
          "fid": 3000031,
          "id": null,
          "siteId": 0,
          "fieldOffice": "",
          "project": null,
          "segment": null,
          "uniqueID": null,
          "panelhook": "",
          "ffid": "20120827-031-01-01-01",
          "mrId": 950,
          "bait": "",
          "species": "WTBS",
          "length": 101,
          "weight": null,
          "fishcount": 1,
          "finCurl": "",
          "otolith": "",
          "rayspine": "",
          "scale": "",
          "ftprefix": "",
          "ftnum": "",
          "ftmr": "",
          "editInitials": "DG"
        }
---
### MoriverDataEntry
- List moriverDataEntry \
  [http://localhost:8080/psapi/moriverDataEntry?fieldId=20120827-031-01-01-01&orderby=f_id%20desc&page=0&size=5](http://localhost:8080/psapi/moriverDataEntry?fieldId=20120827-031-01-01-01&orderby=f_id%20desc&page=0&size=5)
- Create moriverDataEntry \
  `http://localhost:8080/psapi/moriverDataEntry`
    - Example `POST` body
        ```
        {
			    "mrFid": "20121113-003-04-01-04",
				  "siteId": 493,
					"fieldOffice": "",
					"project": null,
				  "segment": null,
				  "season": "ST",
				  "setdate": "2012-11-13T00:00:00Z",
				  "subsample": 4,
				  "subsamplepass": 1,
				  "subsampleROrN": "R",
				  "recorder": "TWH",
				  "gear": "GN18",
				  "gearType": "S",
				  "temp": 5.8,
				  "turbidity": 23,
				  "conductivity": null,
				  "do": null,
				  "distance": null,
				  "width": null,
				  "netrivermile": null,
				  "structurenumber": "",
				  "usgs": "",
				  "riverstage": null,
				  "discharge": null,
				  "u1": "",
				  "u2": "",
				  "u3": "",
				  "u4": "",
				  "u5": "",
				  "u6": "",
				  "u7": "",
				  "macro": "ISB",
				  "meso": "POOL",
				  "habitatrn": "R",
				  "qc": "",
				  "microStructure": "4",
				  "structureFlow": "2",
				  "structureMod": "2",
				  "setSite_1": "1",
				  "setSite_2": "5",
				  "setSite_3": "0",
				  "startTime": "14:38:12",
				  "startLatitude": 40.97622,
				  "startLongitude": -95.82979,
				  "stopTime": "9:58:25",
				  "stopLatitude": null,
				  "stopLongitude": null,
				  "depth1": 4,
				  "velocitybot1": null,
				  "velocity08_1": null,
				  "velocity02or06_1": null,
				  "depth2": 4.3,
				  "velocitybot2": null,
				  "velocity08_2": null,
				  "velocity02or06_2": null,
				  "depth3": 4.2,
				  "velocitybot3": null,
				  "velocity08_3": null,
				  "velocity02or06_3": null,
				  "watervel": null,
				  "cobble": null,
				  "organic": null,
				  "silt": null,
				  "sand": null,
				  "gravel": null,
				  "comments": "no flow taken in eddy \r\n",
				  "complete": null,
				  "checkby": "",
				  "noTurbidity": "",
				  "noVelocity": "",
				  "editInitials": "DG"
				}
- Update moriverDataEntry \
  `http://localhost:8080/psapi/moriverDataEntry`
    - Example `PUT` body
        ```
        {
          "mrId": 300080,
			    "mrFid": "20121113-003-04-01-04",
				  "siteId": 493,
					"fieldOffice": "",
					"project": null,
				  "segment": null,
				  "season": "ST",
				  "setdate": "2012-11-13T00:00:00Z",
				  "subsample": 4,
				  "subsamplepass": 1,
				  "subsampleROrN": "R",
				  "recorder": "TWH",
				  "gear": "GN18",
				  "gearType": "S",
				  "temp": 5.8,
				  "turbidity": 23,
				  "conductivity": null,
				  "do": null,
				  "distance": null,
				  "width": null,
				  "netrivermile": null,
				  "structurenumber": "",
				  "usgs": "",
				  "riverstage": null,
				  "discharge": null,
				  "u1": "",
				  "u2": "",
				  "u3": "",
				  "u4": "",
				  "u5": "",
				  "u6": "",
				  "u7": "",
				  "macro": "ISB",
				  "meso": "POOL",
				  "habitatrn": "R",
				  "qc": "",
				  "microStructure": "4",
				  "structureFlow": "2",
				  "structureMod": "2",
				  "setSite_1": "1",
				  "setSite_2": "5",
				  "setSite_3": "0",
				  "startTime": "14:38:12",
				  "startLatitude": 40.97622,
				  "startLongitude": -95.82979,
				  "stopTime": "9:58:25",
				  "stopLatitude": null,
				  "stopLongitude": null,
				  "depth1": 4,
				  "velocitybot1": null,
				  "velocity08_1": null,
				  "velocity02or06_1": null,
				  "depth2": 4.3,
				  "velocitybot2": null,
				  "velocity08_2": null,
				  "velocity02or06_2": null,
				  "depth3": 4.2,
				  "velocitybot3": null,
				  "velocity08_3": null,
				  "velocity02or06_3": null,
				  "watervel": null,
				  "cobble": null,
				  "organic": null,
				  "silt": null,
				  "sand": null,
				  "gravel": null,
				  "comments": "no flow taken in eddy \r\n",
				  "complete": null,
				  "checkby": "",
				  "noTurbidity": "",
				  "noVelocity": "",
				  "editInitials": "DG"
				}
---
### SupplementalDataEntry
- List supplementalDataEntry \
  [http://localhost:8080/psapi/supplementalDataEntry?fieldId=20120827-031-01-01-01&orderby=f_id%20desc&page=0&size=5](http://localhost:8080/psapi/supplementalDataEntry?fieldId=20120827-031-01-01-01&orderby=f_id%20desc&page=0&size=5)
- Create supplementalDataEntry \
  `http://localhost:8080/psapi/supplementalDataEntry`
    - Example `POST` body
        ```
        {
          "fid": 11781,
          "fFid": "20121024-032-03-01-14-001",
          "mrId": "2412",
          "tagnumber": "4706162c09",
          "pitrn": "N",
          "scuteloc": "N",
          "scutenum": null,
          "scuteloc2": "",
          "scutenum2": null,
          "elhv": "",
          "elcolor": "N",
          "erhv": "",
          "ercolor": "N",
          "cwtyn": "N",
          "dangler": "N",
          "genetic": "Y",
          "geneticsVialNumber": "3860",
          "broodstock": null,
          "hatchWild": null,
          "speciesId": null,
          "archive": null,
          "head": null,
          "snouttomouth": null,
          "inter": null,
          "mouthwidth": null,
          "mIb": null,
          "lOb": null,
          "lIb": null,
          "rIb": null,
          "rOb": null,
          "anal": null,
          "dorsal": null,
          "status": "",
          "hatcheryOrigin": "",
          "sex": "",
          "stage": "",
          "recapture": "",
          "photo": "",
          "geneticNeeds": "",
          "otherTagInfo": "",
          "comments": "-\r\n",
          "editInitials": "DG"
        }
- Update supplementalDataEntry \
  `http://localhost:8080/psapi/supplementalDataEntry`
    - Example `PUT` body
        ```
        {
          "sid": 100041,
          "fid": 11781,
          "fFid": "20121024-032-03-01-14-001",
          "mrId": "2412",
          "tagnumber": "4706162c09",
          "pitrn": "N",
          "scuteloc": "N",
          "scutenum": null,
          "scuteloc2": "",
          "scutenum2": null,
          "elhv": "",
          "elcolor": "N",
          "erhv": "",
          "ercolor": "N",
          "cwtyn": "N",
          "dangler": "N",
          "genetic": "Y",
          "geneticsVialNumber": "3860",
          "broodstock": null,
          "hatchWild": null,
          "speciesId": null,
          "archive": null,
          "head": null,
          "snouttomouth": null,
          "inter": null,
          "mouthwidth": null,
          "mIb": null,
          "lOb": null,
          "lIb": null,
          "rIb": null,
          "rOb": null,
          "anal": null,
          "dorsal": null,
          "status": "",
          "hatcheryOrigin": "",
          "sex": "",
          "stage": "",
          "recapture": "",
          "photo": "",
          "geneticNeeds": "",
          "otherTagInfo": "",
          "comments": "-\r\n",
          "editInitials": "DG"
        }
---
### MoriverDataSummary
- List moriverDataSummary \
  [http://localhost:8080/psapi/moriverDataSummary?year=2021&officeCode=MO&project=1&season=MR&month=10&fromDate=10%2F06%2F2020&toDate=10%2F06%2F2020&orderby=mr_id%20desc&page=0&size=5](http://localhost:8080/psapi/moriverDataSummary?year=2021&officeCode=MO&project=1&season=MR&month=10&fromDate=10%2F06%2F2020&toDate=10%2F06%2F2020&orderby=mr_id%20desc&page=0&size=5)

---
### SuppDataSummary
- List suppDataSummary \
  [http://localhost:8080/psapi/suppDataSummary?year=2021&officeCode=MO&project=1&season=MR&month=10&fromDate=10%2F06%2F2020&toDate=10%2F06%2F2020&orderby=mr_id%20desc&page=0&size=5](http://localhost:8080/psapi/suppDataSummary?year=2021&officeCode=MO&project=1&season=MR&month=10&fromDate=10%2F06%2F2020&toDate=10%2F06%2F2020&orderby=mr_id%20desc&page=0&size=5)

---
### MissouriDataSummary
- List missouriDataSummary \
  [http://localhost:8080/psapi/missouriDataSummary?year=2021&officeCode=MO&project=1&season=MR&month=10&fromDate=10%2F06%2F2020&toDate=10%2F06%2F2020&orderby=mr_id%20desc&page=0&size=5](http://localhost:8080/psapi/missouriDataSummary?year=2021&officeCode=MO&project=1&season=MR&month=10&fromDate=10%2F06%2F2020&toDate=10%2F06%2F2020&orderby=mr_id%20desc&page=0&size=5)

---
### GeneticDataSummary
- List missouriDataSummary \
  [http://localhost:8080/psapi/geneticDataSummary?year=2020&officeCode=KC&project=2&fromDate=06%2F12%2F2020&toDate=06%2F12%2F2020&page=0&size=5](http://localhost:8080/psapi/geneticDataSummary?year=2020&officeCode=KC&project=2&fromDate=06%2F12%2F2020&toDate=06%2F12%2F2020&page=0&size=5)

---
### SearchDataSummary
- List missouriDataSummary \
  [http://localhost:8080/psapi/searchDataSummary?orderby=se_id%20desc&page=0&size=5](http://localhost:8080/psapi/searchDataSummary?orderby=se_id%20desc&page=0&size=5)

---
### Upload
- Upload \
  `http://localhost:8080/psapi/upload`
    - Example `POST` body
        ```
        {
          "editInitials": "DG",
          "siteUpload": {
            "uploadFilename": "pspa_sites_datasheet_20210617_1900_59.csv",
            "Items" :[
              {
                "siteId": 0,
                "siteFid": "F-1",
                "siteYear": 2021,
                "fieldofficeID": "SD",
                "fieldOffice": "KC - Kansas City",
                "projectId": 1,
                "project": "1 - Pallid Sturgeon Population Assessment",
                "segmentId": 28,
                "segment": "28 - Osage River",
                "seasonId": "A0",
                "season": "A0 - Age 0",
                "bend": 2,
                "bendrn": "N",
                "bendRiverMile": 4.8,
                "comments": "test"
              }
            ]
          },
          "fishUpload": {
            "uploadFilename": "fish_datasheet_20210617_1900_59.csv",
            "Items": [
              {
                "siteId": 0,
                "fFid": "20210617-184105056-001-001",
                "mrFid": "20210617-184105056-001",
                "panelhook": "1",
                "bait": "W",
                "species": "PDSG",
                "length": 2,
                "weight": 2,
                "fishcount": 1,
                "finCurl": "Y",
                "otolith": "D",
                "rayspine": "X",
                "scale": "X",
                "ftnum": "45678",
                "ftmr": "L",
                "ftprefix": "BC",
                "comments": "test"
              }
            ]
          },
          "searchUpload":{
            "uploadFilename": "search_effort_20210617_1900_59.csv",
            "Items": [
              {
                "siteId": 0,
                "seFid": "20210617-185747028-001",
                "dsId": 1,
                "siteFid": "F-1",
                "searchDate": "2021-06-17T00:00:00Z",
                "recorder": "NR",
                "searchTypeCode": "BS",
                "searchDay": 12345678,
                "startTime": "18:58:06",
                "startLatitude": 50,
                "startLongitude": -88,
                "stopTime": "18:58:08",
                "stopLatitud": 50,
                "stopLongitude": -88,
                "temp": 30,
                "conductivity": 22
              }
            ]
          },
          "telemetryUpload": {
            "uploadFilename": "telemetry_20210617_1900_59.csv",
            "Items": [
              {
                "tFid": "20210617-185747028-001-001",
                "seFid": "20210617-185747028-001",
                "bend": 2,
                "radioTagNum": 1234567890,
                "frequencyIdCode": 3,
                "captureTime": "18:59:49",
                "captureLatitude": 50,
                "captureLongitude": -88,
                "positionConfidence": 2,
                "macroId": "CONF",
                "mesoId": "CHNB",
                "depth": 1,
                "temp": 30,
                "conductivity": 1,
                "turbidity": 1,
                "silt": 1,
                "sand": 1,
                "gravel": 1,
                "comments": "comments"
              }
            ]
          },
          "procedureUpload": {
            "uploadFilename": "procedure_20210617_1900_59.csv",
            "Items": [
              {
                "f_fid": "20210617-184105056-001-001",
                "purposeCode": "RI",
                "procedurDate": "2021-06-17T00:00:00Z",
                "procedureStartTime": "1:00",
                "procedureEndTime": "2:00",
                "procedureBy": "NR",
                "antibioticInjectionInd": 0,
                "photoDorsalInd": 0,
                "photoVentralInd": 0,
                "photoLeftInd": 0,
                "oldRadioTagNum": 0,
                "oldFrequencyId": 0,
                "dstSerialNum": 12345,
                "dstStartDate": "2021-06-16T00:00:00Z",
                "dstStartTime": "10:30",
                "dstReimplantInd": 0,
                "newRadioTagNum": 1234,
                "newFrequencyId": 0,
                "sexCode": "s",
                "bloodSampleInd": 0,
                "eggSampleInd": 0,
                "comments": "comments",
                "fishHealthComments": "fishHealthComments",
                "evalLocationCode": "GP",
                "spawnCode": "PS",
                "visualReproStatusCode": "R",
                "ultrasoundReproStatusCode": "R",
                "expectedSpawnYear": 1,
                "ltrasoundGonadLength": 1,
                "gonadCondition": "1"
              }
            ]
          },
          "supplementalUpload": {
            "uploadFilename": "supplemental_20210617_1900_59.csv",
            "Items": [
              {
                "siteId": 0,
                "fFid": "20210617-184105056-001-001",
                "mrFid": "20210617-184105056-001",
                "tagnumber": "1234567890",
                "pitrn": "N",
                "scuteloc": "D",
                "scutenum": 2,
                "scuteloc2": "l",
                "scutenum2": 2,
                "elhv": "H",
                "elcolor": "G",
                "erhv": "H",
                "ercolor": "G",
                "cwtyn": "Y",
                "dangler": "Y",
                "genetic": "Y",
                "geneticsVialNumber": "STURG-1",
                "broodstock": 1,
                "hatchWild": 1,
                "speciesId": 1,
                "archive": 1,
                "head": 1,
                "snouttomouth": 1,
                "inter": 1,
                "outhwidth": 1,
                "mIb": 1,
                "lOb": 1,
                "lIb": 1,
                "rIb": 1,
                "rOb": 1,
                "anal": 1,
                "orsal": 1,
                "status": "H",
                "hatcheryOrigin": "H",
                "sex": "s",
                "stage": "g",
                "recapture": "c",
                "photo": "p",
                "geneticNeeds": "geneticNeeds",
                "otherTagInfo": "otherTagInfo",
                "comments": "comments"
              }
            ]
          },
          "moriverUpload": {
            "uploadFilename":"missouri_river_datasheet_20210617_1900_59.csv",
            "Items": [
              {
                "siteId": 0,
                "siteFid": "F-1",
                "mrFid": "20210617-184105056-001",
                "season": "A0",
                "setdate": "2021-06-17T00:00:00Z",
                "subsample": 1,
                "subsamplepass": 1,
                "subsamplen": "R",
                "recorder": "NR",
                "gear": "BSQD",
                "gearType": "E",
                "temp": 1,
                "turbidity": 1,
                "conductivity": 1,
                "do": 1,
                "distance": 1,
                "width": 1,
                "netrivermile": 2.5,
                "structurenumber": "sn",
                "usgs": "usgs",
                "riverstage": 1,
                "discharge": 1,
                "u1": "u1",
                "u2": "u2",
                "u3": "u3",
                "u4": "u4",
                "u5": "u5",
                "u6": "u6",
                "u7": "u7",
                "macro": "BRAD",
                "meso": "CHNB",
                "habitatrn": "N",
                "qc": "q",
                "microStructure": "m",
                "structureFlow": "0",
                "structureMod": "0",
                "setSite_1": "0",
                "setSite_2": "0",
                "setSite_3": "3",
                "startTime": "10:10:00 AM",
                "startLatitude": 36,
                "startLongitude": -88,
                "stopTime": "6:57:10 PM",
                "stopLatitude": 36,
                "stopLongitude": -88,
                "depth1": 1,
                "velocitybot1": 1,
                "velocity08_1": 1,
                "velocity02or06_1": 1,
                "depth2": 1,
                "velocitybot2": 1,
                "velocity08_2": 1,
                "velocity02or06_2": 1,
                "depth3": 1,
                "velocitybot3": 1,
                "velocity08_3": 1,
                "velocity02or06_3": 1,
                "watervel": 1,
                "cobble": 0,
                "organic": 1,
                "silt": 30,
                "sand": 30,
                "gravel": 40,
                "comments": "comments",
                "complete": 1,
                "checkby": "che",
                "noTurbidity": "n",
                "noVelocity": "n"
              }
            ]
          }
        }
        
    - Example `POST` response
        ```
        {
          "uploadSessionId": 5308,
          "uploadedBy": "DeeLiang",
          "siteCntFinal": 1,
          "mrCntFinal": 1,
          "fishCntFinal": 1,
          "searchCntFinal": 1,
          "suppCntFinal": 1,
          "telemetryCntFinal": 1,
          "procedureCntFinal": 1,
          "noSiteCnt": 0,
          "siteMatch": 0,
          "noSiteIDMsg": ""
        }
