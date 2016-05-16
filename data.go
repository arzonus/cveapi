package main

import (
	"github.com/arzonus/cveapi/interfaces"
)

func Update(edh interfaces.ExtdataHandler) {
	edh.PopulateDB()
}
