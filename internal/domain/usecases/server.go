package usecases

import (
	"fmt"
	"log"
	"sync"

	"github.com/hattonuri/soa-hmw-2/internal/config"
	"github.com/hattonuri/soa-hmw-2/internal/generated/proto/base"
)

type Server struct {
	mtx         sync.Mutex
	rooms       map[uint64]*Room
	roomCounter uint64
	maxPlayers  uint64
}

func (s *Server) GetPlayersList(room *base.Room) (*base.Players, error) {
	value, ok := s.rooms[room.RoomId]
	if !ok {
		return nil, fmt.Errorf("not found %d", room.RoomId)
	}
	resultSlice := make([]*base.Player, 0)
	for index, key := range value.Game.Players {
		if key.Alive {
			resultSlice = append(resultSlice, &base.Player{Name: key.Name + fmt.Sprintf(" (%d)", index)})
		}
	}
	result := &base.Players{
		Players: resultSlice,
	}
	return result, nil
}

func NewServer(c *config.Server) *Server {
	return &Server{
		mtx:         sync.Mutex{},
		rooms:       make(map[uint64]*Room, 0),
		roomCounter: 0,
		maxPlayers:  c.MaxPlayers,
	}
}

func (s *Server) CreateRoom() uint64 {
	var room = new(Room)
	var game = new(base.Game)
	game = &base.Game{}
	room.Room = &base.Room{
		Game:       game,
		RoomId:     s.roomCounter,
		MaxPlayers: uint64(s.maxPlayers),
	}
	room.done = make(chan bool)
	s.rooms[s.roomCounter] = room
	s.roomCounter += 1
	return s.roomCounter - 1
}

func (s *Server) getFreeRoomID() uint64 {
	for key, element := range s.rooms {
		if len(element.Players) < int(element.MaxPlayers) {
			return key
		}
	}
	return s.CreateRoom()
}

func (s *Server) getFreeRoom() *Room {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	return s.rooms[s.getFreeRoomID()]
}

func (s *Server) Join(in *base.Player, srv base.MafiaService_JoinServer) error {
	room := s.getFreeRoom()
	if err := room.Join(in, srv); err != nil {
		log.Fatal(err)
		return err
	}
	<-room.done
	return nil
}

func (s *Server) ResponseEvent(event *base.Event) {
	room := s.rooms[event.RoomId]
	room.AddVote(event.SelectedTargetId)
}
