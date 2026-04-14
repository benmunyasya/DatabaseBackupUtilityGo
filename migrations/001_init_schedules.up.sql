CREATE TABLE schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    db_type TEXT NOT NULL,
    db_name TEXT NOT NULL,
    backup_type TEXT NOT NULL,
    cron_expr TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
