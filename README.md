# userAnalytics

A lightweight, self-hosted analytics service built with Go and Turso (libSQL).

Live service here: https://analyticz.vercel.app/

## Quick Start (Client Side)
Add this script to your HTML to start tracking visitors:

```html
<script>
  var dID = "YOUR_DOMAIN_ID" 
</script>
<script async src="https://cdn.jsdelivr.net/gh/NikoJunttila/userAnalytics@main/javascript/tracker.js" type="text/javascript"></script>
```

---

## Environment Variables
The application requires the following environment variables (stored in `.env` or passed to Docker):

| Variable | Description | Example |
| :--- | :--- | :--- |
| `PORT` | The port the HTTP server will listen on. | `8000` |
| `DB_URL` | Turso (libSQL) or local SQLite connection string. | `libsql://your-db-org.turso.io?authToken=...` |
| `emailCode` | Secret / App Password for SMTP email notifications. | `your-secret-code` |

---

## Docker
The project is containerized for easy deployment.

### Build the Image
```bash
docker build -t analytics-app .
```

### Run the Container
```bash
docker run -p 8000:8000 \
  -e PORT=8000 \
  -e DB_URL="libsql://your-db.turso.io?authToken=your-token" \
  -e emailCode="your-email-secret" \
  analytics-app
```

### Automatic Migrations
The application uses **Goose** for schema management. On every startup, it will automatically check for and apply any pending migrations located in `sql/schema/` to the database specified in `DB_URL`.

---

## Development
- **Database**: Ported to Turso/SQLite.
- **Migrations**: Auto-applied via `goose` on startup.
- **Build**: `go build -o analytics .`
