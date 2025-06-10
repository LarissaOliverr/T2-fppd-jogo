// main.go - Loop principal do jogo
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
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
		fmt.Println("Uso: go run main.go <playerID>")
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

	// registra o jogador apos conectar
	var ok bool
	err = client.Call("Servidor.RegistrarJogador", id, &ok)
	if err != nil || !ok {
		log.Fatalf("Erro ao registrar jogador.")
	}


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