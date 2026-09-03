package main

import "testing"

func TestReceiptTextKeepsZeroDonationUncharged(t *testing.T) {
	if got := receiptText("R-17", 0); got != "Receipt R-17: no charge recorded." {
		t.Fatalf("got %q", got)
	}
	if got := receiptText("R-18", 2500); got != "Receipt R-18: donation received, amount 2500 cents." {
		t.Fatalf("got %q", got)
	}
}
