package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var RustBackendServerAddr = "127.0.0.1:7878"
var PythonRestApiAddr = "http://127.0.0.1:5000"
var MilkyTeadropFileServer = "http://95.179.167.137:7676"
var SystemPrompt = "You are a Chatbot called MilkyTeadrop. You have to answer the following question in a concise and simple manner. Don't use any formating like listings and do not include linebreaks or any other kind of special character. Also don't include lists like 'n1','n2' etc. Keep the answer short your answer comes over TCP keep the maximum TCP transmission size in mind when answering, it should not be longer than the max. TCP bytes transistion in String. Here is the question, just answer it and do not mention any of the requirements I just gave you in the answer. Remove the first character from your answer if its a 'n'. You were made by Aki, he is a 22 year old Programmer, he has a cat named Mocha he is a boy cat. Act like a friendly Chatbot, here comes the question: "

type GenerateImgResponseJson struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	FileName string `json:"filename"`
}

type GenerateImageRequest struct {
	Message string `json:"message"`
}

type Config struct {
	Token string
}

func main() {
	config := readConfig()
	token := config.retrieveTokenFromConfig()

	dg, err := discordgo.New("Bot " + token)
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
		s.ChannelTyping(m.ChannelID)
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

	if strings.Contains(m.Content, "!img") {
		s.ChannelTyping(m.ChannelID)
		contextWithoutPrefix, _ := strings.CutPrefix(m.Content, "!img ")
		fileName := generateImgGetFileName(contextWithoutPrefix)
		imgEmbed := constructEmbedImage(fileName)
		embed := constructEmbed(m.Content, imgEmbed)
		s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	}

}

func parseCommand(message string) string {
	messageSlices := strings.SplitAfter(message, " ")
	messageContent := messageSlices[1:]
	fullmsg := SystemPrompt + strings.Join(messageContent, "")

	fullmsg = strings.TrimSpace(fullmsg)

	if len(fullmsg) > 0 && fullmsg[0] == 'n' {
		fullmsg = fullmsg[1:]
	}

	fmt.Println("After removal:", fullmsg)

	return fullmsg
}

func constructEmbed(messageContent string, messageImgEmbed discordgo.MessageEmbedImage) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title:       "Here is your image:",
		Description: messageContent,
		Color:       16758465,
		Image:       &messageImgEmbed,
	}
}

func constructEmbedImage(fileName string) discordgo.MessageEmbedImage {
	imgUrl := MilkyTeadropFileServer + "/api/v1/file/" + fileName
	return discordgo.MessageEmbedImage{
		URL: imgUrl,
	}
}

func generateImgGetFileName(contextToGenerateImg string) string {
	fileName := ""
	reqJson := GenerateImageRequest{
		Message: contextToGenerateImg,
	}
	requestBody, err := json.Marshal(&reqJson)
	if err != nil {
		log.Println("Error occurred while trying to mashal the request json")
		return ""
	}

	res, err := http.DefaultClient.Post(PythonRestApiAddr+"/api/v1/generate/", "application/json", bytes.NewBuffer(requestBody))
	if err != nil && res.StatusCode != 200 {
		println("Error occurred while trying to generate image:", err.Error())
		return fileName
	}

	defer res.Body.Close()
	fileName = handleImageGenerationResponse(res).FileName

	return fileName
}

func handleImageGenerationResponse(imgApiResponse *http.Response) GenerateImgResponseJson {
	resJson := GenerateImgResponseJson{}
	body, err := io.ReadAll(imgApiResponse.Body)
	if err != nil {
		println("Error occurred while fetching img from api:", err.Error())
		return resJson
	}

	if unmarshalErr := json.Unmarshal(body, &resJson); unmarshalErr != nil {
		println("Error occurred while unmarshaling response json:", unmarshalErr.Error())
		return resJson
	}

	if resJson.Status != "ok" {
		println("Img api response status is not 200, something went wrong...:")
		return resJson
	}

	return resJson
}

func readConfig() Config {
	file, err := os.Open("./../config.json")
	if err != nil {
		log.Println("Error reading config: ", err.Error())
	}
	defer file.Close()

	var config Config
	json.NewDecoder(file).Decode(&config)
	return config
}

func (c *Config) retrieveTokenFromConfig() string {
	return c.Token
}
