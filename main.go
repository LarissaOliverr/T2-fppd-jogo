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

var (
	mu          sync.Mutex
	estadoAtual shared.EstadoJogo
)


func main() {
if len(os.Args) < 2 {
		fmt.Println("<playerID>")
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



	// Desenha o estado inicial do jogo
	mu.Lock()
	estado := estadoAtual
	mu.Unlock()

	interfaceDesenharJogo(&jogo, estado)


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
				interfaceDesenharJogo(&jogo, estado)
			case <-stop:
				return
			}
		}
	}()

	// Loop principal
	sequence := 0
	for {
		evento := interfaceLerEventoTeclado()
		if continuar := personagemExecutarAcao(evento, &jogo, client, id, &sequence); !continuar {
			close(stop)
			break
		}
		interfaceDesenharJogo(&jogo, estado)
	}

}