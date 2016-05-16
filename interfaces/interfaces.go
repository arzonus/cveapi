package interfaces

import (
	d "github.com/arzonus/cveapi/domain"
	u "github.com/arzonus/cveapi/usecases"
)

type Interactor interface {
	GetCVEById(CVEId string) d.CVE
	SaveCVE(cve d.CVE)
	PopulateDB(data u.NVDData)
	GetListCVEByProduct(product d.Product) []d.CVE
}

type Logger interface {
	Log(...interface{})
	Info(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Debug(...interface{})
}
