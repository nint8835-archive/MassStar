package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
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

func onMessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	currentUser, _ := s.User("@me")
	log.Debug().
		Str("bot", currentUser.String()).
		Str("message_id", m.MessageID).
		Str("emoji", m.Emoji.Name).
		Msg("Reaction added")

	if m.Emoji.Name == "🔺" {
		log.Info().
			Str("bot", currentUser.String()).
			Str("message_id", m.MessageID).
			Msg("Got trigger emoji")
		s.MessageReactionAdd(m.ChannelID, m.MessageID, "⭐")
	}
}

func onMessageReactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	currentUser, _ := s.User("@me")
	log.Debug().
		Str("bot", currentUser.String()).
		Str("message_id", m.MessageID).
		Str("emoji", m.Emoji.Name).
		Msg("Reaction removed")

	if m.Emoji.Name == "🔺" {
		log.Info().
			Str("bot", currentUser.String()).
			Str("message_id", m.MessageID).
			Msg("Removed trigger emoji")
		s.MessageReactionRemove(m.ChannelID, m.MessageID, "⭐", "@me")
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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
		client.AddHandler(onMessageReactionAdd)
		client.AddHandler(onMessageReactionRemove)
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
