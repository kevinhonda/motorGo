package sendStuffs

import (
	"encoding/csv"
	"sync"
	"time"

	//"errors"
	"fmt"
	"os"

	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/kafy11/gosocket/log"

	//"github.com/valyala/fasthttp"

	"motorv2/pkg/db"
	awsC "motorv2/src/awsControllers"

	//"reflect"
	"unicode/utf8"
)

func BusSend(layout, query string, busType ...string) {
	log.Info("BusSend - LAYOUT: ", layout)
	fmt.Println("BusSend - LAYOUT: " + layout)

	//Inicio do temporizador
	execStart := time.Now()
	//log.Info("BusSend - QUERY: ", query)
	//Consulta no banco
	rows, err := db.SqlRun(query, layout)
	if err != nil {
		log.Error(fmt.Sprintf("SqlRun %s - Erro no banco: %s", layout, err))
		fmt.Printf("%s - Erro no banco: %s", layout, err)
		return
	}
	defer rows.Close()

	//Transformando em CSV
	csvName, rowCount := rowToCsv(rows, layout)

	//Logando tempo de consulta
	execEnd := time.Now()
	execDuration := execEnd.Sub(execStart).Seconds()
	fmt.Println(fmt.Printf("%s - O tempo total da consulta no banco levou %f segundos para executar.\n", layout, execDuration))
	log.Info(fmt.Sprintf("%s - O tempo total da consulta no banco levou %f segundos para executar.\n", layout, execDuration))

	//Listando CSVs do layout
	fmt.Println(csvName)
	filesName := GetFilesName(layout)

	count := len(filesName)
	counts := strconv.Itoa(count)
	rowCounts := strconv.Itoa(rowCount)

	fmt.Println(fmt.Printf("%s - filesName: Total de arquivos %s", layout, counts))
	log.Info("fileNames - ", filesName)
	log.Info("Tamanho desse fileNames é ", count)

	signedStart := time.Now()
	if rowCount > 0 {
		preSigneds := make(map[string]interface{})
		//Gerando pre signed
		if len(busType) > 0 {
			preSigneds, _ = awsC.GetInfos("Adm_get_pre_signed?table_name=" + layout + "&table_size=" + rowCounts + "&block_size=" + strconv.Itoa(100000) + "&mini_bus=y")
		} else {
			preSigneds, _ = awsC.GetInfos("Adm_get_pre_signed?table_name=" + layout + "&table_size=" + rowCounts + "&block_size=" + strconv.Itoa(100000))
		}

		signedEnd := time.Now()
		signedDuration := signedEnd.Sub(signedStart).Seconds()
		//fmt.Println(fmt.Sprintf("%s - O tempo total de gerar preSigneds foi de  %f segundos.\n", layout, signedDuration))
		fmt.Printf("%s - O tempo total de gerar preSigneds foi de  %f segundos.\n", layout, signedDuration)
		log.Info(fmt.Sprintf("%s - O tempo total de gerar preSigneds foi de  %f segundos.\n", layout, signedDuration))
		//log.Info("Link gerado com sucesso")

		//Enviando CSVs para s3
		var wg sync.WaitGroup
		sendStart := time.Now()

		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				startNe := time.Now()
				awsC.SendToS3(filesName[i], preSigneds[strconv.Itoa(i+1)].(string))
				//fmt.Printf("Uau mas que belo pre signed: %s", preSigneds[strconv.Itoa(i+1)].(string))
				layoutEnd := time.Now()
				layoutDuration := layoutEnd.Sub(startNe).Seconds()
				fmt.Println(fmt.Printf("%s - O tempo de envio do %s - %s levou %f segundos para executar.\n", layout, filesName[i], strconv.Itoa(i), layoutDuration))
				log.Info(fmt.Sprintf("%s - O tempo de envio do %s - %s levou %f segundos para executar.\n", layout, filesName[i], strconv.Itoa(i), layoutDuration))
			}(i)
		}

		wg.Wait()

		sendEnd := time.Now()
		sendDuration := sendEnd.Sub(sendStart).Seconds()
		fmt.Println(fmt.Printf("%s - O tempo total do envio para S3 foi de %f segundos para executar.\n", layout, sendDuration))
		log.Info(fmt.Sprintf("%s - O tempo total do envio para S3 foi de %f segundos para executar.\n", layout, sendDuration))

		execDuration = sendEnd.Sub(execStart).Seconds()
		fmt.Println(fmt.Printf("BusSend - O tempo total do envio %s foi de %f segundos para executar.\n", layout, execDuration))
		log.Info(fmt.Sprintf("BusSend - O tempo total do envio %s foi de %f segundos para executar.\n", layout, execDuration))

		//rowCountStr := strconv.Itoa(rowCount)

		//splitCsv(layout, csvName, rowCountStr)
		//log.Info("Link gerado com sucesso")
		// Nome do arquivo CSV de saída

		// Tente excluir o arquivo
		/*
			err = os.Remove(csvName)
			if err != nil {
				log.Error(fmt.Printf("%s - Erro ao apagar o arquivo: %s", layout, err))
				fmt.Println(fmt.Printf("%s - Erro ao apagar o arquivo: %s", layout, err))
			}
		*/
		log.Info("Layout " + layout + " enviado com sucesso")
	} else {
		log.Info("O layout " + layout + "não tem dados a ser enviado")
		log.Info("RMV " + filesName[0])
		err = os.Remove(filesName[0])
		if err != nil {
			log.Error(fmt.Sprintf("%s RMV - %s", layout, err))
			return
		}
	}

}

