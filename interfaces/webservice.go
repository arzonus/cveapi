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
	name := getParamFromURL("name", req)
	version := getParamFromURL("version", req)
	product := d.Product{
		Name:    name,
		Version: version,
	}
	cves := handler.Interactor.GetListCVEByProduct(product)
	data, err := json.Marshal(cves)
	if err != nil {
		handler.Logger.Info(err)
	}

	res.Write(data)

}

func getParamFromURL(param string, req *http.Request) string {
	return mux.Vars(req)[param]
}
