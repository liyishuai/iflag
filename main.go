package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// https://stackoverflow.com/a/31129967
var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

type HsStock struct {
	ResultCode, Reason string
	Result             []struct {
		Data struct {
			Gid, IncrePer, Increase, Name,
			TodayStartPri, YestodEndPri, NowPri, TodayMax, TodayMin,
			CompetitivePri, ReservePri, TraNumber, TraAmount,
			BuyOne, BuyOnePri,
			BuyTwo, BuyTwoPri,
			BuyThree, BuyThreePri,
			BuyFour, BuyFourPri,
			BuyFive, BuyFivePri,
			SellOne, SellOnePri,
			SellTwo, SellTwoPri,
			SellThree, SellThreePri,
			SellFour, SellFourPri,
			SellFive, SellFivePri,
			Date, Time string
		}
		DapanData struct {
			Dot, Name, NowPic, Rate, TraAmount, TraNumber string
		}
		GoPicture struct {
			MinUrl, DayUrl, WeekUrl, MonthUrl string
		}
	}
	ErrorCode int
}

var stock_key = os.Getenv("STOCK_KEY")

func getHsStock(gid string) (HsStock, error) {
	hsStock := HsStock{}
	err := getJson("http://web.juhe.cn:8080/finance/stock/hs?gid="+
		gid+"&key="+stock_key, &hsStock)
	return hsStock, err
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	tgPath := os.Getenv("TELEGRAM_PATH")
	if tgPath == "" {
		log.Fatal("$TELEGRAM_PATH must be set")
	}
	if stock_key == "" {
		log.Fatal("$STOCK_KEY must be set")
	}
	bot, err := tgbotapi.NewBotAPIWithClient(
		os.Getenv("TELEGRAM_APITOKEN"),
		myClient)
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)
	webhook := tgbotapi.NewWebhook("https://iflag.herokuapp.com/" + tgPath)
	webhook.MaxConnections = 100
	if _, err := bot.SetWebhook(webhook); err != nil {
		log.Fatal(err)
	}
	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	updates := bot.ListenForWebhook("/" + tgPath)
	go http.ListenAndServe(":"+port, nil)
	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		bot.Send(msg)
	}
}
