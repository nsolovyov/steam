package main

import (
	"log"
	"os"
	"time"

	"github.com/asamy45/steam"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	timeTip, err := steam.GetTimeTip()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Time tip: %#v\n", timeTip)

	timeDiff := time.Duration(timeTip.Time - time.Now().Unix())
	session := steam.Session{}
	if err := session.Login(os.Getenv("steamAccount"), os.Getenv("steamPassword"), os.Getenv("steamSharedSecret"), timeDiff); err != nil {
		log.Fatal(err)
	}
	log.Print("Login successful")

	key, err := session.GetWebAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Key: ", key)

	sent, _, err := session.GetTradeOffers(steam.TradeFilterSentOffers|steam.TradeFilterRecvOffers, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	var receiptID uint64
	for k := range sent {
		offer := sent[k]
		var sid steam.SteamID
		sid.Parse(offer.Partner, steam.AccountInstanceDesktop, steam.AccountTypeIndividual, steam.UniversePublic)

		if receiptID == 0 && len(offer.ReceiveItems) != 0 && offer.State == steam.TradeStateAccepted {
			receiptID = offer.ReceiptID
		}

		log.Printf("Offer id: %d, Receipt ID: %d", offer.ID, offer.ReceiptID)
		log.Printf("Offer partner SteamID 64: %d", uint64(sid))
	}

	items, err := session.GetTradeReceivedItems(receiptID)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items {
		log.Printf("New asset id: %d", item.AssetID)
	}

	log.Println("Bye!")
}