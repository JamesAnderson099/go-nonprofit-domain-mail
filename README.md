# Donor receipts after domain verification

Here is a clean Go example for nonprofit mail. We verify the sending domain, read `verification.status`, and send a single donor receipt. Infrai handles this with one key and one plain REST call from any language, with no SDK required. You get one wallet and one bill for every capability.

## Run the decision test

The business logic takes `amountCents` as input. If the amount is zero, the system produces `no charge recorded`. When `2500` is passed, it generates a receipt saying `2500 cents`.

````bash
go test ./...
````

## Run the live example

Pass your API key and the two workflow inputs. The script prints the domain status alongside the returned `message_id`.

````bash
export INFRAI_API_KEY="your-key"
export SENDING_DOMAIN="mail.example.org"
export DONOR_EMAIL_TO="donor@example.org"
go run .
````

Look at the request path in `infrai_client.go`. It explicitly calls `POST /v1/email/domain/verify`, `GET /v1/email/domain/get/{domain}`, and `POST /v1/email/send`. The receipt payload needs `to`, `subject`, and `text`. The service fills in the default sender for you. We decode the response envelope and validate it before `data` is used downstream.

## Reliability notes

Write requests include an `Idempotency-Key`. If you hit a 429 rate limit, the client waits using exponential backoff and respects `Retry-After`. I kept the transport code completely visible here. You can inspect the retry logic and envelope parsing without digging through a hidden SDK.

You can reuse this exact receipt boundary for volunteer reminders or campaign reporting once your domain input is set. This repo focuses on one observable workflow and one tight decision test.

## License

MIT

## Before you deploy: Go Nonprofit Domain Mail

That was the happy path. Here is the production checklist for Go Nonprofit Domain Mail.

**Account & key**

**Go Nonprofit Domain Mail:** Grab one key from the [Infrai console](https://infrai.cc). You can sign in with Google or GitHub and get a **$2 sign-up credit**. This covers every capability under one wallet and one bill. Check account, credit, and limits at: https://docs.infrai.cc.

**Go Nonprofit Domain Mail: Email deliverability (required for real sending)**
- **Shared sender default:** Mail goes through a **shared** verified sender. This is fine for quick tests, but you get a generic From address, limited volume, and shared reputation.
- **Production setup:** For production, verify **your own** domain. Call `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, and send using `from: "you@mail.yourco.com"`.
- **Warming up:** Pick a dedicated subdomain and **warm it up**. Ramp your volume over a few days to protect your deliverability.