func GetFilesName(layout string) []string {
	// Obter o diretório atual
	folderName, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter o diretório atual:", err)
		log.Error("Erro ao obter o diretório atual:", err)
		return nil
	}

	// Ler o diretório
	files, err := os.ReadDir(folderName)
	if err != nil {
		fmt.Println("Erro ao ler o diretório:", err)
		log.Error("Erro ao ler o diretório:", err)
		return nil
	}
	fmt.Println("Files - ", files)
	log.Info("Files - ", files)
	// Regex para verificar o padrão `_número`
	//numberPattern := regexp.MustCompile(`_\d+`)
	numberPattern := regexp.MustCompile(`^` + layout + `_\d+.*$`)

	var fileNames []string
	for _, file := range files {
		// Ignorar diretórios
		if file.IsDir() {
			continue
		}

		// Verificar se o nome do arquivo começa com o layout e contém "_número"
		/*
			if strings.Contains(file.Name(), layout+"_") && numberPattern.MatchString(file.Name()) {
				fileNames = append(fileNames, file.Name())
			}
		*/

		if numberPattern.MatchString(strings.ToUpper(file.Name())) {
			fileNames = append(fileNames, file.Name())
		}

	}

	// Exibir os arquivos encontrados (opcional)
	fmt.Println("Arquivos encontrados no diretório:", folderName)
	for _, name := range fileNames {
		fmt.Println(name)
	}

	return fileNames
}

func rowToCsv(rows *sqlx.Rows, layout string) (string, int) {
	reg := regexp.MustCompile(`[^A-Za-z0-9ÁÀÂÃÉÊÍÓÔÕÚÇáàâãéêíóôõúç_.@,!?()\-:; ]`)
	fileCount := 1
	lineCount := 0
	totalLineCount := 0
	/*
		if !rows.Next() {
			// Fecha as rows e retorna indicando que está vazio
			rows.Close()
			return "", 0
		}
	*/
	fileName := createFileName(layout, fileCount)
	file, writer, err := createCsvFile(fileName)
	if err != nil {

		log.Error(layout, " - Não foi possível criar o arquivo CSV:", err)
		fmt.Println(layout, " - Não foi possível criar o arquivo CSV:", err)
		return "", 0
	}
	defer file.Close()
	defer writer.Flush()

	columns, err := rows.Columns()
	if err != nil {
		log.Error(layout, " - ", err)
		fmt.Println(layout, " - ", err)
		return "", 0
	}

	if err := writer.Write(columns); err != nil {
		log.Error(layout, " - Erro ao escrever cabeçalhos no CSV:", err)
		return "", 0
	}

	for rows.Next() {
		if lineCount > 0 && lineCount%100000 == 0 {
			writer.Flush()
			file.Close()

			fileCount++
			fileName = createFileName(layout, fileCount)
			file, writer, err = createCsvFile(fileName)
			if err != nil {
				//return lineCount, err
				return "", 0
			}
			defer file.Close()
			defer writer.Flush()

			if err := writer.Write(columns); err != nil {
				log.Error(layout, " - Erro ao escrever cabeçalhos no CSV:", err)
				//return lineCount, err
				return "", 0
			}
		}

		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Error(layout, " - ", err)
			//return lineCount, err
			return "", 0
		}

		record := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				record[i] = "NULL"
			} else {
				nVal := interfaceToUTF8(val)
				nVal = reg.ReplaceAllString(nVal, "")
				record[i] = fmt.Sprintf("%s", nVal)
				//log.Info(fmt.Sprintf("É UTF-8 válido? %t\n", utf8.ValidString(nVal)))
				//log.Info(fmt.Sprintf("Tipo (reflect): %v\n", reflect.TypeOf(nVal)))

				//log.Info(fmt.Sprintf("%s - %v", layout, val))
			}
		}

		//log.Info("CSV ", layout, " - ", record)

		if err := writer.Write(record); err != nil {
			log.Error(layout, " - Erro ao escrever registro no CSV:", err)
			//return lineCount, err
			return "", 0
		}

		lineCount++
		totalLineCount++
	}

	if err := rows.Err(); err != nil {
		log.Error(layout, " - ", err)
		//return lineCount, err
		return "", 0
	}

	log.Info("Dados exportados com sucesso")
	//return lineCount, nil
	return "", totalLineCount
}

func createFileName(layout string, count int) string {
	return layout + "_" + strconv.Itoa(count) + ".csv"
}

func createCsvFile(fileName string) (*os.File, *csv.Writer, error) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Error("Não foi possível criar o arquivo CSV:", err)
		return nil, nil, err
	}

	writer := csv.NewWriter(file)
	writer.Comma = ';' // Define o separador como ponto e vírgula
	return file, writer, nil
}

func interfaceToUTF8(val interface{}) string {
	switch v := val.(type) {
	case string:
		if utf8.ValidString(v) {
			return v
		}
		log.Error("string contém UTF-8 inválido")
		return ""
	case []byte:
		if utf8.Valid(v) {
			return string(v)
		}
		log.Error("[]byte contém UTF-8 inválido")
		return ""
	default:
		// Converte outros tipos para string
		str := fmt.Sprintf("%v", v)
		if !utf8.ValidString(str) {
			log.Error("a conversão resultou em UTF-8 inválido")
			return ""
		}
		return str
	}
}
