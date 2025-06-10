package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
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

type UserService struct {
    mu    sync.Mutex
    users map[int]User
    nextID int
}

func (s *UserService) CreateUser(req *CreateUserRequest, resp *User) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.nextID++
    user := User{ID: s.nextID, PosX: req.PosX, PosY: req.PosY}
    s.users[user.ID] = user
    *resp = user
    return nil
}

func (s *UserService) GetUser(req *GetUserRequest, resp *User) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    user, ok := s.users[req.ID]
    if !ok {
        return errors.New("usuário não encontrado")
    }
    *resp = user
    return nil
}

func main() {
    service := &UserService{
        users: make(map[int]User),
        nextID: 0,
    }

    rpc.Register(service)
    listener, err := net.Listen("tcp", ":8932")
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