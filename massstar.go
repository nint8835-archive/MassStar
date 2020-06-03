package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// BotToken represents an individual token to be used
type BotToken struct {
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}

// ClientObj represents an individual instantiated client and any associated data
type ClientObj struct {
	Client *discordgo.Session
	Token  BotToken
}

var clients []ClientObj
var processedMessages = map[string]bool{}
var processedMessagesLock sync.Mutex
var starRegex *regexp.Regexp

func onMessageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {
	processedMessagesLock.Lock()
	defer processedMessagesLock.Unlock()

	channel, _ := s.Channel(m.ChannelID)
	currentUser, _ := s.User("@me")
	log.Info().
		Str("channel", channel.Name).
		Str("content", m.Content).
		Str("bot", currentUser.String()).
		Str("message_id", m.ID).
		Bool("already_seen", processedMessages[m.ID]).
		Msg("Message created")

	if processedMessages[m.ID] {
		return
	}
	processedMessages[m.ID] = true

	// match := starRegex.FindStringSubmatch(m.Content)
	// if len(match) == 1 {

	// }
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	starRegex = regexp.MustCompile("^\\.s (\\d+)$")

	var tokens []BotToken

	log.Debug().Msg("Loading tokens")
	tokenFile, _ := os.Open("tokens.json")
	defer tokenFile.Close()
	tokenBytes, _ := ioutil.ReadAll(tokenFile)
	json.Unmarshal(tokenBytes, &tokens)
	log.Debug().Interface("tokens", tokens).Msg("Tokens loaded")

	log.Debug().Msg("Initializing clients")
	for _, token := range tokens {
		log.Debug().Str("name", token.Name).Msg("Initializing client")
		client, _ := discordgo.New("Bot " + token.Token)
		client.AddHandler(onMessageCreated)
		_ = client.Open()
		clients = append(clients, ClientObj{client, token})
	}
	log.Debug().Msg("Clients initialized")

	log.Info().Msg("MassStar running. Press Ctrl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Debug().Msg("Cleaning up Discord sessions")

	for _, clientObj := range clients {
		log.Debug().Str("name", clientObj.Token.Name).Msg("Cleaning up session")
		clientObj.Client.Close()
	}
}
