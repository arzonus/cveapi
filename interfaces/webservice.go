package interfaces

import (
	"encoding/json"
	d "github.com/arzonus/cveapi/domain"
	"github.com/gorilla/mux"
	"net/http"
)

type WebserviceHandler struct {
	Interactor Interactor
	Logger     Logger
}

func (handler WebserviceHandler) GetCVEById(res http.ResponseWriter, req *http.Request) {
	handler.Logger.Info("GetCVEById")
	cveId := getParamFromURL("cveId", req)
	cve := handler.Interactor.GetCVEById(cveId)
	res.Write([]byte(cve.Id))
}

func (handler WebserviceHandler) GetListCVEByProduct(res http.ResponseWriter, req *http.Request) {
	handler.Logger.Info("Get List CVE By Product")
	handler.Logger.Info(req.URL)

	name := getParamFromURL("name", req)
	version := getParamFromURL("version", req)
	product := d.Product{
		Name:    name,
		Version: version,
	}
	cves := handler.Interactor.GetListCVEByProduct(product)

	type jCVE struct {
		CVEs []d.CVE
	}

	data, err := json.Marshal(jCVE{CVEs: cves})
	handler.Logger.Info(data)

	if err != nil {
		handler.Logger.Info(err)
	}

	res.Write(data)

}

func getParamFromURL(param string, req *http.Request) string {
	return mux.Vars(req)[param]
}
