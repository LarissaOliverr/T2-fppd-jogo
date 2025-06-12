// personagem.go - Funções para movimentação e ações do personagem
package main

import (
	"jogo/shared"
	"log"
	"net/rpc"
)

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w':
		dy = -1 // Move para cima
	case 'a':
		dx = -1 // Move para a esquerda
	case 's':
		dy = 1 // Move para baixo
	case 'd':
		dx = 1 // Move para a direita
	}

	nx, ny := jogo.PosX+dx, jogo.PosY+dy
	// Verifica se o movimento é permitido e realiza a movimentação
	if jogoPodeMoverPara(jogo, nx, ny) {
		jogoMoverElemento(jogo, jogo.PosX, jogo.PosY, dx, dy)
		jogo.PosX, jogo.PosY = nx, ny
	}
}

// Define o que ocorre quando o jogador pressiona a tecla de interação
// Neste exemplo, apenas exibe uma mensagem de status
// Você pode expandir essa função para incluir lógica de interação com objetos
// Define o que ocorre quando o jogador pressiona a tecla de interação

func personagemInteragir(jogo *Jogo) {
	// Verifica os 4 blocos ao redor
	direcoes := []struct{ dx, dy int }{
		{0, -1}, // cima
		{0, 1},  // baixo
		{-1, 0}, // esquerda
		{1, 0},  // direita
	}

	for _, d := range direcoes {
		x, y := jogo.PosX+d.dx, jogo.PosY+d.dy

		if y >= 0 && y < len(jogo.Mapa) && x >= 0 && x < len(jogo.Mapa[y]) {
			elem := &jogo.Mapa[y][x]

			// Verifica se é um botão
			if elem.simbolo == Botao.simbolo {
				jogo.BotaoBool = !jogo.BotaoBool

				if elem.cor == CorVerde {
					elem.cor = CorVermelho
				} else {
					elem.cor = CorVerde
				}

				ativarPortal(jogo)

				var ack bool
				_ = jogo.Cliente.Call("Servidor.AtualizarEstadoLogico", shared.EstadoJogo{
					BotaoAtivo:  jogo.BotaoBool,
					PortalAtivo: jogo.PortalAtivo,
				}, &ack)

				// Verifica se é um portal
				if elem.simbolo == Portal.simbolo {
					if jogo.PortalAtivo == true {
						jogo.StatusMsg = "Você escapou em segurança!"
					} else {
						jogo.StatusMsg = "Portal está desativado..."
					}
				}
			}

		}
	}

}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo, client *rpc.Client, id string, sequence *int) bool {
	switch ev.Tipo {
	case "sair":
		var ack bool
		_ = client.Call("Servidor.DesconectarJogador", id, &ack)
		return false

	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move localmente
		personagemMover(ev.Tecla, jogo)

		// Atualiza sequence e envia ao servidor
		*sequence = *sequence + 1
		mov := shared.Movimento{
			ID:       id,
			PosX:     jogo.PosX,
			PosY:     jogo.PosY,
			Sequence: *sequence,
		}
		var ack bool
		err := client.Call("Servidor.AtualizarMovimento", mov, &ack)
		if err != nil {
			log.Println("Erro ao atualizar posição no servidor:", err)
		}
	}
	return true // Continua o jogo
}
