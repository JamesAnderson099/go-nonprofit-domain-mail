package main

import "fmt"

func receiptText(receiptID string, amountCents int) string {
	if amountCents <= 0 {
		return fmt.Sprintf("Receipt %s: no charge recorded.", receiptID)
	}
	return fmt.Sprintf("Receipt %s: donation received, amount %d cents.", receiptID, amountCents)
}
