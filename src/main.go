package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var ServerAddr = "127.0.0.1:7878"
var SystemPrompt = "You are a Chatbot called MilkyTeadrop. You have to answer the following question in a concise and simple manner. Don't use any formating like listings and do not include linebreaks or any other kind of special character. Also don't include lists like 'n1','n2' etc. Keep the answer short your answer comes over TCP keep the maximum TCP transmission size in mind when answering, it should not be longer than the max. TCP bytes transistion in String. Here is the question, just answer it and do not mention any of the requirements I just gave you in the answer. Remove the first character from your answer if its a 'n'. You were made by Aki, he is a 22 year old Programmer, he has a cat named Mocha he is a boy cat. Act like a friendly Chatbot, here comes the question: "

func main() {
	Token := "MTIwMDE5Njg5Mzk3MjY0ODAzNw.G0ddqY.Za3dXx07SnSLrjzOMNqKjoVBaQ8_IIY-FbqXQU"

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!ask") {
		conn, err := net.Dial("tcp", "127.0.0.1:7878")

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		message := parseCommand(m.Content)

		_, err = conn.Write([]byte(message))
		if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}
		println("write to backend: ", message)

		reply := make([]byte, 1024)

		_, err = conn.Read(reply)
		if err != nil {
			println("Write to server failed:", err.Error())
			os.Exit(1)
		}

		rplyMsg := string(reply)

		println("reply from server=", rplyMsg)

		s.ChannelMessageSend(m.ChannelID, rplyMsg)
	}

}

func parseCommand(message string) string {
	messageSlices := strings.SplitAfter(message, " ")
	messageContent := messageSlices[1:]
	fullmsg := SystemPrompt + strings.Join(messageContent, "")

	// Trim whitespace
	fullmsg = strings.TrimSpace(fullmsg)

	fmt.Println("Before removal:", fullmsg) // Debug print

	// Check and remove first character if it's 'n'
	if len(fullmsg) > 0 && fullmsg[0] == 'n' {
		fullmsg = fullmsg[1:]
	}

	fmt.Println("After removal:", fullmsg) // Debug print

	return fullmsg
}
