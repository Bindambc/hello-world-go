package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)
import "os"
import "time"
import "bufio"
import "strconv"

const delay = 5
const ok = 200
const tentativas = 3

func main() {

	exibeIntroducao()

	for {
		fmt.Println("1 - Iniciar Monitoramento")
		fmt.Println("2 - Exibir Log")
		fmt.Println("0 - Sair do programa")

		comando := getComando()

		switch comando {
		case 1:
			fmt.Println("Iniciando o monitoramento")
			monitoramento()
		case 2:
			fmt.Println("Exibindo logs")
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa")
			os.Exit(0)
		default:
			fmt.Println("Comando incorreto")
			os.Exit(-1)
		}
	}

}

func monitoramento() {

	sites := leSitesArquivo()
	fmt.Println(sites)

	for i := 0; i <= tentativas; i++ {
		for _, site := range sites {

			fmt.Println("Monitorando...")

			testaSite(site)

		}

		time.Sleep(delay * time.Second)
	}
}

func testaSite(site string) {
	response, _ := http.Get(site)

	if response != nil {

		if response.StatusCode == ok {
			fmt.Println("Site ", site, "está online")
			registraLog(site, true)
		} else {
			fmt.Println("Site", site, "está com problemas", response.StatusCode)
			registraLog(site, false)
		}

	} else {
		fmt.Println("Site", site, "não está respondendo")
		registraLog(site, false)
	}
}

func getComando() int8 {
	var comando int8

	_, err := fmt.Scan(&comando)

	if err != nil {
		fmt.Println("Erro inesperado:", err)
		return 0
	}
	fmt.Println("O comando escolhido foi", comando)
	return comando
}

func exibeIntroducao() {
	var nome = "Mauricio"
	fmt.Println("Olá", nome)
}

func leSitesArquivo() []string {
	arquivo, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println(err)
	}

	leitor := bufio.NewReader(arquivo)

	var linhas []string

	for {
		linha, err := leitor.ReadString('\n')

		if err == io.EOF {
			break
		}

		linha = strings.TrimSpace(linha)
		log.Default().Println(linha)
		linhas = append(linhas, linha)
	}

	err = arquivo.Close()
	if err != nil {
		log.Fatal(err)
	}

	return linhas
}

func registraLog(site string, status bool) {

	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}

	_, err = arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " | " + site + " - online: " + strconv.FormatBool(status) + "\n")
	if err != nil {
		return
	}

	err = arquivo.Close()
	if err != nil {
		fmt.Println(err)
	}

}

func imprimeLogs() {
	arquivo, err := os.ReadFile("log.txt")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(arquivo))
}
