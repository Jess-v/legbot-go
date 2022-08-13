package main

import (
	"bufio"
	"errors"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
)

var folderPath string = "/app/users/"

func main() {
	configureLogs()
	token := os.Getenv("DISCORD_API_TOKEN")
	log.Info("Starting Discord bot with API token: " + token)

	discordBot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Error("Failed creating Discord session: " + err.Error())
	}

	discordBot.Identify.Intents |= discordgo.IntentMessageContent
	discordBot.Identify.Intents |= discordgo.IntentsGuildMessages
	discordBot.AddHandler(messageCreate)

	err = discordBot.Open()
	if err != nil {
		log.Error("Failed opening connection:" + err.Error())
	}

	log.Info("Bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discordBot.Close()
}

func configureLogs() {
	log.SetFormatter(&ecslogrus.Formatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true)
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	filePath := folderPath + message.Author.ID + ".txt"

	messageLower := strings.ToLower(message.Content)

	if strings.Contains(messageLower, "left leg next week") {
		wasWritten := writeConfig(filePath, "left leg")
		if wasWritten == false {
			log.Error("Unable to write config for " + message.Author.ID)
			session.ChannelMessageSend(message.ChannelID, "Oh no! Jess messed something up with me and I wasn't able to save this 😔")
			return
		} else {
			log.Info("Config written for " + message.Author.ID)
			session.ChannelMessageSend(message.ChannelID, "I got ya")
			return
		}
	} else if strings.Contains(messageLower, "right leg next week") {
		wasWritten := writeConfig(filePath, "right leg")
		if wasWritten == false {
			log.Error("Unable to write config for " + message.Author.ID)
			session.ChannelMessageSend(message.ChannelID, "Oh no! Jess messed something up with me and I wasn't able to save this 😔")
			return
		} else {
			log.Info("Config written for " + message.Author.ID)
			session.ChannelMessageSend(message.ChannelID, "I got ya")
		}
	} else if strings.Contains(messageLower, "which leg") {
		if idExists(filePath) {
			response := getConfig(filePath)
			session.ChannelMessageSend(message.ChannelID, response)
			log.Info("Config retrieved for " + message.Author.ID)
			return
		} else {
			session.ChannelMessageSend(message.ChannelID, "Sorry! I don't seem to have that info for you.")
			log.Warn("Config file not found for " + message.Author.ID)
			return
		}
	} else if strings.Contains(messageLower, "good bot") {
		session.ChannelMessageSend(message.ChannelID, "😊")
	}
	return
}

func getConfig(filePath string) string {
	var config string
	file, err := os.Open(filePath)
	if err != nil {
		log.Error("Failure opening file " + err.Error())
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Error("Failure closing file " + err.Error())
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		config = scanner.Text()
	}
	return config
}

func idExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		log.Error("Cannot determine if" + filePath + "exists " + err.Error())
		return false
	}
}

func writeConfig(filePath string, content string) bool {
	file, err := os.Create(filePath)
	if err != nil {
		log.Error("Cannot create or open file. " + err.Error())
		return false
	}
	_, err = file.WriteString(content)
	if err != nil {
		log.Error("Cannot write to file." + err.Error())
		return false
	}
	err = file.Close()
	if err != nil {
		log.Error("Failure closing file." + err.Error())
		return false
	}
	return true
}
