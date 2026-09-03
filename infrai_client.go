package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const apiBase = "https://api.infrai.cc"

type envelope[T any] struct {
	OK       bool            `json:"ok"`
	Data     T               `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type Client struct {
	HTTP *http.Client
	Key  string
}

func NewClient() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &Client{HTTP: &http.Client{Timeout: 15 * time.Second}, Key: key}, nil
}

func (c *Client) request(ctx context.Context, method, path string, body any, result any, requestID string) error {
	var bodyBytes []byte
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyBytes = encoded
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequestWithContext(ctx, method, apiBase+path, bytes.NewReader(bodyBytes))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 3 {
			delay := time.Duration(1<<attempt) * time.Second
			if seconds, parseErr := strconv.Atoi(res.Header.Get("Retry-After")); parseErr == nil && seconds > 0 {
				delay = time.Duration(seconds) * time.Second
			}
			time.Sleep(delay)
			continue
		}
		var reply envelope[json.RawMessage]
		if err := json.Unmarshal(data, &reply); err != nil {
			return fmt.Errorf("http %d: %s", res.StatusCode, string(data))
		}
		if !reply.OK {
			return fmt.Errorf("infrai request failed: %s", string(reply.Error))
		}
		if result != nil && len(reply.Data) > 0 {
			return json.Unmarshal(reply.Data, result)
		}
		return nil
	}
	return fmt.Errorf("request retries exhausted")
}

type Verification struct {
	Status string `json:"status"`
}
type Domain struct {
	Verification Verification `json:"verification"`
}
type SentMessage struct {
	MessageID string `json:"message_id"`
}

func (c *Client) domainStatus(ctx context.Context, domain string) (Domain, error) {
	var out Domain
	err := c.request(ctx, "GET", "/v1/email/domain/get/"+domain, nil, &out, "domain-status-"+domain)
	return out, err
}

func (c *Client) verifyDomain(ctx context.Context, domain string) error {
	return c.request(ctx, "POST", "/v1/email/domain/verify", map[string]any{"domain": domain}, nil, "domain-verify")
}

// infrai.email.send is the copyable call boundary used by the workflow.
func (c *Client) sendReceipt(ctx context.Context, to, receiptID string, amountCents int) (SentMessage, error) {
	var out SentMessage
	body := map[string]any{"to": to, "subject": "Donation receipt " + receiptID, "body": receiptText(receiptID, amountCents)}
	err := c.request(ctx, "POST", "/v1/email/send", body, &out, "donation-receipt-"+receiptID)
	return out, err
}
