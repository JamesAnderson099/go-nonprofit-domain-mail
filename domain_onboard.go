package main

import (
	"context"
	"fmt"
)

func onboardDomain(ctx context.Context, c *Client, domain string) error {
	if err := c.verifyDomain(ctx, domain); err != nil {
		return err
	}
	status, err := c.domainStatus(ctx, domain)
	if err != nil {
		return err
	}
	if status.Verification.Status == "" {
		return fmt.Errorf("domain verification status is empty")
	}
	fmt.Println("domain verification status:", status.Verification.Status)
	return nil
}
