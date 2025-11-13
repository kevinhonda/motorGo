package handExec

import (
	"fmt"
	"motorv2/pkg/ws"

	"bufio"
	"motorv2/src/handReturn"
	"motorv2/src/handSend"
	"motorv2/src/update"
	"os"

	//"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kafy11/gosocket/log"
)

//var dbs *sqlx.DB

func Exec() {

	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	//Pergunta com resposta aberta
	scanner := bufio.NewScanner(os.Stdin)
	//fmt.Printf("Deseja enviar os layouts manualmente?")
	fmt.Println("Deseja executar qual ação?")
	fmt.Println("1 - Envio")
	fmt.Println("2 - Retorno")
	fmt.Println("3 - Atualizas")

	ws.SetBaseUrl(os.Getenv("WS_BASE_URL"))
	ws.SetAuth(os.Getenv("WS_AUTH_USER"), os.Getenv("WS_AUTH_PASS"))

	scanner.Scan()
	varValue := scanner.Text()
	switch varValue {
	case "1":
		handSend.HandSend()
	case "2":
		handReturn.HandReturn()
	case "3":
		fmt.Println("Joga o link parças")
		//scanner.Scan()
		//attUrl := scanner.Text()
		err := update.Do("https://1drv.ms/u/c/6a491274b9a3b081/ET1Ly32-WolKt7l-bDr240sBU_N451LmvjGaTLK7WUVLog?e=7ph3hI")
		if err != nil {
			fmt.Println("Algo cagou:", err)
			log.Error("Algo cagou:", err)
		}

		defer log.Fatal("Stop service")

		fmt.Println("Atualizado com sucesso!")

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
