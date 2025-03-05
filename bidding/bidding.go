package bidding

import (
	"github/Yhrone/go_Bidding_Auction/models"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRouter(port string) *gin.Engine {
	r := gin.Default()

	r.POST("/bid", func(c *gin.Context) {
		var req models.AdRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Bidding service failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if rand.Float64() < 0.3 {
			c.Status(http.StatusNoContent)
		}
		time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)

		bidPrice := rand.Float64()*9.9 + 0.1
		adObject := models.AdObject{
			AdID:     "ad_" + strconv.Itoa(rand.Intn(9000)+1000), // Random AdID (e.g., ad_1234)
			BidPrice: float64(int(bidPrice*100)) / 100,           // Round to 2 decimal places
		}
		c.JSON(http.StatusOK, models.BidResponse{
			Status:   http.StatusOK,
			AdObject: &adObject,
		})
	})
	return r
}
