/*
=============================================
Author:				Michael Kaiser
Create date:		2025-08-20
Description:		Script to setup the needed database schemas and other structures for the insight DV.
Modification:
YYYY-MM-DD	mkaiser
=============================================
*/

DROP SCHEMA IF EXISTS cfg;
DROP SCHEMA IF EXISTS inb;
DROP SCHEMA IF EXISTS rdv;
DROP SCHEMA IF EXISTS udm;

CREATE SCHEMA cfg;
CREATE SCHEMA inb;
CREATE SCHEMA rdv;
CREATE SCHEMA udm;