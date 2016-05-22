//
// CVE API (CVA)
//
// Author: Vsevolod Kaloshin
//
package main

import (
	"github.com/arzonus/cveapi/infrastructure"
	"github.com/arzonus/cveapi/interfaces"
	"github.com/arzonus/cveapi/usecases"
	"github.com/codegangsta/cli"
	"os"
)

const (
	DBTYPE    = "mssql"
	NVDCVEURL = "http://static.nvd.nist.gov/feeds/xml/cve/nvdcve-2.0-"
	NVDCPEURL = "http://static.nvd.nist.gov/feeds/xml/cpe/dictionary/official-cpe-dictionary_v2.3.xml.zip"
)

var (
	MSSQLURI  string
	TEMPPATH  string
	STARTYEAR string
	ENDYEAR   string
)

func main() {

	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "CVE API"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			EnvVar:      "CVA_DB_MSSQL",
			Name:        "db-mssql",
			Value:       "server=lab-eu-lastbackend.database.windows.net;user id=lastbackend;password=fw9UJEfbpFijdlA2YDk8;database=db-00-lab-eu;sslmode=disable",
			Usage:       "Connection string for connect to MSSQL DB",
			Destination: &MSSQLURI,
		},
		cli.StringFlag{
			EnvVar:      "CVA_TEMP_PATH",
			Name:        "temp-path",
			Value:       "/tmp/",
			Usage:       "Set temporaty folder",
			Destination: &TEMPPATH,
		},
		cli.StringFlag{
			EnvVar:      "CVA_START_YEAR",
			Name:        "start-year",
			Value:       "2003",
			Usage:       "Set start year",
			Destination: &STARTYEAR,
		},
		cli.StringFlag{
			EnvVar:      "CVA_END_YEAR",
			Name:        "end-year",
			Value:       "2017",
			Usage:       "Set end year",
			Destination: &ENDYEAR,
		},
	}

	app.Action = Initialize
	app.Run(os.Args)

	select {}
}

func Initialize(c *cli.Context) {

	dbHandler := infrastructure.NewMSSQLHandler(MSSQLURI)
	dbHandler.Logger = infrastructure.Logger{}

	handlers := make(map[string]interfaces.DBHandler)
	handlers["DBCVERepo"] = dbHandler
	handlers["DBProductRepo"] = dbHandler

	Interactor := new(usecases.Interactor)

	CVERepository := interfaces.NewDBCVERepo(handlers)
	CVERepository.Logger = infrastructure.Logger{}

	ProductRepository := interfaces.NewDBProductRepo(handlers)
	ProductRepository.Logger = infrastructure.Logger{}

	Interactor.CVERepository = CVERepository
	Interactor.ProductRepository = ProductRepository
	Interactor.Logger = infrastructure.Logger{}

	webserviceHandler := interfaces.WebserviceHandler{}
	webserviceHandler.Logger = infrastructure.Logger{}
	webserviceHandler.Interactor = Interactor

	extdataHandler := interfaces.ExtdataHandler{
		StartYear: STARTYEAR,
		EndYear:   ENDYEAR,
		NVDCVEURL: NVDCVEURL,
		NVDCPEURL: NVDCPEURL,
		Temp:      TEMPPATH,
	}
	extdataHandler.Logger = infrastructure.Logger{}
	extdataHandler.Interactor = Interactor

	go HTTP(webserviceHandler)
	go Update(extdataHandler)

	select {}

}
