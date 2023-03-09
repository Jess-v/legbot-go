package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/caarlos0/env"
	log "github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"

	"main/internal/pkg/database"
)

type environmentVariables struct {
	BotToken         string `env:"BOT_TOKEN,required"`
	LogLevel         string `env:"LOG_LEVEL" envDefault:"info"`
	DatabaseHost     string `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort     string `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseName     string `env:"DATABASE_NAME" envDefault:"legbot"`
	DatabaseUser     string `env:"DATABASE_USER" envDefault:"postgres"`
	DatabasePassword string `env:"DATABASE_PASSWORD" envDefault:"postgres"`
}

var db database.DatabaseConnection

func main() {
	envVars := &environmentVariables{}
	err := env.Parse(envVars)
	if err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	bot := initBot(envVars)

	err = bot.Open()
	if err != nil {
		log.Fatalf("Failed to open bot connection: %v", err)
	}

	defer bot.Close()

	log.Info("Bot is now running. Press CTRL-C to exit.")

	log.Info("Registering commands...")
	registerCommands(bot)

	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Info("Press CTRL-C to exit.")
	<-stop

}

func configureLogs(logLevel string) {
	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.SetFormatter(&ecslogrus.Formatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
}

func initBot(envVars *environmentVariables) *discordgo.Session {
	configureLogs(envVars.LogLevel)

	log.WithFields(log.Fields{
		"logLevel":     envVars.LogLevel,
		"databaseHost": envVars.DatabaseHost,
		"databasePort": envVars.DatabasePort,
		"databaseUser": envVars.DatabaseUser,
	}).Info("initializing bot...")

	db = database.DatabaseConnection{
		Host:     envVars.DatabaseHost,
		Port:     envVars.DatabasePort,
		Name:     envVars.DatabaseName,
		User:     envVars.DatabaseUser,
		Password: envVars.DatabasePassword,
	}

	log.Info("Initializing database...")
	err := db.InitDb()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Info("Database initialized!")

	discordBot, err := discordgo.New("Bot " + envVars.BotToken)

	return discordBot
}

var (
	integerOptionMinValue    = 1
	dmPermission             = false
	defaultMemberPermissions = discordgo.PermissionSendMessages | discordgo.PermissionEmbedLinks | discordgo.PermissionAddReactions

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "set",
			Description: "Set the location for your next HRT dosage",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "location",
					Description: "The location you want to set",
					Required:    true,
				},
			},
		},
		{
			Name:        "where",
			Description: "Get the location for your next HRT dosage",
		},
		{
			Name:        "praise",
			Description: "Praise the bot, tell it that it did a good job",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"set":    handleSet,
		"where":  handleWhere,
		"praise": handlePraise,
	}
)

func handleSet(s *discordgo.Session, i *discordgo.InteractionCreate) {
	location := i.ApplicationCommandData().Options[0].StringValue()
	userExists, err := db.CheckUser(i.Member.User.ID)
	if err != nil {
		log.Error(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There was an error checking if you're in the database :(",
				Flags:   64,
			},
		})
		return
	}
	if userExists {
		err := db.UpdateUser(i.Member.User.ID, location)
		if err != nil {
			log.Error(err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "There was an error updating your location",
					Flags:   64,
				},
			})
			return
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Your location has been updated to %s", location),
				Flags:   64,
			},
		})
		return
	} else {
		err := db.NewUser(i.Member.User.ID, location)
		if err != nil {
			log.Error(err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "There was an error adding you to the database :(",
					Flags:   64,
				},
			})
			return
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("I'll remember that for ya! Your next dose will be applied to: %s", location),
				Flags:   64,
			},
		})
	}
}

func handleWhere(s *discordgo.Session, i *discordgo.InteractionCreate) {
	location, err := db.GetLocation(i.Member.User.ID)
	if err != nil {
		log.Error(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There was an error getting your location :(",
				Flags:   64,
			},
		})
		return
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Apply your dose to the following location: %s", location),
			Flags:   64,
		},
	})
}

func handlePraise(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "😳",
			Flags:   64,
		},
	})
}

func registerCommands(s *discordgo.Session) {
	for _, command := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", command)
		if err != nil {
			log.Error(err)
		}
	}
}

func removeCommands(s *discordgo.Session) {
	for _, command := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", command.Name)
		if err != nil {
			log.Error(err)
		}
	}
}
