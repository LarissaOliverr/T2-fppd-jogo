package shared

type EstadoPlayer struct {
	ID       string
	PosX, PosY     int
	Sequence int
}

type EstadoJogo struct {
	Players      map[string]EstadoPlayer
	Mapa         [][]rune
	BotaoAtivo   bool
	PortalAtivo  bool
}


type Movimento struct {
	ID       string
	PosX   int
	PosY   int
	Sequence int
}