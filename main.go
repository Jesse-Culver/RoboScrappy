package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// Variables used throughout
var (
	Token string
	log   = logrus.New()
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Begin log setup
	file, err := os.OpenFile("Discord.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Create status message
	dg.AddHandler(func(dg *discordgo.Session, ready *discordgo.Ready) {
		err = dg.UpdateStatus(0, "Hello UNT!")
		if err != nil {
			fmt.Println("Error attempting to set status")
		}
	})

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	// Ignore all messages created by bots
	if message.Author.Bot {
		return
	}

	// Check if we get an error finding channel, guild, or member
	// Thank you to clinet for snippet https://github.com/clinet/clinet/blob/master/messages.go
	channel, err := session.State.Channel(message.ChannelID)
	if err != nil {
		return //Error finding the channel
	}
	guild, err := session.State.Guild(channel.GuildID)
	if err != nil {
		return //Error finding the guild
	}
	content := message.Content
	if content == "" {
		return //The message was empty
	}
	member, err := session.GuildMember(guild.ID, message.Author.ID)
	if err != nil {
		return //Error finding the guild member
	}

	fmt.Println("Message from " + member.User.Username + " received: " + message.Content)
	log.WithFields(logrus.Fields{
		"User":    member.User.Username,
		"UserID":  member.User.ID,
		"Channel": message.ChannelID,
		"Message": message.Content,
	}).Info("User Message")

	// if message.Content == "$Ping" {
	// 	session.ChannelMessageSend(message.ChannelID, "Pong!")
	// }
}
