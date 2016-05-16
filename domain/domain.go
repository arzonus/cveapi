package domain

type CVERepository interface {
	Save(cve CVE)
	GetById(cveId string) CVE
	GetListByProduct(product Product) []CVE
	GetList() []CVE
}

type ProductRepository interface {
	Save(product Product)
	GetById(cpeId string) Product
	GetList() []Product
}

type CVE struct {
	Id      string
	Summary string
	CPEs    []string
}

type Product struct {
	Name    string
	Vendor  string
	Version string
	Type    string
	CPE     string
	Title   string
	CVEs    []string
}
