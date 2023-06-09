package usecases

import (
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/hattonuri/soa-hmw-2/internal/generated/proto/base"
)

type Room struct {
	*base.Room
	game *Game
	mut  sync.Mutex
	wg   sync.WaitGroup
	done chan bool
}

func (r *Room) GenerateIntoPlayers(players []*base.Player) []*base.Player {
	roles := make([]base.Role, 0)
	roles = append(roles, base.Role_MAFIA)
	roles = append(roles, base.Role_MEDIC)
	for i := 2; i < len(players); i++ {
		roles = append(roles, base.Role_CITIZEN)
	}
	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})
	for i := 0; i < len(players); i++ {
		players[i].Role = roles[i]
		players[i].Alive = true
	}
	return players
}

func (r *Room) Join(player *base.Player, srv base.MafiaService_JoinServer) error {
	err := func() error {
		r.mut.Lock()
		defer r.mut.Unlock()
		if uint64(len(r.Players)) == r.MaxPlayers {
			return fmt.Errorf("room %d overflow", r.RoomId)
		}
		if uint64(len(r.Players)) == 0 {
			r.createGame()
			r.wg.Add(int(r.MaxPlayers))
			r.Players = make([]*base.Player, 0)
		}
		r.Players = append(r.Players, player)
		log.Printf("New player in room %d (%d/%d)! hello %s!!!",
			r.RoomId,
			len(r.Players),
			r.MaxPlayers,
			player.Name,
		)
		r.game.mutex.Lock()
		r.game.clients = append(r.game.clients, &GameClient{
			game: r.game,
			srv:  srv,
		})
		r.game.mutex.Unlock()
		if uint64(len(r.Players)) == r.MaxPlayers {
			r.Game.Players = r.GenerateIntoPlayers(r.Players)
			go func() {
				r.game.Run()
			}()
		}
		return nil
	}()
	if err != nil {
		return err
	}
	r.wg.Done()
	log.Printf("Player %s in room %d wait start game...", player.Name, r.RoomId)
	r.wg.Wait()
	return nil
}

func (r *Room) AddVote(selectedTargetId uint64) {
	r.game.mutex.Lock()
	r.Game.SelectedTargetIds = append(r.Game.SelectedTargetIds, selectedTargetId)
	r.game.wg.Done()
	r.game.mutex.Unlock()
}

func (r *Room) createGame() {
	r.game = &Game{
		room:    r,
		game:    r.Game,
		wg:      sync.WaitGroup{},
		mutex:   sync.Mutex{},
		clients: make([]*GameClient, 0),
	}
}
