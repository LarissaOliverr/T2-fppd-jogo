// jogo.go - Funções para manipular os elementos do jogo, como carregar o mapa e mover o personagem
package main

import (
	"bufio"
	"os"
)

// Elemento representa qualquer objeto do mapa (parede, personagem, vegetação, etc)
type Elemento struct {
	simbolo   rune
	cor       Cor
	corFundo  Cor
	tangivel  bool // Indica se o elemento bloqueia passagem
}

// Jogo contém o estado atual do jogo
type Jogo struct {
	Mapa            [][]Elemento // grade 2D representando o mapa
	PosX, PosY      int          // posição atual do personagem
	UltimoVisitado  Elemento     // elemento que estava na posição do personagem antes de mover
	StatusMsg       string       // mensagem para a barra de status

	BotaoBool		bool		 // boledano para ativação do botão
	PortalAtivo		bool		 // variavel que verifica se os personagens clicaram no botão
	InimigoDir      int          // direcional para movimento do inimigo

}

// Elementos visuais do jogo
var (
	//elementos base
	Personagem = Elemento{'☺', CorCinzaEscuro, CorPadrao, true}
	Inimigo    = Elemento{'☠', CorVermelho, CorPadrao, true}
	Parede     = Elemento{'▤', CorParede, CorFundoParede, true}
	Vegetacao  = Elemento{'♣', CorVerde, CorPadrao, false}
	Vazio      = Elemento{' ', CorPadrao, CorPadrao, false}

	//elementos adicionais
	Portal	   = Elemento{'✷', CorPadrao, CorPadrao, true}
	Botao	   = Elemento{'⏺', CorVermelho, CorPadrao, true}
)

// Cria e retorna uma nova instância do jogo
func jogoNovo() Jogo {
	// O ultimo elemento visitado é inicializado como vazio
	// pois o jogo começa com o personagem em uma posição vazia
	return Jogo{
		UltimoVisitado: Vazio,
		BotaoBool: false,
		PortalAtivo: false,
		InimigoDir: 1,
	}
}

// Lê um arquivo texto linha por linha e constrói o mapa do jogo
func jogoCarregarMapa(nome string, jogo *Jogo) error {
	arq, err := os.Open(nome)
	if err != nil {
		return err
	}
	defer arq.Close()

	scanner := bufio.NewScanner(arq)
	y := 0
	for scanner.Scan() {
		linha := scanner.Text()
		var linhaElems []Elemento
		for x, ch := range linha {
			e := Vazio
			switch ch {
			case Portal.simbolo:
				e = Portal
			case Botao.simbolo:
				e= Botao
			case Parede.simbolo:
				e = Parede
			case Inimigo.simbolo:
				e = Inimigo
			case Vegetacao.simbolo:
				e = Vegetacao
			case Personagem.simbolo:
				jogo.PosX, jogo.PosY = x, y // registra a posição inicial do personagem
			}
			linhaElems = append(linhaElems, e)
		}
		jogo.Mapa = append(jogo.Mapa, linhaElems)
		y++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Verifica se o personagem pode se mover para a posição (x, y)
func jogoPodeMoverPara(jogo *Jogo, x, y int) bool {
	// Verifica se a coordenada Y está dentro dos limites verticais do mapa
	if y < 0 || y >= len(jogo.Mapa) {
		return false
	}

	// Verifica se a coordenada X está dentro dos limites horizontais do mapa
	if x < 0 || x >= len(jogo.Mapa[y]) {
		return false
	}

	// Verifica se o elemento de destino é tangível (bloqueia passagem)
	if jogo.Mapa[y][x].tangivel {
		return false
	}

	// Pode mover para a posição
	return true
}

// Move um elemento para a nova posição
func jogoMoverElemento(jogo *Jogo, x, y, dx, dy int) {
	nx, ny := x+dx, y+dy

	// Obtem elemento atual na posição
	elemento := jogo.Mapa[y][x] // guarda o conteúdo atual da posição

	jogo.Mapa[y][x] = jogo.UltimoVisitado     // restaura o conteúdo anterior
	jogo.UltimoVisitado = jogo.Mapa[ny][nx]   // guarda o conteúdo atual da nova posição
	jogo.Mapa[ny][nx] = elemento              // move o elemento
}

func InimigoMover(jogo *Jogo) {
	var novosInimigos []struct{ X, Y int }

	// Primeiro, encontre todos os inimigos
	for y, linha := range jogo.Mapa {
		for x, elem := range linha {
			if elem.simbolo == Inimigo.simbolo {
				novosInimigos = append(novosInimigos, struct{ X, Y int }{x, y})
			}
		}
	}

	// Depois, tente mover cada um
	for _, pos := range novosInimigos {
		x, y := pos.X, pos.Y
		nx := x + jogo.InimigoDir

		if nx < 0 || nx >= len(jogo.Mapa[y]) || jogo.Mapa[y][nx].tangivel || jogo.Mapa[y][nx].simbolo == Inimigo.simbolo {
			// Inverte direção se não puder andar
			jogo.InimigoDir *= -1
			nx = x + jogo.InimigoDir
			if nx < 0 || nx >= len(jogo.Mapa[y]) || jogo.Mapa[y][nx].tangivel || jogo.Mapa[y][nx].simbolo == Inimigo.simbolo {
				continue // ainda bloqueado
			}
		}

		// Move o inimigo
		jogo.Mapa[y][x] = Vazio
		jogo.Mapa[y][nx] = Inimigo
	}
}


