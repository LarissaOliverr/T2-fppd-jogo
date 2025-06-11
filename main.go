// main.go - Loop principal do jogo
package main

import (
	"bufio"
	"fmt"
	"jogo/shared"
	"log"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mu          sync.Mutex
	estadoAtual shared.EstadoJogo
)


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Use: jogo.exe <playerID>")
		return
	}
	id := os.Args[1]

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite o endereco do servidor: ")
	endereco, _ := reader.ReadString('\n')
	endereco = strings.TrimSpace(endereco)

	client, err := rpc.Dial("tcp", endereco)
	if err != nil {
		log.Fatal("Erro ao conectar ao servidor RPC:", err)
	}

	var ok bool
	err = client.Call("Servidor.RegistrarJogador", id, &ok)
	if err != nil || !ok {
		log.Fatalf("Erro ao registrar jogador.")
	}

	interfaceIniciar()
	defer interfaceFinalizar()

	mapaFile := "mapa.txt"
	if len(os.Args) > 2 {
		mapaFile = os.Args[2]
	}

	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	jogo.ID = id
	jogo.Cliente = client

	go func() {
		for {
			var estado shared.EstadoJogo
			err := client.Call("Servidor.GetEstadoJogo", id, &estado)
			if err != nil {
				log.Println("Erro ao obter estado do jogo:", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			mu.Lock()
			estadoAtual = estado
			mu.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Loop principal
	sequence := 0
	stop := make(chan struct{})

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				InimigoMover(&jogo)
				mu.Lock()
				estado := estadoAtual
				mu.Unlock()
				interfaceDesenharJogo(&jogo, estado)
			case <-stop:
				return
			}
		}
	}()

	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, client, id, &sequence); !continuar {
			close(stop)
			break
		}
		jogoAtualizarEstadoMultiplayer(&jogo)
		mu.Lock()
		estado := estadoAtual
		mu.Unlock()
		interfaceDesenharJogo(&jogo, estado)
	}
}
