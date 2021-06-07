package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	stock_key := os.Getenv("STOCK_KEY")
	if stock_key == "" {
		log.Fatal("$STOCK_KEY must be set")
	}
	r := gin.Default()
	r.GET("/hs", func(c *gin.Context) {
		gid := c.Query("gid")
		hsStock := HsStock{}
		if err := getJson("http://web.juhe.cn:8080/finance/stock/hs?gid="+
			gid+"&key="+stock_key, &hsStock); err != nil {
			c.String(http.StatusBadGateway, err.Error())
		} else {
			c.JSON(http.StatusOK, hsStock)
		}

	})
	r.Run(":" + port)
}
