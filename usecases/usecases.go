package usecases

import (
	d "github.com/arzonus/cveapi/domain"
)

type Interactor struct {
	CVERepository     d.CVERepository
	ProductRepository d.ProductRepository
	Logger            Logger
}

type NVDData struct {
	CVEs     []d.CVE
	Products []d.Product
}

type Logger interface {
	Log(...interface{})
	Info(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Debug(...interface{})
}

func (inter *Interactor) PopulateDB(data NVDData) {

	maxGos := 25 //Max concurrent connection to MSSQL
	guard := make(chan struct{}, maxGos)

	inter.Logger.Info("Start send CVE Info")
	lenCVEs := len(data.CVEs)
	lenProds := len(data.Products)
	for i, v := range data.CVEs {
		guard <- struct{}{}
		go func(cve d.CVE) {
			inter.SaveCVE(cve)
			<-guard
		}(v)
		inter.Logger.Info("Left CVE: ", lenCVEs-i)
	}
	inter.Logger.Info("Start send Product Info")

	for i, v := range data.Products {
		guard <- struct{}{}
		go func(product d.Product) {
			inter.SaveProduct(product)
			<-guard
		}(v)
		inter.Logger.Info("Left Product: ", lenProds-i)
	}
}

func (inter *Interactor) GetCVEById(cveId string) d.CVE {
	return d.CVE{}
}

func (inter *Interactor) SaveCVE(cve d.CVE) {
	inter.CVERepository.Save(cve)
}

func (inter *Interactor) SaveProduct(product d.Product) {
	inter.ProductRepository.Save(product)
}

func (inter *Interactor) GetListCVEByProduct(product d.Product) []d.CVE {
	return inter.CVERepository.GetListByProduct(product)
}
