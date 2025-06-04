// main.go - Loop principal do jogo
package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"time"
)

type User struct {
	ID				int			 // numero de identificação do usuario do servidor
    PosX, PosY      int          // posição atual do personagem do usuario
}

type CreateUserRequest struct {
    PosX, PosY      int          // posição atual do personagem do usuario
}

type GetUserRequest struct {
    ID int
}

func main() {
	if len(os.Args) < 2 {
        fmt.Println("Uso: go run client.go <endereço:porta>")
        os.Exit(1)
    }

    serverAddr := os.Args[1] // Ex: "localhost:1234"

    client, err := rpc.Dial("tcp", serverAddr)
    if err != nil {
        log.Fatal("Erro ao conectar:", err)
    }

	createReq := CreateUserRequest{PosX: }



	// Inicializa a interface (termbox)
	interfaceIniciar()
	defer interfaceFinalizar()

	// Usa "mapa.txt" como arquivo padrão ou lê o primeiro argumento
	mapaFile := "mapa.txt"
	if len(os.Args) > 1 {
		mapaFile = os.Args[1]
	}

	// Inicializa o jogo
	jogo := jogoNovo()
	if err := jogoCarregarMapa(mapaFile, &jogo); err != nil {
		panic(err)
	}

	// Desenha o estado inicial do jogo
	interfaceDesenharJogo(&jogo)

// Canal de parada
	stop := make(chan struct{})

	// Goroutine que move o inimigo sozinho
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				InimigoMover(&jogo)
				interfaceDesenharJogo(&jogo)
			case <-stop:
				return
			}
		}
	}()

	// Loop principal
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo); !continuar {
			close(stop)
			break
		}
		interfaceDesenharJogo(&jogo)
	}
}