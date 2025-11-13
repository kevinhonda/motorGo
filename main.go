package main

import (
	"bufio"
	"motorv2/src/diagnosis"
	"motorv2/src/handExec"
	"motorv2/src/installWebsocket"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/svc/mgr"

	//_ "github.com/denisenkom/go-mssqldb" // Importe o driver do SQL Server

	"github.com/joho/godotenv"
	"github.com/kardianos/service"

	"fmt"

	"motorv2/src/conf_env"
	"motorv2/src/mortor"
	"os"

	"github.com/kafy11/gosocket/log"
	src "golang.org/x/sys/windows/svc"
)

type myService struct {
	mode string
}

//go build

func main() {
	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	err := setWorkingDirectory()
	if err != nil {
		log.Fatal("Erro setando o working directory", err)
	}

	// Obtém o diretório de trabalho atual
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório de trabalho:", err)
		os.Exit(1)
	}
	// Printa o diretório de trabalho atual
	log.Info("Diretório de trabalho atual:" + dir)
	fmt.Println("Diretório de trabalho atual:", dir)

	var ambient string
	if isProduction() {
		fmt.Println("Ambiente de Produção detectado")
		ambient = "PRD"
	} else {
		fmt.Println("Ambiente de QA detectado")
		ambient = "QA"
	}
	log.Info("Ambiente: " + ambient)
	conf_env.SetUrl(ambient)

	svcConfig := &service.Config{
		Name:        fmt.Sprintf("MotorV2.%s", ambient),
		DisplayName: fmt.Sprintf(".MotorV2.%s", ambient),
		Description: fmt.Sprintf("MotorV2.%s", ambient),
		Arguments:   []string{"motor"},
	}

	serviceName := fmt.Sprintf("MotorV2.%s", ambient)

	//prg := &myService{}
	prg := &myService{mode: mode}
	svc, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	if service.Interactive() {
		if len(os.Args) > 1 {
			err := service.Control(svc, os.Args[1])
			if err != nil {
				log.Fatal(fmt.Sprintf("erro ao executar comando %s: %v", os.Args[1], err))
			}
			return
		}
	}

	isService, err := src.IsWindowsService()
	if err != nil {
		fmt.Println("Erro ao verificar se o código está sendo executado como um serviço do Windows:", err)
		return
	}

	if isService {
		log.Info("O código está sendo executado como um serviço do Windows.")
	} else {
		log.Info("O código não está sendo executado como um serviço do Windows.")
		exists, err := serviceExists(serviceName)
		if err != nil {
			log.Fatal("Erro ao verificar serviço:", err)
		}

		if exists {
			log.Info("O serviço já existe.")
			fmt.Printf("O serviço já %s existe.\n", serviceName)
			scanner := bufio.NewScanner(os.Stdin)

			for l := 0; l < 1; {
				fmt.Println("Escolha uma das ações")
				fmt.Println("1 - Execução manual")
				fmt.Println("2 - Diagnostico")
				fmt.Println("3 - Instalacao Websocket")
				fmt.Println("4 - Sair")

				scanner.Scan()
				varValue := scanner.Text()
				switch varValue {
				case "1":
					l = 2
					log.Info("Começo execução manual")
					handExec.Exec()
				case "2":
					l = 2
					log.Info("Começo diagnostico")
					diagnosis.MotorDiagnosis()
				case "3":
					l = 2
					log.Info("Começo Instalacao Websocket")
					installWebsocket.InstallWebsocket("install")
				case "4":
					l = 2
					fmt.Println("Pressione Enter para sair...")
					fmt.Scanln()
					os.Exit(0)
				default:
					fmt.Println("////////////////////////////////////")
					fmt.Println("//////////Resposta inválida/////////")
					fmt.Println("////////////////////////////////////")
					l = 0
					//os.Exit(0)
				}
			}

			os.Exit(0)
		} else {
			fmt.Printf("Vamos Configurar? ")
			log.Info("Vamos Configurar? ")
			conf_env.SetConf()
			err = svc.Install()
			if err != nil {
				log.Error(err)
			}
			fmt.Println("Service Criado")
			installWebsocket.InstallWebsocket("install")
			fmt.Println("Pressione Enter para sair...")
			fmt.Scanln()
			os.Exit(0)
		}
	}

	args := os.Args
	log.Info("Args 1")
	log.Info(args[1])
	switch args[1] {
	case "motor":
		log.Info("Iniciando o serviço " + serviceName)
		/**/
		err = svc.Install()
		if err != nil {
			log.Error(err)
		}

		installed, err := serviceExists(serviceName)
		if err != nil {
			log.Fatal("Erro ao verificar serviço:", err)
		}

		if installed {
			fmt.Println("Serviço criado")
		}

		err = svc.Run()
		if err != nil {
			log.Fatal(err)
		}
	case "websocket":
		log.Info("Iniciando o serviço Websocket")
		installWebsocket.InstallWebsocket("run")
	}

}
func isProduction() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}

	exeName := strings.ToLower(filepath.Base(exePath))

	// Verifica se contém indicadores de produção
	return strings.Contains(exeName, "prd") ||
		strings.Contains(exeName, "prod") ||
		strings.Contains(exeName, "production") ||
		strings.Contains(exeName, "producao")
}

func NewMyService(mode string) *myService {
	return &myService{
		mode: strings.ToLower(mode), // Já converte para minúsculo
	}
}

/**/
func serviceExists(serviceName string) (bool, error) {
	m, err := mgr.Connect()
	if err != nil {
		return false, err
	}
	defer m.Disconnect()

	service, err := m.OpenService(serviceName)
	if err != nil {
		return false, nil
	}
	defer service.Close()

	return true, nil
}

func setWorkingDirectory() error {
	fullexecpath, err := os.Executable()
	if err != nil {
		return err
	}

	dir, _ := filepath.Split(fullexecpath)

	os.Chdir(dir)
	return nil
}

func (m *myService) Start(s service.Service) error {
	log.Info("Iniciando o serviço")
	go m.run()
	return nil
}

func (m *myService) Stop(s service.Service) error {
	// Faça o que for necessário para parar o serviço
	log.Info("Parando o serviço")
	return nil
}

func (m *myService) run() {
	log.Info("mode:", m.mode)
	switch strings.ToLower(m.mode) {
	case "motor":
		err := godotenv.Load()
		if err != nil {
			log.Error("Error loading .env file")
		}

		mortor.SetConfig(&mortor.MotorConfigParams{
			WsBaseUrl:  os.Getenv("WS_BASE_URL"),
			WsAuthUser: os.Getenv("WS_AUTH_USER"),
			WsAuthPass: os.Getenv("WS_AUTH_PASS"),
		})
		log.Info("Iniciando o serviço Motor RUNERAS")
		err = mortor.Run()
		if err != nil {
			log.Error(err)
		}
	case "Websocket":
		log.Info("Iniciando o serviço Websocket run")
		installWebsocket.InstallWebsocket("run")

	}
	log.Info("Run")
	return
}
