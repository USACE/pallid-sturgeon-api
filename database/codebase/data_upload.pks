create or replace PACKAGE          data_upload AS
/******************************************************************************
   NAME:       MEDIA
   PURPOSE:  upload data file

   REVISIONS:
   Ver        Date        Author           Description
   ---------  ----------  ---------------  ------------------------------------
   1.0        1/20/2012   J. Foster        1. Created this package.
   1.1         25FEB2020 JDK         1. Added uploadDSSearchDatasheet
   1.2        03APR2020 JDK         1. Work
******************************************************************************/

PROCEDURE uploadDSSearchDatasheet(p_user IN upload_supplemental.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_sfidMatch OUT number,
                                p_fileBrowseSup IN varchar2 default null
                                );

  PROCEDURE uploadMRdatasheet (
                                p_user IN upload_mr.uploaded_by%TYPE,
                                p_complete IN upload_mr.complete%TYPE,
                                p_checkby IN upload_mr.checkby%TYPE,
                                p_cnt OUT number,
                                p_mrfidMatch OUT number,
                                p_fileBrowseMR IN varchar2 default null
                            );
                            
  PROCEDURE uploadFishDatasheet (
                                p_user IN upload_fish.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_ffidMatch OUT number,
                                p_fileBrowseF IN varchar2 default null
                            );
                            
  PROCEDURE uploadSuppDatasheet (
                                p_user IN upload_supplemental.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_sfidMatch OUT number,
                                p_fileBrowseSup IN varchar2 default null
                            );
                            
  PROCEDURE uploadSiteDatasheet (
                                p_user IN upload_sites.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_siteMatch OUT number,
                                p_fileBrowseSite IN varchar2 default null
                            );
                            
PROCEDURE uploadSiteDatasheetCheck (
                                p_user IN upload_sites.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_siteMatch OUT number,
                                p_fileBrowseSite IN varchar2 default null,
                                p_siteSessionID OUT number
                            );
PROCEDURE uploadSiteDatasheetCheck2020 (
                                p_user IN upload_sites.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_siteMatch OUT number,
                                p_fileBrowseSite IN varchar2 default null,
                                p_siteSessionID OUT number
                            );                            
                            
                            
PROCEDURE uploadMRdatasheetCheck (
                                p_user IN upload_mr.uploaded_by%TYPE,
                                -- p_complete IN upload_mr.complete%TYPE,
                                p_checkby IN upload_mr.checkby%TYPE,
                                p_cnt OUT number,
                                p_mrfidMatch OUT number,
                                p_fileBrowseMR IN varchar2 default null,                                
                                p_mrSessionID OUT number
                            );
                            
PROCEDURE uploadFishDatasheetCheck (
                                p_user IN upload_fish.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_ffidMatch OUT number,
                                p_fileBrowseF IN varchar2 default null,
                                p_fishSessionID OUT number,
                                p_mrSessionID in number
                            );
                            
PROCEDURE uploadSuppDatasheetCheck (
                                p_user IN upload_supplemental.uploaded_by%TYPE,
                                p_cnt OUT number,
                                p_sfidMatch OUT number,
                                p_fileBrowseSup IN varchar2 default null,
                                p_suppSessionID OUT number,
                                p_mrSessionID in number,
                                p_plus_cnt OUT number,
                                p_plus_FID OUT varchar2
                            );
                            
PROCEDURE uploadFinal (
                                p_user IN ds_sites_check.uploaded_by%TYPE,
                                p_site_cnt_final OUT number,
                                p_siteSessionID in number,
                                p_mrSessionID in number,
                                p_mr_cnt_final OUT number,
                                p_fishSessionID in number,
                                p_fishCntFinal OUT number,
                                p_suppSessionID in number,
                                p_suppCntFinal OUT number,
                                p_plus_cnt in number,
                                p_noSite_cnt OUT number,
                                p_noSiteID_msg OUT varchar2
                            );

END data_upload;