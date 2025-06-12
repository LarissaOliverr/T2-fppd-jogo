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

type ComandoInteracao struct {
    ID       string
    Sequence int
    Tipo     string 
}

type Movimento struct {
	ID       string
	PosX   int
	PosY   int
	Sequence int
}