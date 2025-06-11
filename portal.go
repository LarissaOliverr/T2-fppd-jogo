package main


func ativarPortal(jogo *Jogo) {
	for y := range jogo.Mapa {
		for x := range jogo.Mapa[y] {
			elem := &jogo.Mapa[y][x]

			if elem.simbolo == Portal.simbolo {
				if jogo.BotaoBool {
					elem.cor = CorAzul
					jogo.PortalAtivo = true
					jogo.StatusMsg = "O portal foi ativado!"
				} else {
					elem.cor = CorPadrao
					jogo.PortalAtivo = false
					jogo.StatusMsg = "O portal foi desativado."
				}
			}
		}
	}
}
