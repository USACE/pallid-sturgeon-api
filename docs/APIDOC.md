# Instrumentation API

Table of Contents
- [Pallid Sturgeon API](#pallidsturgeon-api)
    - [seasons](#seasons)
    - [uploads](#uploads)

---
### Seasons
- List Seasons \
  [http://localhost:8080/psapi/seasons](http://localhost:8080/psapi/seasons)

---
### Uploads
- Get Upload Session Id \
  [http://localhost:8080/psapi/uploadSessionId](http://localhost:8080/psapi/uploadSessionId)
- Site Upload \
  `http://localhost:8080/psapi/siteUpload`
    - Example `POST` body
        ```
        [
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
            "comments": "test",
            "uploadSessionId": 5308,
            "uploadFilename": "pspa_sites_datasheet_20210617_1900_59.csv"
          }
       ]
- Fish Upload \
  `http://localhost:8080/psapi/fishUpload`
    - Example `POST` body
        ```
        [
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
            "comments": "test",
            "uploadSessionId": 5308,
            "uploadFilename": "fish_datasheet_20210617_1900_59.csv"
          }
       ]
- Search Upload \
  `http://localhost:8080/psapi/searchUpload`
    - Example `POST` body
        ```
        [
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
            "conductivity": 22,
            "uploadSessionId": 5308,
            "uploadFilename": "search_effort_20210617_1900_59.csv"
          }
       ]
- Telemetry Upload \
  `http://localhost:8080/psapi/telemetryUpload`
    - Example `POST` body
        ```
        [
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
            "comments": "comments",
            "uploadSessionId": 5308,
            "uploadFilename": "telemetry_20210617_1900_59.csv"
          }
       ]
- Procedure Upload \
  `http://localhost:8080/psapi/procedureUpload`
    - Example `POST` body
        ```
        [
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
            "gonadCondition": "1",
            "uploadSessionId": 5308,
            "uploadFilename": "procedure_20210617_1900_59.csv"
          }
       ]
- Supplemental Upload \
  `http://localhost:8080/psapi/supplementalUpload`
    - Example `POST` body
        ```
        [
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
            "comments": "comments",
            "uploadSessionId": 5308,
            "uploadFilename": "supplemental_20210617_1900_59.csv"
          }
       ]
- Moriver Upload \
  `http://localhost:8080/psapi/moriverUpload`
    - Example `POST` body
        ```
        [
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
            "noVelocity": "n",
            "uploadSessionId": 5308,
            "uploadFilename":"missouri_river_datasheet_20210617_1900_59.csv"
          }
       ]
 - Store Procedure \
  `http://localhost:8080/psapi/storeProcedure/1606`
    - Example `GET` response
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