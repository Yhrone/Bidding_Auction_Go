package auction

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github/Yhrone/go_Bidding_Auction/models"

	"github.com/gin-gonic/gin"
)

var biddingServices = []string{
	"http://localhost:8001/bid",
	"http://localhost:8002/bid",
	"http://localhost:8003/bid",
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/auction", func(c *gin.Context) {

		var req models.AdRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Printf("Auction service failed to parse request: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Channel to collect bids with a timeout
		bidChan := make(chan *models.AdObject, len(biddingServices))
		var wg sync.WaitGroup
		client := &http.Client{Timeout: 200 * time.Millisecond}
		body, _ := json.Marshal(req)

		// Start bidding requests concurrently
		for _, url := range biddingServices {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
				if err != nil {
					log.Printf("Bid from %s failed or timed out: %v", url, err)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode == http.StatusOK {
					var bidResp models.BidResponse
					if err := json.NewDecoder(resp.Body).Decode(&bidResp); err != nil {
						log.Printf("Failed to decode bid from %s: %v", url, err)
					} else if bidResp.AdObject != nil {
						bidChan <- bidResp.AdObject
					}
				}
			}(url)
		}

		// Close channel after all bids are done or timeout occurs
		go func() {
			wg.Wait()
			close(bidChan)
		}()

		// Collect bids with a 200ms deadline
		var highestBid *models.AdObject
		timeout := time.After(200 * time.Millisecond)

		for {
			select {
			case bid, ok := <-bidChan:
				if !ok { // Channel closed, all bids collected
					goto done
				}
				if highestBid == nil || bid.BidPrice > highestBid.BidPrice {
					highestBid = bid
				}
			case <-timeout:
				log.Printf("Auction timed out after 200ms, using bids received so far")
				goto done
			}
		}

	done:
		if highestBid == nil {
			log.Printf("No valid bids received within 200ms")
			c.JSON(http.StatusNoContent, nil)
			return
		}

		log.Printf("Winner selected: AdID=%s, BidPrice=%.2f", highestBid.AdID, highestBid.BidPrice)
		c.JSON(http.StatusOK, models.BidResponse{
			Status:   http.StatusOK,
			AdObject: highestBid,
		})
	})

	return r
}
