package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	to, domain := os.Getenv("DONOR_EMAIL_TO"), os.Getenv("SENDING_DOMAIN")
	if to == "" || domain == "" {
		panic("DONOR_EMAIL_TO and SENDING_DOMAIN are required")
	}
	c, err := NewClient()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	if err := onboardDomain(ctx, c, domain); err != nil {
		panic(err)
	}
	message, err := c.sendReceipt(ctx, to, "DON-1001", 2500)
	if err != nil {
		panic(err)
	}
	fmt.Println("receipt sent:", message.MessageID)
}
