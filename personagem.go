// personagem.go - Funções para movimentação e ações do personagem
package main

// Atualiza a posição do personagem com base na tecla pressionada (WASD)
func personagemMover(tecla rune, jogo *Jogo) {
	dx, dy := 0, 0
	switch tecla {
	case 'w': dy = -1 // Move para cima
	case 'a': dx = -1 // Move para a esquerda
	case 's': dy = 1  // Move para baixo
	case 'd': dx = 1  // Move para a direita
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

			// Verifica se é um botão (baseado no símbolo)
			if elem.simbolo == Botao.simbolo {
				jogo.BotaoBool = !jogo.BotaoBool

				if elem.cor == CorVerde {
					elem.cor = CorVermelho
				} else {
					elem.cor = CorVerde
				}

				ativarPortal(jogo)
			}
		}
	}
}

// Processa o evento do teclado e executa a ação correspondente
func personagemExecutarAcao(ev EventoTeclado, jogo *Jogo) bool {
	switch ev.Tipo {
	case "sair":
		// Retorna false para indicar que o jogo deve terminar
		return false
	case "interagir":
		// Executa a ação de interação
		personagemInteragir(jogo)
	case "mover":
		// Move o personagem com base na tecla
		personagemMover(ev.Tecla, jogo)
	}
	return true // Continua o jogo
}
