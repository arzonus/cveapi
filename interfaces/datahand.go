package interfaces

import (
	"archive/zip"
	"encoding/xml"
	d "github.com/arzonus/cveapi/domain"
	u "github.com/arzonus/cveapi/usecases"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ExtdataHandler struct {
	Interactor Interactor
	Logger     Logger
	StartYear  string
	EndYear    string
	Temp       string
	NVDCVEURL  string
	NVDCPEURL  string
	NVDData    NVDData
}

type NVDData struct {
	CVEs []CVEData
	CPEs []CPEData
}

type CVEs struct {
	Entries []CVEData `xml:"entry"`
}

type CVEData struct {
	Id      string   `xml:"cve-id"`
	Summary string   `xml:"summary"`
	CPEs    []string `xml:"vulnerable-software-list>product"`
}

type CPEs struct {
	Entries []CPEData `xml:"cpe-item"`
}

type CPEData struct {
	Id    string `xml:"name,attr"`
	Title string `xml:"title"`
}

func (handler *ExtdataHandler) PopulateDB() {
	handler.Logger.Info("Populate Database")
	handler.getNVDCPE()
	//handler.getNVDCVE()
	handler.Interactor.PopulateDB(handler.convNVDIntToNVDUse())
	handler.Logger.Info("There was stoped")
}

func (handler *ExtdataHandler) convNVDIntToNVDUse() u.NVDData {

	handler.Logger.Info("Start converting CVE")
	var cves []d.CVE
	var products []d.Product

	for _, v := range handler.NVDData.CVEs {
		cve := d.CVE{
			Id:      v.Id,
			Summary: strings.Replace(v.Summary, "'", " ", -1),
			CPEs:    v.CPEs,
		}
		cves = append(cves, cve)
	}
	handler.Logger.Info("End converting CVE")

	handler.Logger.Info("Start converting CPE")

	for _, v := range handler.NVDData.CPEs {
		product := parseCPE(v.Id)
		product.Title = strings.Replace(v.Title, "'", " ", -1)
		//product.CVEs = findCVEs(product.CPE, cves)
		products = append(products, product)
	}

	handler.Logger.Info("End convering CPE")

	return u.NVDData{
		CVEs:     cves,
		Products: products,
	}
}

func (handler *ExtdataHandler) getNVDCPE() {

	url := handler.NVDCPEURL
	filenameZIP := "official-cpe-dictionary_v2.3.zip"
	filenameXML := "official-cpe-dictionary_v2.3.xml"

	err := handler.downloadData(url, filenameZIP)

	if err != nil {
		handler.Logger.Error(err)
	}

	err = handler.unzipData(filenameZIP)

	if err != nil {
		handler.Logger.Error(err)
	}

	cpes, err := handler.parseNVDCPEData(filenameXML)

	if err != nil {
		handler.Logger.Error(err)
	}

	handler.NVDData.CPEs = cpes.Entries
}

func findCVEs(cpe string, ucves []d.CVE) []string {
	var cves []string
	for _, cve := range ucves {
		for _, v := range cve.CPEs {
			if cpe == v {
				cves = append(cves, cve.Id)
				break
			}
		}
	}
	return cves
}

func (handler *ExtdataHandler) getNVDCVE() {

	sy, _ := strconv.ParseInt(handler.StartYear, 10, 64)
	ey, _ := strconv.ParseInt(handler.EndYear, 10, 64)

	var cves []CVEData

	countGo := ey - sy
	done := make(chan struct{}, countGo)

	for i := sy; i < ey; i++ {

		go func(i int64) {
			year := strconv.FormatInt(i, 10)

			url := handler.NVDCVEURL + year + ".xml.zip"
			filenameZIP := "nvd-" + year + ".zip"
			filenameXML := "nvdcve-2.0-" + year + ".xml"

			err := handler.downloadData(url, filenameZIP)

			if err != nil {
				handler.Logger.Error(err)
			}

			err = handler.unzipData(filenameZIP)

			if err != nil {
				handler.Logger.Error(err)
			}

			cve, err := handler.parseNVDCVEData(filenameXML)

			if err != nil {
				handler.Logger.Error(err)
			}

			cves = append(cves, cve.Entries...)
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < int(countGo); i++ {
		<-done
	}

	handler.NVDData.CVEs = cves
}

func (handler *ExtdataHandler) parseNVDCVEData(filename string) (CVEs, error) {

	var cves CVEs

	path := handler.Temp + filename
	handler.Logger.Info("Read file", path)

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return cves, err
	}
	handler.Logger.Info("Parse file", path)

	err = xml.Unmarshal(data, &cves)
	handler.Logger.Info("End parse", path)

	if err != nil {
		return cves, err
	}
	return cves, nil
}

func (handler *ExtdataHandler) parseNVDCPEData(filename string) (CPEs, error) {

	var cpes CPEs

	path := handler.Temp + filename
	handler.Logger.Info("Read file", path)

	data, err := ioutil.ReadFile(path)

	if err != nil {
		return cpes, err
	}

	handler.Logger.Info("Parse file", path)

	err = xml.Unmarshal(data, &cpes)
	handler.Logger.Info("End parse", path)

	if err != nil {
		return cpes, err
	}

	return cpes, nil
}

func (handler *ExtdataHandler) downloadData(url string, filename string) error {

	path := handler.Temp + filename
	handler.Logger.Info("Create file", path)

	out, err := os.Create(path)
	defer out.Close()

	if err != nil {
		return err
	}

	handler.Logger.Info("Start download data")
	handler.Logger.Info("URL", url)
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	handler.Logger.Info("End download data")

	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return err
	}

	return nil
}

func (handler *ExtdataHandler) unzipData(archive string) error {
	handler.Logger.Info("Start unzip data")
	archive = handler.Temp + archive
	target := handler.Temp
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		handler.Logger.Info(path)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		fileReader.Close()

		out, err := os.Create(path)
		defer out.Close()

		if _, err := io.Copy(out, fileReader); err != nil {
			return err
		}
	}
	handler.Logger.Info("End unzip data")
	return nil
}

func parseCPE(cpe string) d.Product {

	str := strings.Split(cpe, ":")

	return d.Product{
		Version: str[4],
		Name:    str[3],
		Vendor:  str[2],
		Type:    str[1][1:],
		CPE:     cpe,
	}
}
