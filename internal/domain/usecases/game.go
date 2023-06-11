package usecases

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/hattonuri/soa-hmw-2/internal/generated/proto/base"
	"github.com/streadway/amqp"
)

type Game struct {
	room *Room

	clients []*GameClient

	game  *base.Game
	wg    sync.WaitGroup
	mutex sync.Mutex
}

func (g *Game) NotifyPlayerGroup(event *base.Event, selector func(int) bool, wait bool) {
	cnt := 0
	for i := 0; i < len(g.clients); i++ {
		if selector(i) {
			cnt += 1
		}
	}
	if wait {
		g.mutex.Lock()
		g.wg.Add(cnt)
		defer g.wg.Wait()
		g.game.SelectedTargetIds = []uint64{}
		g.mutex.Unlock()
	}
	if event.NeedsAnswer {
		event.Text = event.Text + "\nEnter GetPlayers for get players"
	}
	for i, client := range g.clients {
		if selector(i) {
			log.Println(event.Text, client.game.room.Players[i].Name)
			event.RoomId = g.room.RoomId
			err := client.srv.Send(event)
			if err != nil {
				log.Fatalf("failed to send event to client: %v", err)
				return
			}
		}
	}
}

func (g *Game) GetMaxVotedIndex() int {
	cnt := make([]int, len(g.clients))
	result := 0
	for _, cur := range g.game.SelectedTargetIds {
		cnt[cur]++
		if cnt[result] < cnt[cur] {
			result = int(cur)
		}
	}
	return result
}

func (g *Game) RunChat() {
	amqpConn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("amqp dial: %v", err)
		return
	}
	amqpChannel, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("amqp channel: %v", err)
		return
	}
	defer amqpChannel.Close()
	q, err := amqpChannel.QueueDeclare(
		fmt.Sprintf("room-%d", g.room.RoomId),
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf("amqp queue declare: %v", err)
		return
	}

	log.Println("Listening to ", q.Name)
	messages, err := amqpChannel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("amqp channel consume: %v", err)
		return
	}
	for chatMessage := range messages {
		msg := string(chatMessage.Body)
		if strings.Contains(msg, "-ALLCHAT-") {
			tokens := strings.Split(msg, "-ALLCHAT-")
			amqpSender := tokens[0]
			text := tokens[1]
			g.NotifyPlayerGroup(
				&base.Event{
					Text: amqpSender + ":" + text,
				},
				func(i int) bool {
					player := g.game.Players[i]
					if player.Alive && player.Name != amqpSender {
						fmt.Println("ROFLAN ", player.Name)
					}
					return player.Alive && player.Name != amqpSender
				},
				false,
			)
		} else if strings.Contains(msg, "-MAFIACHAT-") {
			tokens := strings.Split(msg, "-MAFIACHAT-")
			amqpSender := tokens[0]
			text := tokens[1]
			g.NotifyPlayerGroup(
				&base.Event{
					Text: amqpSender + ":" + text,
				},
				func(i int) bool {
					player := g.game.Players[i]
					if player.Alive && player.Name != amqpSender && player.Role == base.Role_MAFIA {
						fmt.Println("ZZZ ", player.Name)
					}
					return player.Alive && player.Name != amqpSender && player.Role == base.Role_MAFIA
				},
				false,
			)
		}
	}
}

func (g *Game) Run() {
	defer close(g.room.done)

	go g.RunChat()
	for {
		cntMafia := 0
		cntNotMafia := 0
		for index := range g.clients {
			player := g.game.Players[index]
			if !player.Alive {
				continue
			}
			log.Println(index, player.Name, player.Role.String())
			if player.Role == base.Role_MAFIA {
				cntMafia++
			} else {
				cntNotMafia++
			}
		}
		log.Printf("Room %d, %d citisens vs %d mafia", g.room.RoomId, cntNotMafia, cntMafia)
		if cntMafia == 0 {
			g.NotifyPlayerGroup(
				&base.Event{
					Text: "Citizens win!",
				},
				func(i int) bool {
					return true
				},
				false,
			)
			return
		}
		if cntNotMafia <= cntMafia {
			g.NotifyPlayerGroup(
				&base.Event{
					Text: "Mafia win!",
				},
				func(i int) bool {
					return true
				},
				false,
			)
			return
		}
		g.NotifyPlayerGroup(
			&base.Event{
				Text: "Night...",
			},
			func(i int) bool {
				return true
			},
			false,
		)
		whomToHeal := -1
		g.NotifyPlayerGroup(
			&base.Event{
				Text:        "Medic protect...",
				RoomId:      g.room.RoomId,
				NeedsAnswer: true,
			},
			func(i int) bool {
				player := g.game.Players[i]
				return player.Alive && player.Role == base.Role_MEDIC
			},
			true,
		)
		whomToHeal = g.GetMaxVotedIndex()
		log.Printf("Try get protectID: %d heal %s", whomToHeal, g.game.Players[whomToHeal].Name)
		g.NotifyPlayerGroup(
			&base.Event{
				Text:        "Mafia kill...",
				RoomId:      g.room.RoomId,
				NeedsAnswer: true,
			},
			func(i int) bool {
				player := g.game.Players[i]
				return player.Alive && player.Role == base.Role_MAFIA
			},
			true,
		)
		maxIndex := g.GetMaxVotedIndex()
		log.Printf("Try get killedID: %d kill %s", maxIndex, g.game.Players[maxIndex].Name)
		if maxIndex == whomToHeal {
			g.NotifyPlayerGroup(
				&base.Event{
					Text: "All right...all people alive",
				},
				func(i int) bool {
					return true
				},
				false,
			)
		} else {
			g.mutex.Lock()
			g.game.Players[maxIndex].Alive = false
			g.mutex.Unlock()
			g.NotifyPlayerGroup(
				&base.Event{
					Text: fmt.Sprintf("User %s has been killed", g.game.Players[maxIndex].Name),
				},
				func(i int) bool {
					return true
				},
				false,
			)
		}
		g.NotifyPlayerGroup(
			&base.Event{
				Text:        "Time to vote...",
				NeedsAnswer: true,
			},
			func(i int) bool {
				return g.game.Players[i].Alive
			},
			true,
		)
		maxIndex = g.GetMaxVotedIndex()
		g.NotifyPlayerGroup(
			&base.Event{
				Text: fmt.Sprintf(
					"User %s has been kicked and his role: %s",
					g.game.Players[maxIndex].Name,
					g.game.Players[maxIndex].Role.String()),
			},
			func(i int) bool {
				return true
			},
			false,
		)
		g.mutex.Lock()
		g.game.Players[maxIndex].Alive = false
		g.mutex.Unlock()
	}
}

type GameClient struct {
	game *Game
	srv  base.MafiaService_JoinServer
}
