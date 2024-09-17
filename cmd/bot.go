package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/thenanox/ficsit-discord-bot/cmd/slashcommands/ping"
	"github.com/thenanox/ficsit-discord-bot/cmd/slashcommands/pioneers"
	"github.com/thenanox/ficsit-discord-bot/internal/satisfactory"
	"github.com/zekrotja/ken"
)

var numberOfPlayers int = 0

func Execute(token string) error {
	var terminationSignalChannel = make(chan os.Signal, 1)
	signal.Notify(terminationSignalChannel, os.Interrupt, syscall.SIGTERM)

	// cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	waitGroup := &sync.WaitGroup{}

	startHealthServer("0.0.0.0:8080")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		cancel()
		return err
	}
	defer session.Close()

	session.Identify.Intents = discordgo.IntentsAll

	err = registerSlashCommands(ctx, session)
	if err != nil {
		cancel()
		return err
	}

	routineManager(ctx, session, waitGroup)
	if err != nil {
		cancel()
		return err
	}

	err = session.Open()
	if err != nil {
		cancel()
		return err
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	for {
		select {
		case <-terminationSignalChannel:
			cancel()
			waitGroup.Wait()
			close(terminationSignalChannel)
			os.Exit(0)
		}
	}
}

func startHealthServer(addr string) *http.Server {
	// HTTP Gateway server
	router := mux.NewRouter()
	router.HandleFunc("/api/health", HandleHealth)

	srv := &http.Server{
		Addr:              addr,
		WriteTimeout:      time.Second * 15,
		ReadTimeout:       time.Second * 15,
		IdleTimeout:       time.Second * 60,
		ReadHeaderTimeout: time.Second * 15,
		Handler:           handlers.LoggingHandler(os.Stdout, router),
	}
	go func() {
		fmt.Printf("starting HTTP main server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("error starting HTTP main server: %v", err))
		}
	}()

	return srv
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	return
}

func registerSlashCommands(_ context.Context, session *discordgo.Session) error {
	k, err := ken.New(session)
	if err != nil {
		return err
	}
	err = k.RegisterCommands(
		new(ping.PingCommand),
		new(pioneers.PioneersCommand),
	)
	if err != nil {
		return err
	}
	defer k.Unregister()

	return nil
}

func routineManager(ctx context.Context, session *discordgo.Session, waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	go spawnCheckSatisfactoryServer(ctx, session, waitGroup)
}

func spawnCheckSatisfactoryServer(ctx context.Context, session *discordgo.Session, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	defer func() {
		if x := recover(); x != nil {
			waitGroup.Add(1)
			go spawnCheckSatisfactoryServer(ctx, session, waitGroup)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := checkSatisfactoryServer(ctx, session)
			if err != nil {
				return
			}
			t := time.Duration(60) * time.Second
			err = sleep(ctx, t)
			if err != nil {
				return
			}
		}
	}
}

func checkSatisfactoryServer(_ context.Context, session *discordgo.Session) error {
	response, err := satisfactory.QueryServerState()
	if err != nil {
		return err
	}
	err = session.UpdateCustomStatus(fmt.Sprintf("%d/%d pioneers", response.Data.ServerGameState.NumConnectedPlayers, response.Data.ServerGameState.PlayerLimit))
	if err != nil {
		return err
	}
	if numberOfPlayers < response.Data.ServerGameState.NumConnectedPlayers {
		numberOfPlayers = response.Data.ServerGameState.NumConnectedPlayers
		channel := os.Getenv("DISCORD_CHANNEL")
		role := os.Getenv("DISCORD_ROLE")
		_, err = session.ChannelMessageSend(channel, fmt.Sprintf("%s A new pioneer entered in the server. %d/%d pioneers in MASSAGE-2(A-B)b", role, response.Data.ServerGameState.NumConnectedPlayers, response.Data.ServerGameState.PlayerLimit))
		if err != nil {
			return err
		}
	} else if numberOfPlayers > response.Data.ServerGameState.NumConnectedPlayers {
		numberOfPlayers = response.Data.ServerGameState.NumConnectedPlayers
		channel := os.Getenv("DISCORD_CHANNEL")
		role := os.Getenv("DISCORD_ROLE")
		_, err = session.ChannelMessageSend(channel, fmt.Sprintf("%s A pioneer exited the server. %d/%d pioneers in MASSAGE-2(A-B)b", role, response.Data.ServerGameState.NumConnectedPlayers, response.Data.ServerGameState.PlayerLimit))
		if err != nil {
			return err
		}
	}
	return nil
}

func sleep(ctx context.Context, t time.Duration) error {
	timeoutchan := make(chan bool)
	go func() {
		<-time.After(t)
		timeoutchan <- true
	}()

	select {
	case <-timeoutchan:
		break
	case <-ctx.Done():
		return errors.New("terminated sleep due to a context cancellation")
	}
	return nil
}
