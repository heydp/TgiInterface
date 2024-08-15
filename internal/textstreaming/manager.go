package textstreaming

import (
	"heydp/TgiInterface/clients"
	"heydp/TgiInterface/internal/redis"
	"io"
	"log"
)

const (
	DOWN = "down"
)

type Manager struct {
	redisClient *redis.RedisDb
	tgiClients  map[string]*clients.Client
}

func NewManager(redisClient *redis.RedisDb, tgiClients map[string]*clients.Client) *Manager {
	return &Manager{
		redisClient: redisClient,
		tgiClients:  tgiClients,
	}
}

func (m *Manager) HealthCheck() {
	for key, val := range m.tgiClients {
		resp, err := val.HealthCheck()
		if err != nil {
			m.redisClient.Update(key, DOWN)
		} else {
			m.redisClient.Update(key, *resp)
		}
	}
}

func (m *Manager) TextStreaming(req Request) *io.PipeReader {
	// var respStr string
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		flag := false
		for key, val := range m.tgiClients {
			healthCheck, err := m.redisClient.Find(key)
			if err == nil && *healthCheck == DOWN {
				continue
			}
			resp, err := val.GenerateTextResponse(req.Text)
			if err != nil {
				log.Println("the err is : ", err.Error())
				m.redisClient.Update(key, DOWN)
				continue
			} else {
				writer.Write([]byte(*resp))
				flag = true
				break
			}
		}
		if flag == false {
			writer.Write([]byte("currently unable to find response"))
		}
	}()

	return reader

	// return &Response{respStr}, nil
}
