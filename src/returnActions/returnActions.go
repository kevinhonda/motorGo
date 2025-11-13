package returnActions

import (
	"motorv2/src/returnLayout3wm"
	"motorv2/src/returnLayoutDrp"
	"motorv2/src/returnLayoutOc"
	"motorv2/src/returnLayoutSug"

	//"motorv2/src/returnFuncs"
	"github.com/kafy11/gosocket/log"
)

var layoutReturn string

type ReturnTest struct{}

// Definir o tipo OrdemHeader como uma string

func (ra ReturnTest) ReturnStuffs(layouts string, querys map[string]interface{}) {
	//func ReturnStuffs(layouts string) {
	//body, _ := getReturn("Po_get_id?status=2")
	log.Info("ReturnStuffs " + layouts + " querys:")
	switch layouts {
	case "Mrp_po":
		returnLayoutOc.ReturnStuffs("Mrp_po", querys)
	case "Sug_id":
		returnLayoutSug.ReturnStuffs("Sug_id", querys)
	case "Drp_id":
		returnLayoutDrp.ReturnStuffs("Drp_id", querys)
	case "3WM":
		returnLayout3wm.ReturnStuffs("3WM", querys)
	}

}
