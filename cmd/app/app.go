package main

import (
	"encoding/json"
	"fmt"
	"heydp/TgiInterface/clients"
	"heydp/TgiInterface/internal/redis"
	"heydp/TgiInterface/internal/textstreaming"
	"log"
	"net/http"
	"os"
)

func main() {
	os.Exit(mainerr())
}

func mainerr() int {
	errChan := make(chan error)
	go func(errchan chan error) {
		err := Run()
		if err != nil {
			errChan <- err
		}
	}(errChan)

	select {
	case err := <-errChan:
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
}

func Run() error {
	var err error
	fmt.Println("starting the text streaming app")

	srv := NewWebServer()

	env := GetAppEnv()
	configPath := fmt.Sprintf("configs/%s.json", env)

	file, err := os.ReadFile(configPath)
	if err != nil {
		errMsg := fmt.Sprintf("error in reading configFile, err - %v", err.Error())
		fmt.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	var appConfig AppConfig
	err = json.Unmarshal(file, &appConfig)
	if err != nil {
		errMsg := fmt.Sprintf("error in unmarshalling configFile, err - %v", err.Error())
		fmt.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	redisCreds := redis.RedisCreds{
		Host:     appConfig.Redis.Host,
		Port:     appConfig.Redis.Port,
		UserName: appConfig.Redis.UserName,
		Password: appConfig.Redis.Password,
		Database: 0,
	}

	redisClient, err := redisCreds.GiveRedisClient()
	if err != nil {
		log.Printf("failed to connect to redis, err - %v\n", err.Error())
		return err
	}

	var tgiClients = make(map[string]*clients.Client)
	for _, val := range appConfig.TgiServices {
		curClient := clients.NewClient(&http.Client{}, val.Scheme, val.Host, val.Port, val.Path, val.Token)
		tgiClients[val.ServiceId] = curClient
	}

	ctrls, ictrls := setUpRoutes(redisClient, tgiClients)

	var serverResources ServerResources
	serverResources.ctrls = ctrls
	serverResources.itcrls = ictrls

	srv.InitRouter(&serverResources)
	fmt.Println("initializing servers")

	srv.initialized = true
	err = srv.Start(appConfig)
	if err != nil {
		return err
	}

	return nil
}

func setUpRoutes(redisClient *redis.RedisDb, tgiClients map[string]*clients.Client) ([]Controllers, []Controllers) {

	var controller []Controllers
	var icontroller []Controllers

	textStreamingManager := textstreaming.NewManager(redisClient, tgiClients)
	textStreamingController := textstreaming.NewController(textStreamingManager)
	controller = append(controller, textStreamingController)

	return controller, icontroller
}
