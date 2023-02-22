package webhook

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/lushenle/mmchatgpt/config"
	"github.com/lushenle/mmchatgpt/gpt"
	"github.com/mattermost/mattermost-server/v5/model"
)

type WhSvrParam struct {
	Port     int
	CertFile string
	KeyFile  string
}

type WebHookServer struct {
	Server *http.Server // http server
}

func (s *WebHookServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mmURL := config.GetMattermostURL()
	mmToken := config.GetMattermostToken()
	//botUsername := config.GetBotUsername()

	// Parse the request body as a JSON object.
	var data map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Failed to parse message:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Ensure that the file_ids field is initialized as an empty array.
	data["file_ids"] = &model.StringArray{}

	// Extract the message fields from the JSON object.
	post := &model.Post{}

	var trim string
	if triggerWord, ok := data["trigger_word"].(string); ok {
		trim = triggerWord
	}

	if text, ok := data["text"].(string); ok {
		post.Message = strings.TrimSpace(strings.TrimPrefix(text, trim))
	} else {
		log.Println("Failed to parse message: missing or invalid 'text' field")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if channelID, ok := data["channel_id"].(string); ok {
		post.ChannelId = channelID
	} else {
		log.Println("Failed to parse message: missing or invalid 'channel_id' field")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Only respond to messages that mention the bot.
	//if !strings.Contains(post.Message, "@"+botUsername) {
	//	w.WriteHeader(http.StatusOK)
	//	return
	//}

	// Send the message to the ChatGPT API and get the response.
	response, err := gpt.GenerateResponse(post.Message)
	if err != nil {
		log.Println("Failed to generate response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Post the response back to the channel.
	params := &model.Post{
		ChannelId: post.ChannelId,
		Message:   fmt.Sprintf("---\n%s\n---", response),
	}
	client := model.NewAPIv4Client(mmURL)
	client.SetOAuthToken(mmToken)
	createdPost, resp := client.CreatePost(params)
	if resp.Error != nil {
		log.Println("Failed to post message:", resp.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Posted response to channel %s: %s", createdPost.ChannelId, createdPost.Id)

	w.WriteHeader(http.StatusOK)
}
