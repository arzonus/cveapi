package interfaces

import (
	"fmt"
	d "github.com/arzonus/cveapi/domain"
)

type DBHandler interface {
	Execute(statement string)
	Query(statement string) Row
}

type Row interface {
	Scan(dest ...interface{})
	Next() bool
}

type DBRepo struct {
	dbHandlers map[string]DBHandler
	dbHandler  DBHandler
	Logger     Logger
}

type DBCVERepo DBRepo

func NewDBCVERepo(dbHandlers map[string]DBHandler) *DBCVERepo {
	dbCVERepo := new(DBCVERepo)
	dbCVERepo.dbHandlers = dbHandlers
	dbCVERepo.dbHandler = dbHandlers["DBCVERepo"]
	return dbCVERepo
}

func (repo *DBCVERepo) Save(cve d.CVE) {
	query := fmt.Sprintf("INSERT INTO cves (cve_id, summary) VALUES ('%v', '%q');", cve.Id, cve.Summary)
	repo.dbHandler.Execute(query)

	lenCPEs := len(cve.CPEs)

	if lenCPEs > 0 {

		query = fmt.Sprintf("INSERT INTO cpes_cves (cve_id, cpe_id) VALUES")

		for i, v := range cve.CPEs {

			if i < len(cve.CPEs)-1 {
				query = query + fmt.Sprintf(" ('%v', '%v'),", cve.Id, v)
			} else {
				query = query + fmt.Sprintf(" ('%v', '%v');", cve.Id, v)
			}
		}

		repo.dbHandler.Execute(query)
	}
}

func (repo *DBCVERepo) GetById(cveId string) d.CVE {
	row := repo.dbHandler.Query(fmt.Sprintf("SELECT cve_id, summary FROM cves WHERE cve_id = '%v'", cveId))
	var (
		dcveId   string
		dsummary string
	)
	row.Next()
	row.Scan(&dcveId, &dsummary)
	repo.Logger.Info("dcveId: ", dcveId)
	repo.Logger.Info("dsummary: ", dsummary)
	return d.CVE{Id: dcveId, Summary: dsummary}
}

func (repo *DBCVERepo) GetList() []d.CVE {
	return []d.CVE{}
}

func (repo *DBCVERepo) GetListByProduct(product d.Product) []d.CVE {
	query := fmt.Sprintf("SELECT c.cve_id, c.summary FROM cves AS c RIGHT JOIN (SELECT cve_id FROM cpes_cves AS cc LEFT JOIN products AS p ON cc.cpe_id=p.cpe_id WHERE p.name='%v' AND p.product_version='%v') AS b ON b.cve_id=c.cve_id",
		product.Name, product.Version)

	repo.Logger.Info(query)
	row := repo.dbHandler.Query(query)

	var cves []d.CVE

	for row.Next() {
		var cve d.CVE
		row.Scan(&cve.Id, &cve.Summary)
		cves = append(cves, cve)
	}

	return cves
}

type DBProductRepo DBRepo

func NewDBProductRepo(dbHandlers map[string]DBHandler) *DBProductRepo {
	dbProductRepo := new(DBProductRepo)
	dbProductRepo.dbHandlers = dbHandlers
	dbProductRepo.dbHandler = dbHandlers["DBProductRepo"]
	return dbProductRepo
}

func (repo *DBProductRepo) Save(product d.Product) {
	query := fmt.Sprintf("INSERT INTO products (name, vendor, product_version, type, cpe_id, title) VALUES ('%v','%v','%v','%v','%v','%q');",
		product.Name,
		product.Vendor,
		product.Version,
		product.Type,
		product.CPE,
		product.Title)

	repo.dbHandler.Execute(query)

}

func (repo *DBProductRepo) GetById(cpeId string) d.Product {
	return d.Product{}
}

func (repo *DBProductRepo) GetList() []d.Product {
	return []d.Product{}
}
