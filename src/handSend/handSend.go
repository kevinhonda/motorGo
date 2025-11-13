package handSend

import (
	"fmt"

	"motorv2/src/awsControllers"
	"motorv2/src/sendSchedule"

	"bufio"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/kafy11/gosocket/log"
)

func HandSend() {
	var macID string

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Erro ao buscar interfaces de rede:", err)
		return
	}

	// Loop pelas interfaces para encontrar o endereço MAC
	for _, iface := range interfaces {
		// Ignora interfaces loopback e interfaces sem endereço MAC
		if iface.Flags&net.FlagLoopback == 0 && len(iface.HardwareAddr) > 0 {
			macAddr := iface.HardwareAddr
			fmt.Println("Endereço MAC:", macAddr.String())
			log.Info("Endereço MAC:" + macAddr.String())
			macID = macAddr.String()
		}
	}

	//Pergunta com resposta aberta
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Deseja executar qual ação?")
	fmt.Println("1 - Envio de todos os layouts")

	fmt.Println("2 - Envio de layouts especificos")
	scanner.Scan()
	sendValue := scanner.Text()
	switch sendValue {
	case "1":
		log.Info("Envio all - Inicio")
		fmt.Println("Envio all - Inicio")

		sendSchedule.BusConfig(macID, "")
		fmt.Println("Envio all - Final")
		log.Info("Envio all - Final")
	case "2":
		bodyResp, _ := awsControllers.GetInfos("Adm_get_wes_env?id=" + macID)
		//layoutsSend := bodyResp["LAYOUTS"].(string)
		if layoutsSendStr, exists := bodyResp["LAYOUTS"].(string); exists {
			log.Info("Layouts para envio: " + layoutsSendStr)
			fmt.Println("Layouts para envio: " + layoutsSendStr)
			layoutStr := strings.ReplaceAll(strings.ToUpper(layoutsSendStr), " ", "")
			log.Info("Layouts: " + layoutStr)
			layouts := strings.Split(layoutStr, ",")
			var layoutsSendName string

			for w := 0; w < 1; {
				for index, layout := range layouts {
					fmt.Printf("%d - %s\n", index, layout)
				}
				fmt.Println("Quais layouts deseja enviar?")
				scanner.Scan()
				layoutsValues := scanner.Text()
				layoutsValues = strings.ReplaceAll(strings.ToUpper(layoutsValues), " ", "")
				layoutsVavalue := strings.Split(layoutsValues, ",")
				layoutsSendName = ""
				w = len(layoutsVavalue)
				for _, layoutIndex := range layoutsVavalue {
					// Convert string to integer
					index, err := strconv.Atoi(layoutIndex)
					if err != nil {
						fmt.Println("///////////////////////////////////////")
						fmt.Printf("Erro: '%s' não é um número válido\n", layoutIndex)
						fmt.Println("///////////////////////////////////////")
						w = w - 1
						continue
					}

					// Check if index is within bounds
					if index < 0 || index >= len(layouts) {
						fmt.Println("///////////////////////////////////////")
						fmt.Printf("Erro: índice %d fora dos limites (0-%d)\n", index, len(layouts)-1)
						fmt.Println("///////////////////////////////////////")
						w = w - 1
						continue
					}

					layoutsSendName += layouts[index]
					layoutsSendName += ", "
				}

				if len(layoutsSendName) > 0 {
					layoutsSendName = layoutsSendName[:len(layoutsSendName)-2]
				}

				fmt.Println("Layouts para envio: " + layoutsSendName)
				fmt.Println("Deseja enviar esse layous " + layoutsSendName)
				if w > 0 {
					for iw := 0; iw < 1; {
						scanner.Scan()
						varValue := strings.ToLower(scanner.Text())

						switch varValue {
						case "sim", "s":
							w = 1
							iw = 2
						case "não", "n", "nao":
							w = -1
							iw = 2
						default:
							//fmt.Println("Resposta inválida. Por favor, responda com 'sim', 's', 'não' ou 'n'.")
							fmt.Println("////////////////////////////////////")
							fmt.Println("//////////Resposta inválida/////////")
							fmt.Println("////////////////////////////////////")
						}
					}
				}
			}
			log.Info("Envio parts - Inicio")
			fmt.Println("Envio parts - Inicio")
			//sendLayouts(layoutsSendName)
			sendSchedule.BusConfig(macID, layoutsSendName)
		} else {
			log.Error("Nenhum layout configurado para envio.")
			fmt.Println("Nenhum layout configurado para envio.")
		}

	default:
		fmt.Println("Resposta invalida")
		//os.Exit(0)
		fmt.Println("Pressione Enter para sair...")
		fmt.Scanln()
		os.Exit(0)
	}
	fmt.Println("Pressione Enter para sair...")
	fmt.Scanln()
	os.Exit(0)
}
