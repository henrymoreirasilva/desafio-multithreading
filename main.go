package main

import (
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		println("Informe o número do CEP como parâmetro.")
		return
	}
	c1 := make(chan string)
	c2 := make(chan string)

	for _, cep := range os.Args[1:] {
		go func() {
			req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
			if err != nil {
				c1 <- ""
			}
			defer req.Body.Close()
			res, err := io.ReadAll(req.Body)
			if err != nil {
				c1 <- ""
			}
			c1 <- string(res)

		}()

		go func() {
			req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
			if err != nil {
				c2 <- ""
			}
			defer req.Body.Close()
			res, err := io.ReadAll(req.Body)
			if err != nil {
				c2 <- ""
			}
			c2 <- string(res)
		}()

		select {
		case msg1 := <-c1:
			print("####\n")
			print("Resposta " + cep + " de: https://brasilapi.com.br \n\n")
			println(msg1)
			print("\n\n")
		case msg2 := <-c2:
			print("####\n")
			println("Resposta para " + cep + " de: http://viacep.com.br \n\n")
			println(msg2)
			print("\n\n")
		case <-time.After(time.Second):
			println("Timeout\n\n")
		}
	}
}
