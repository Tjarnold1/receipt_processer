package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode"
)

const dateLayout string = "2006-01-02"
const timeLayout string = "15:04"

// Writing custom unmarshalers because I want to throw parsing errors when we ingest receipts. Not when calculating points

type PurchaseDate struct {
	Date time.Time
}

func (pd *PurchaseDate) UnmarshalJSON(b []byte) error {
	parsedDate, err := time.Parse(dateLayout, strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}
	pd.Date = parsedDate
	return nil
}

type PurchaseTime struct {
	Time time.Time
}

func (pt *PurchaseTime) UnmarshalJSON(b []byte) error {
	parsedTime, err := time.Parse(timeLayout, strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}
	pt.Time = parsedTime
	return nil
}

type Item struct {
	ShortDescription string  `json:"shortDescription" binding:"required"`
	Price            float64 `json:"price,string" binding:"required"`
}

type Receipt struct {
	Retailer     string       `json:"retailer" binding:"required"`
	PurchaseDate PurchaseDate `json:"purchaseDate" binding:"required"`
	PurchaseTime PurchaseTime `json:"purchaseTime" binding:"required"`
	Items        []Item       `json:"items" binding:"required,dive"`
	Total        float64      `json:"total,string" binding:"required"`
}

var storage = make(map[uuid.UUID]Receipt)

func main() {
	router := gin.Default()
	router.POST("receipts/process", processReceipt)
	router.GET("receipts/:id/points", getReceiptPoints)

	router.Run("localhost:8080")
}

func processReceipt(c *gin.Context) {
	var receipt Receipt

	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	/*
		Considered processing the receipt here, and just storing the points scored since it would take up so much less
		memory. I decided to keep the entire receipt because in the future we may change the scoring rules and
		retroactively re-calculate points, or add some other new receipt related feature. We can always get rid of data
		later. We can't get it back.
	*/
	storage[id] = receipt
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getReceiptPoints(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing provided ID"})
	}
	receipt, ok := storage[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
	}

	points := receipt.CalculatePoints()

	c.JSON(http.StatusOK, gin.H{"points": points})
}

func (r *Receipt) CalculatePoints() int {
	pointsRewarded := 0
	pointsRewarded += calculateRetailerNamePoints(r.Retailer)
	pointsRewarded += calculateRoundTotalPoints(r.Total)
	pointsRewarded += calculateQuarterTotalPoints(r.Total)
	pointsRewarded += calculateItemPairPoints(r.Items)
	pointsRewarded += calculateItemNameLengthPoints(r.Items)
	pointsRewarded += calculateOddPurchaseDatePoints(r.PurchaseDate)
	pointsRewarded += calculateHappyHourPoints(r.PurchaseTime)
	return pointsRewarded
}

func calculateRetailerNamePoints(retailer string) int {
	points := 0
	for _, r := range strings.TrimSpace(retailer) {
		if unicode.IsDigit(r) || unicode.IsLetter(r) {
			points++
		}
	}
	return points
}

func calculateRoundTotalPoints(total float64) int {
	if total == math.Trunc(total) {
		return 50
	}
	return 0
}

func calculateQuarterTotalPoints(total float64) int {
	quotient := total / 0.25
	if quotient == math.Trunc(quotient) {
		return 25
	}
	return 0
}

func calculateItemPairPoints(items []Item) int {
	return (len(items) / 2) * 5
}

func calculateItemNameLengthPoints(items []Item) int {
	points := 0
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}
	return points
}

func calculateOddPurchaseDatePoints(date PurchaseDate) int {
	if date.Date.Day()%2 == 1 {
		return 6
	}
	return 0
}

func calculateHappyHourPoints(time PurchaseTime) int {
	if time.Time.Hour() >= 14 && time.Time.Hour() < 16 {
		return 10
	}
	return 0
}
