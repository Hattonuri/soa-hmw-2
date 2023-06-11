package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hattonuri/soa-hmw-2/internal/config"
	"github.com/hattonuri/soa-hmw-2/internal/generated/proto/base"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func FNV32a(text string) uint32 {
	algorithm := fnv.New32a()
	algorithm.Write([]byte(text))
	return algorithm.Sum32()
}

func main() {
	cfg := &config.Client{}
	config.InitClient(cfg)
	rand.Seed(time.Now().UnixNano() + int64(FNV32a(cfg.Hostname)))
	user := cfg.User
	role := cfg.Role
	if role == "bot" {
		user += fmt.Sprintf("%d", rand.Int31())
	}
	addr := ":8080"
	if len(cfg.Hostname) > 0 {
		addr = "server" + addr
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := base.NewMafiaServiceClient(conn)
	stream, err := client.Join(context.Background(), &base.Player{Name: user})
	if err != nil {
		log.Fatal(err)
	}
	var amqpChannel *amqp.Channel
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		amqpConn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			continue
		}
		amqpChannel, err = amqpConn.Channel()
		if err != nil {
			continue
		}
		defer amqpChannel.Close()
		break
	}
	if err != nil {
		log.Fatalf("amqp dial: %v", err)
		return
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("stream finished")
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", resp.Text)
		for i := 0; ; i++ {
			if resp.NeedsAnswer {
				var input string
				players, err := client.GetPlayersList(context.Background(),
					&base.Room{
						RoomId: resp.RoomId,
					})
				if i != 0 {
					if role != "bot" {
						fmt.Scanln(&input)
					} else {
						input = fmt.Sprintf("%d", rand.Int31()%int32(len(players.Players)))
					}
				}
				if input == "GetPlayers" || i == 0 {
					if err != nil {
						log.Fatal(err)
					}
					for _, value := range players.Players {
						fmt.Println(value.Name)
					}
					continue
				}
				if len(input) != 0 && input[0] == '!' {
					err = amqpChannel.Publish(
						"",
						fmt.Sprintf("room-%v", resp.RoomId),
						false,
						false,
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        []byte(user + "-ALLCHAT-" + input[1:]),
						},
					)
					if err != nil {
						log.Println(err)
					}
					continue
				}
				if len(input) != 0 && input[0] == '?' {
					err = amqpChannel.Publish(
						"",
						fmt.Sprintf("room-%v", resp.RoomId),
						false,
						false,
						amqp.Publishing{
							ContentType: "text/plain",
							Body:        []byte(user + "-MAFIACHAT-" + input[1:]),
						},
					)
					if err != nil {
						log.Println(err)
					}
					continue
				}
				val, err := strconv.ParseUint(input, 10, 32)
				if err != nil {
					log.Fatal(err)
					return
				}
				resp.SelectedTargetId = val
				client.ResponseEvent(context.Background(), resp)
			}
			break
		}
	}
	log.Printf("Finish")
}
