# Instrumentation API

Table of Contents
- [Pallid Sturgeon API](#pallidsturgeon-api)
    - [seasons](#seasons)
    - [uploads](#uploads)

---
### Uploads
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
            "fieldOffice": "SD - South Dakota Game Fish & Parks",
            "projectId": 1,
            "project": "1 - Pallid Sturgeon Population Assessment",
            "segmentId": 0,
            "segment": "7 - Gavins Point Dam to Ponca",
            "seasonId": "FC",
            "season": "FC - Fish Community",
            "bend": 10,
            "bendrn": "R",
            "bendRiverMile": 799.5,
            "comments": "test",
            "uploadSessionId": 78,
            "uploadFilename": "test_datasheet.csv"
          },
          {
            "siteId": 1,
            "siteFid": "F-1",
            "siteYear": 2021,
            "fieldofficeID": "SD",
            "fieldOffice": "SD - South Dakota Game Fish & Parks",
            "projectId": 1,
            "project": "1 - Pallid Sturgeon Population Assessment",
            "segmentId": 0,
            "segment": "7 - Gavins Point Dam to Ponca",
            "seasonId": "FC",
            "season": "FC - Fish Community",
            "bend": 10,
            "bendrn": "R",
            "bendRiverMile": 799.5,
            "comments": "test",
            "uploadSessionId": 78,
            "uploadFilename": "test_datasheet.csv"
          }
       ]
