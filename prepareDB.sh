SERVER=dev-eu-lastbackend.database.windows.net
USER=lastbackend
PASS=xlo0yZTeBOtWbDGidGx4
DB=db-03-dev-eu

echo "DROP DATABASES"
tsql -S $SERVER -U $USER -P $PASS -D $DB <<EOF
USE [db-03-dev-eu]
GO
DECLARE @name VARCHAR(128)
DECLARE @SQL VARCHAR(254)
 
SELECT @name = (SELECT TOP 1 [name]
                FROM sysobjects
                WHERE [type] = 'U' AND category = 0
                ORDER BY [name])
WHILE @name IS NOT NULL
  BEGIN
    SELECT @SQL = 'DROP TABLE [dbo].[' + RTRIM(@name) + ']'
    EXEC (@SQL)
    PRINT 'Dropped Table: ' + @name
    SELECT @name = (SELECT TOP 1 [name]
                    FROM sysobjects
                    WHERE [type] = 'U' AND category = 0 AND [name] > @name
                    ORDER BY [name])
  END
GO
GO
CREATE TABLE cves (
    cve_id VARCHAR(256),
    summary VARCHAR(MAX)
);
GO
GO
CREATE TABLE cpes_cves (
    cve_id VARCHAR(256),
    cpe_id VARCHAR(256)
);
GO
GO
CREATE TABLE products (
    name VARCHAR(256),
    vendor VARCHAR(256),
    product_version VARCHAR(256),
    type VARCHAR(256),
    cpe_id VARCHAR(256),
    title VARCHAR(256)
);
GO