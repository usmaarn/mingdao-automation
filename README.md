# mingdao-automation

A local HTTP service that drives a Mingdao web form through Playwright. Send a
JSON payload describing one or more records and the service fills and saves each
record in the browser automatically.

## Requirements

- [Bun](https://bun.sh/)
- A Mingdao account with access to the target app

## Setup

```bash
bun install
bunx playwright install chromium
```

The Mingdao page URL is configured in `src/index.ts`:

```ts
const url = "https://mingdaoyun.palmpay-inc.com/app/...";
```

Change it there if you need to target a different view.

## Running

```bash
bun start
```

What happens on startup:

1. A Chromium window opens using the persistent profile in `./browser-data`.
2. The configured Mingdao URL is loaded.
3. The add-record modal is opened and reset.
4. The HTTP server starts on a random free port and prints the URL, e.g.
   `Server running at http://localhost:57123/`.

On the first run you must log in manually in the opened browser. The session is
stored in `./browser-data`, so later runs reuse it. Keep the browser window open
while the service is running; do not close it or interact with the form.

To use a fixed port:

```bash
PORT=8080 bun start
```

## Sending records

Send a `POST` request with a JSON body. The service responds immediately with
`{"message":"queued"}` and processes records one at a time in the background.
Watch the terminal for `Record N completed.` log lines.

```bash
curl -X POST http://localhost:8080/ \
  -H "Content-Type: application/json" \
  -d '{
    "merchantName": "Acme Store",
    "data": [
      {
        "summary": "POS terminal not printing receipts",
        "solution": "Replaced paper roll and restarted the terminal.",
        "startDate": "2026-10-29 09:00",
        "endDate": "2026-10-29 10:30",
        "type": "issue locating"
      }
    ]
  }'
```

### Payload fields

| Field          | Required | Description                                                        |
| -------------- | -------- | ------------------------------------------------------------------ |
| `merchantName` | no       | Merchant to link to every record. Leave empty to skip linking.     |
| `data`         | yes      | Non-empty array of records (see below).                            |

Each record in `data`:

| Field       | Required | Description                                            |
| ----------- | -------- | ------------------------------------------------------ |
| `summary`   | yes      | Non-empty summary text.                                |
| `solution`  | yes      | Non-empty solution text.                               |
| `startDate` | yes      | Format `YYYY-MM-DD HH:mm`.                             |
| `endDate`   | yes      | Format `YYYY-MM-DD HH:mm`.                             |
| `type`      | yes      | One of the supported types listed below.               |

Supported `type` values:

- `issue locating`
- `documentation` (or the legacy spelling `documentantion`)
- `phone call`
- `meeting`

The value is matched case-insensitively.

### Responses

| Status | Meaning                                                    |
| ------ | ---------------------------------------------------------- |
| `200`  | `{"message":"queued"}` — records accepted for processing. |
| `400`  | Invalid JSON or failed validation (see `errors`).         |
| `405`  | Method other than `POST`.                                  |

Validation errors use Zod's flattened format:

```json
{
  "errors": {
    "formErrors": [],
    "fieldErrors": {
      "data": ["At least one record is required"]
    }
  }
}
```

## Notes and limitations

- Records are processed sequentially under a mutex; a long batch delays
  subsequent requests.
- Processing is asynchronous. A `200` means the payload passed validation, not
  that the browser work succeeded. Check the terminal for errors.
- If the merchant name is not found, the record is still saved without a linked
  merchant and a warning is logged.
- The current implementation does not upload screenshots; `screenshotPath` is
  defined but unused by the HTTP API.
- `browser-data` is git-ignored and contains your login session. Delete it to
  force a fresh login.

## Development

```bash
bun run typecheck
```
