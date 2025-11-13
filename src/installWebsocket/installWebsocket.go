package installWebsocket

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kafy11/gosocket/log"
	"github.com/kardianos/service"

	// "gt-gitlab.gtplan.net/websocket/websocket-v2/client/pkg/file"
	// "gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/build"
	// "gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/companyClient"
	// "gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/config"
	// "gt-gitlab.gtplan.net/websocket/websocket-v2/client/src/runFlags"

	"motorv2/src/websocket/pkg/file"
	"motorv2/src/websocket/src/build"
	"motorv2/src/websocket/src/companyClient"
	"motorv2/src/websocket/src/config"
	"motorv2/src/websocket/src/runFlags"
)

const serviceName = ".GTPlanWebsocketSvc"
const serviceDescription = "GTPlan Websocket"

type program struct {
	wsClient *companyClient.Client
	service  service.Service
}

func InstallWebsocket(param string) {
	build.Init()

	if !*runFlags.Flags.Dev {
		err := setWorkingDirectory()
		if err != nil {
			log.Fatal("Erro setando o working directory", err)
		}
	}

	loadDotEnv()

	wsClient, err := companyClient.New(config.GetWSClientParams())
	if err != nil {
		log.Fatal("Falha ao criar o client ", err)
	}

	if *runFlags.Flags.Test {
		wsClient.Test()
		os.Exit(0)
	}

	config.SetWebserviceParams()

	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
		Arguments:   []string{"websocket"},
	}

	prg := &program{
		wsClient: wsClient,
	}

	s, err := service.New(prg, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.service = s

	if param == "install" {
		if !*runFlags.Flags.Dev {
			err = s.Install()
			if err != nil {
				log.Error("Falha ao instalar o serviço", err)
			}
		}
	}

	dirname := "." + string(filepath.Separator)
	log.Info("Deletando arquivos .old da pasta", dirname)
	file.DeleteAllWithExtension(dirname, ".old")

	if param == "run" {
		err = s.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
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

func loadDotEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Erro no .env", err)
	}
}

func (p *program) Run() {
	defer func() {
		if service.Interactive() {
			p.Stop(p.service)
		} else {
			p.service.Stop()
		}
	}()

	err := config.ConnectDB()
	if err != nil {
		log.Error(err)
		return
	}

	//////////////////////////////////////////////////////////////////////////////
	/*
		if !*runFlags.Flags.Dev {
			err := setWorkingDirectory()
			if err != nil {
				log.Fatal("Erro setando o working directory", err)
			}
		}

		loadDotEnv()

		wsClient, err := companyClient.New(config.GetWSClientParams())
		if err != nil {
			log.Fatal("Falha ao criar o client ", err)
		}

		if *runFlags.Flags.Test {
			wsClient.Test()
			os.Exit(0)
		}
	*/
	//////////////////////////////////////////////////////////////////////////////
	p.wsClient.Run()

	return
}

func (p program) Start(s service.Service) error {
	log.Info("Iniciando serviço")
	log.Info("Iniciando o cliente websocket")
	go p.Run()
	return nil
}

func (p program) Stop(s service.Service) error {
	log.Info("Parando serviço")
	return p.wsClient.Stop()
}

func (p program) Restart() error {
	return p.service.Restart()
}
