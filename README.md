# Database Backup Utility (backup-cli)

A developer‑friendly CLI tool to backup and restore multiple databases with compression, scheduling, and cloud storage support.  
Built in Go, it uses [Cobra](https://github.com/spf13/cobra) for CLI commands, [golang-migrate](https://github.com/golang-migrate/migrate) for schema migrations, and [robfig/cron](https://github.com/robfig/cron) for scheduling.
https://roadmap.sh/projects/database-backup-utility
---

## ✨ Features
- **Manual backups**: Run full, incremental, or differential backups for Postgres, MySQL, MongoDB, or SQLite.
- **Restore support**: Restore databases from compressed backup files.
- **Scheduling**: Persist backup jobs in SQLite and run them automatically using cron expressions.
- **Compression**: Backups are saved as `.sql.gz` files to save space.
- **Persistence**: All schedules are stored in `schedules.db` and reloaded at startup.
- **Logging**: Clear `[INFO]`, `[SUCCESS]`, `[ERROR]` messages for transparency.

---

## 📂 Project Structure
cmd/backup-cli/       # CLI entrypoint and commands (backup, restore, schedule, config)
internal/db/          # Database adapters (Postgres, MySQL, Mongo, SQLite)
internal/backup/      # Backup runner logic
internal/scheduler/   # ScheduleManager with cron integration
internal/log/         # Logging utility
migrations/           # SQL migrations for schedules table

---

## 🚀 Getting Started

### Prerequisites
- Go 1.20+
- Postgres/MySQL/MongoDB/SQLite installed locally
- Environment variables set for DB connection:
POSTGRES_DB_HOST=localhost
POSTGRES_DB_PORT=5432
POSTGRES_DB_USER=youruser
POSTGRES_DB_PASSWORD=yourpass
### Build
```bash
go build -o backup-cli ./cmd/backup-cli
Run
Manual backup:

bash
./backup-cli backup --db postgres --type full --db-name taskmanager
Schedule a backup:

bash
./backup-cli schedule --db postgres --db-name taskmanager --backup full --cron "0 2 * * *"
Restore:

bash
./backup-cli restore --db postgres --db-name taskmanager --file backups/taskmanager_full_20260414_122300_backup.sql.gz
Migrations
The schedules table is managed by golang-migrate. Migration files live in migrations/ and follow the convention:

Code
001_init_schedules.up.sql
001_init_schedules.down.sql
🕒 Scheduler
Jobs are stored in schedules.db.

On startup, ScheduleManager loads all jobs and registers them with cron.

New jobs created via schedule command are persisted and registered immediately.

The process must remain running to execute scheduled jobs.