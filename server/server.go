package main

import (
	"fmt"
	"jogo/shared"
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Servidor struct {
	mu     sync.Mutex
	estado shared.EstadoJogo
}

// RegistrarJogador adiciona um jogador ao estado do jogo
func (s *Servidor) RegistrarJogador(id string, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.estado.Players == nil {
		s.estado.Players = make(map[string]shared.EstadoPlayer)
	}

	// Evita duplicação
	if _, existe := s.estado.Players[id]; !existe {
		s.estado.Players[id] = shared.EstadoPlayer{
			ID:       id,
			PosX:     1,
			PosY:     1,
			Sequence: 0,
		}
	}

	*reply = true
	return nil
}

// GetEstadoJogo retorna o estado atual do jogo
func (s *Servidor) GetEstadoJogo(id string, estado *shared.EstadoJogo) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	*estado = s.estado
	return nil
}

// AtualizarMovimento permite que o cliente envie um movimento
func (s *Servidor) AtualizarMovimento(mov shared.Movimento, reply *bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, ok := s.estado.Players[mov.ID]
	if !ok {
		*reply = false
		return fmt.Errorf("jogador não encontrado")
	}

	// Atualiza se a sequência for mais nova
	if mov.Sequence > player.Sequence {
		player.PosX = mov.PosX
		player.PosY = mov.PosY
		player.Sequence = mov.Sequence
		s.estado.Players[mov.ID] = player
	}

	*reply = true
	return nil
}

func main() {
	servidor := new(Servidor)

	rpc.Register(servidor)
	listener, err := net.Listen("tcp", "0.0.0.0:8932")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Servidor RPC iniciado em :8932")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
