package returnEng

import (
	"fmt"
	"motorv2/pkg/db"
	"motorv2/src/returnActions"
	"motorv2/src/returnFuncs"

	"motorv2/src/sendStuffs"

	//"motorv2/pkg/ws"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kafy11/gosocket/log"
)

var dbs *sqlx.DB

func ReturnAction() {
	//ws.SetAuth("11C16AF52264426C02A28813E4D4F1B5", "C51CE410C124A10E0DB5E4B97FC2AF39")

	returnFuncsInstance := returnFuncs.GetTest{}
	bodyResp, _ := returnFuncsInstance.GetReturn("Adm_get_wes_config")
	log.Info("ReturnAction - Começo")
	var returnQuerys map[string]interface{}

	//log.Info("ReturnAction - dataQuerys")
	envs := bodyResp["env"].(map[string]interface{})
	layoutsMini := envs["LAYOUTS_MINI"].(string)
	layoutsMini = strings.ReplaceAll(strings.ToUpper(layoutsMini), " ", "")
	layoutsMiniList := strings.Split(layoutsMini, ",")

	dataQuerys := bodyResp["query"].(map[string]interface{})
	//Setando querys de retorno
	if returnQuery, exists := dataQuerys["RETURN"].(map[string]interface{}); exists {
		//returnQuerys = returnQuery["RETURN"].(map[string]interface{})
		returnQuerys = returnQuery

		log.Info("returnQuerys:")
		log.Info(returnQuerys)
	} else {
		log.Error("Sem querys de RETORNO")
	}

	miniQuerys := make(map[string]string)
	if sendQuery, exists := dataQuerys["NEW"].(map[string]interface{}); exists {
		for _, layoutMini := range layoutsMiniList {
			if sendQueryr, exists := sendQuery[layoutMini].(string); exists && sendQueryr != "" {
				//log.Info("Query encontrada para layout:", layoutMini)
				miniQuerys[layoutMini] = sendQueryr
			} else {
				log.Error("Layout sem query válida:", layoutMini)
			}
		}
	}

	log.Info(fmt.Sprintf("Layouts Mini: %s ", layoutsMini))
	/*
		log.Info("Querys Layout Mini")
		log.Info(miniQuerys)
		for _, layoutMini := range layoutsMiniList {
			log.Info(fmt.Sprintf("Layouts mini %s: %s", layoutMini, miniQuerys[layoutMini]))
		}
	*/
	if returnQuerys != nil || len(miniQuerys) > 0 {
		var err error
		// Cria o ticker UMA vez fora do loop
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop() // Garante que os recursos serão liberados

		for {
			<-ticker.C // Espera 5 minutos

			if dbs, err = db.OpenDB(); err != nil {
				log.Error("Error initializing database connection:", err)
				continue
			}

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				if returnQuerys != nil {
					if returnLayouts, exists := bodyResp["env"].(map[string]interface{})["RETURN"].(string); exists && returnLayouts != "" {
						log.Info("A chave RETURN existe")
						RunReturns(returnLayouts, returnQuerys)
					} else {
						log.Error("Não possui query para a company ou para o layout.")
					}
				}
			}()

			go func() {
				defer wg.Done()
				if len(miniQuerys) > 0 {
					RunMinis(miniQuerys)
				}
			}()
			log.Info("Antes de wg.Wait()")
			wg.Wait()
			log.Info("Depois wg.Wait()")
			if err := dbs.Close(); err != nil {
				log.Error("Erro ao fechar a conexão com o banco de dados:", err)
			} else {
				log.Info("Conexão com o banco de dados fechada com sucesso")
			}

			if dbs.Stats().OpenConnections > 0 {
				log.Error("Ainda há conexões abertas após Close()")
			} else {
				log.Info("Todas as conexões foram fechadas com sucesso")
			}
			log.Info("Finalzin do loop")
		}
	} else {
		log.Error("Não possui query de RETURN para a company.")
	}

}

func RunMinis(sendQuerys map[string]string) {
	var wg sync.WaitGroup

	for layout, sendQuery := range sendQuerys {
		wg.Add(1)
		go func(layout, query string) {
			defer wg.Done()
			//log.Info("Mini layouts: ", layout)
			//log.Info("Mini sendQuery: ", sendQuery)
			sendStuffs.BusSend(layout, sendQuery, "y")
		}(layout, sendQuery)
	}
	wg.Wait()
}

func RunReturns(returnLayouts string, dataQuery map[string]interface{}) {
	log.Info("RunReturns - Começo da exe")
	layouts := strings.Split(returnLayouts, ",")
	var wg sync.WaitGroup
	returnActionsInstance := returnActions.ReturnTest{}

	for i := 0; i < len(layouts); i++ {
		log.Info("RunReturns - Rodando: ", layouts[i])
		wg.Add(1)
		go func(layout string) {
			defer wg.Done()
			layout = strings.TrimSpace(layout)
			if data, exists := dataQuery[strings.ToUpper(layout)].(map[string]interface{}); exists {
				returnActionsInstance.ReturnStuffs(layout, data)
				log.Info(data)
				log.Info(fmt.Sprintf("RunReturns: Executando retorno do layout: %s ", layout))
			} else {
				log.Error(dataQuery)
				log.Error(layout + " - Não possui query de para a company ou para o layout.")
			}
		}(layouts[i])
	}
	wg.Wait()

	returnFuncs.DeleteFolder()
}
