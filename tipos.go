// tipos compartilhados entre cliente e servidor

package main

type MovimentoArgs struct {
	ID      string
	Direcao rune
}

type EstadoArgs struct {
	ID string
}

type Jogador struct {
	ID     string
	PosX   int
	PosY   int
	Vida   int
	Status string
}

type EstadoJogo struct {
	Jogadores map[string]Jogador
}