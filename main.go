package main

import (
	"fmt"
	"log"
	"net/http"
//	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	bot, err := linebot.New(
		//os.Getenv("CHANNEL_SECRET"),
		//os.Getenv("CHANNEL_TOKEN"),
        "e0fca7fa7bd74ac93fd14c8fcec2c12e",
        "JX8eE8e8WzXBajYpHzYIe4W5mxeCLBbcYdcLBSZkbX00rBikrPkL272gpTfYg8UTQ3YJj0sfUoLDvc6GHHpahIPDG67WYFxUSo5C1Geq7mFNPGf2zuVdI8FaS2NoZSfmyveSwmSAUUaXm2a0tTi6rgdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				//メッセージがテキスト形式の場合
				case *linebot.TextMessage:
					replyMessage := message.Text
					_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
					if err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
