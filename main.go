package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type response struct {
	Results results `json:"results"`
}

type results struct {
	Shop []shop `json:"shop"`
}

type shop struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func getRestInfo(lat string, lng string) string {
	apikey := "e937fa83d59190c6"
	url := fmt.Sprintf("https://webservice.recruit.co.jp/hotpepper/gourmet/v1/?key=%s&format=json&lat=%s&lng=%s", apikey, lat, lng)
	//リクエストしてボディを取得
	resp, err := http.Get(url) //Get Method issues a Get to the specified URL（httpリクエストのに対するレスポンスを返す）
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()                //http.Getによってrespには常にnil以外のresp.Bodyが含まれる。　resp.Bodyからの読み取りが完了したら閉じる。
	body, err := ioutil.ReadAll(resp.Body) //戻り値 []byte resp.Body >> レスポンスの中身
	if err != nil {
		log.Fatal(err)
	}
	var data response
	if err := json.Unmarshal(body, &data); err != nil { //json>>構造体
		log.Fatal(err)
	}

	info := ""
	for _, shop := range data.Results.Shop {
		info += shop.Name + "\n" + shop.Address + "\n\n"
	}
	return info
}

func sendRestInfo(bot *linebot.Client, event *linebot.Event) {
	msg := event.Message.(*linebot.LocationMessage)

	lat := strconv.FormatFloat(msg.Latitude, 'f', 2, 64)  //緯度取得 f=実数表現実現  第3引数=桁数 第4引数=bitSize=64
	lng := strconv.FormatFloat(msg.Longitude, 'f', 2, 64) //経度取得

	//replyMsg := fmt.Sprintf("緯度:%s\n経度:%s", lat, lng)
	replyMsg := getRestInfo(lat, lng)
	_, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMsg)).Do()
	if err != nil {
		log.Print(err)
	}
}

func main() {
	bot, err := linebot.New(
		"e0fca7fa7bd74ac93fd14c8fcec2c12e", //発行されたチャネルシークレット
		"JX8eE8e8WzXBajYpHzYIe4W5mxeCLBbcYdcLBSZkbX00rBikrPkL272gpTfYg8UTQ3YJj0sfUoLDvc6GHHpahIPDG67WYFxUSo5C1Geq7mFNPGf2zuVdI8FaS2NoZSfmyveSwmSAUUaXm2a0tTi6rgdB04t89/1O/w1cDnyilFU=", //発行されたチャネルアクセストークン
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	parrot := func(w http.ResponseWriter, req *http.Request) {
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
				case *linebot.TextMessage: //文字列ならば
					replyMessage := message.Text //Lineで送信された内容
					_, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
					if err != nil {
						log.Print(err)
					}
				case *linebot.StickerMessage: //スタンプの場合
					replyMessage := fmt.Sprintf(
						"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.LocationMessage:
					sendRestInfo(bot, event)
				}
			}
		}
	}
	http.HandleFunc("/callback", parrot)
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":2022", nil); err != nil {
		log.Fatal(err)
	}
}
