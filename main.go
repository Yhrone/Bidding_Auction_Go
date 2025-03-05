package main

import (
	"fmt"
	"log"

	"github/Yhrone/go_Bidding_Auction/auction"
	"github/Yhrone/go_Bidding_Auction/bidding"
)

func main() {
	biddingPorts := []string{"8001", "8002", "8003"}

	log.Println("Starting ad auction system...")

	for _, port := range biddingPorts {
		go func(p string) {
			log.Printf("Starting bidding service on port %s", p)
			router := bidding.SetupRouter(p)
			addr := fmt.Sprintf(":%s", p)
			if err := router.Run(addr); err != nil {
				log.Printf("Bidding service on port %s failed: %v", p, err)
			}
		}(port)
	}

	log.Println("Starting auction service on port 8000")
	auctionRouter := auction.SetupRouter()
	if err := auctionRouter.Run(":8000"); err != nil {
		log.Fatalf("Auction service failed: %v", err)
	}
}
