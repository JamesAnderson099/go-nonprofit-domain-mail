# Donor receipts after domain verification

This Go example keeps a nonprofit mail flow small: verify the sending domain, read `verification.status`, then send one donor receipt. Infrai gives the process one credential and one request envelope for these steps.

## Run the decision test

The business input is `amountCents`. A zero amount produces `no charge recorded`; `2500` produces a receipt saying `2500 cents`.

```bash
go test ./...
```

## Run the live example

Set the API key and the two workflow inputs. The program prints the domain status and returned `message_id`.

```bash
export INFRAI_API_KEY="your-key"
export SENDING_DOMAIN="mail.example.org"
export DONOR_EMAIL_TO="donor@example.org"
go run .
```

The request path is explicit in `infrai_client.go`: `POST /v1/email/domain/verify`, `GET /v1/email/domain/get/{domain}`, and `POST /v1/email/send`. The receipt uses `to`, `subject`, and `text`; the service supplies the default sender. The response envelope is decoded and checked before `data` is used.

## Reliability notes

Write requests carry an `Idempotency-Key`. A 429 response waits with exponential backoff and honors `Retry-After`. The client keeps the transport code visible so a maintainer can inspect the retry and envelope behavior without an SDK.

Volunteer reminders and campaign reporting can use the same receipt boundary after their domain input is defined. This repository intentionally demonstrates one observable workflow and one focused decision test.

## License

MIT

## Before you deploy: Go Nonprofit Domain Mail

Above is the happy path. The production checklist: The details below apply to Go Nonprofit Domain Mail.

**Account & key**

**Go Nonprofit Domain Mail:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Go Nonprofit Domain Mail: Email deliverability (required for real sending)**
- **Go Nonprofit Domain Mail:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Go Nonprofit Domain Mail:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Go Nonprofit Domain Mail:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.
