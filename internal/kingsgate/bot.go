package kingsgate

import (
	"fmt"
	"kingsgate/internal/kingsgate/commands"
	"kingsgate/internal/kingsgate/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

var (
	session *discordgo.Session
)

// Run sets up the session and runs the bot
func Run() error {
	fmt.Println("Initialising bot base config")
	err := config.InitBaseConfig()
	if err != nil {
		return errors.Wrap(err, "Error loading bot config")
	}

	fmt.Println("Creating discord session")
	session, err = discordgo.New("Bot " + config.Config.Token)
	if err != nil {
		return errors.Wrap(err, "Error creating discord session")
	}

	session.AddHandler(ready)

	// Add commands to session
	fmt.Println("Registering commands to session")
	commands.Register(session)

	// Open the session's web-socket
	fmt.Println("Opening session")
	session.Open()

	fmt.Println("Bot running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	session.Close()

	return nil
}

func ready(session *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	session.UpdateStatus(0, "the bone xylophone")
}
