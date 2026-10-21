# Deployment, Packaging & Operational Guide

## 1. Building the Desktop Application

### 1.1 Standalone Windows Executable
To build a single portable Windows executable without an installer:

```powershell
# Step 1: Compile the Vue 3 frontend
npm --prefix frontend run build

# Step 2: Build the standalone Go executable
go build -ldflags "-s -w -H windowsgui" -o SchoolManagement.exe ./cmd/desktop
```
*(The `-H windowsgui` flag suppresses the background Windows console window).*

---

### 1.2 Windows Installer Package (NSIS)
Wails comes with built-in NSIS installer packaging support:

```powershell
wails build -nsis
```
- The resulting installer is created at:
  ```text
  build/bin/school-management-installer.exe
  ```
- **Per-User Installation:** Installs by default for the current user without requiring Administrator/UAC elevation.
- **Data Preservation:** Uninstalling the application binary does **not** delete student data or financial records.

---

## 2. File System & Application Data Directories

All application state is stored strictly within OS-approved user application directories:

```text
%LOCALAPPDATA%\SchoolManagement\
├── data/
│   ├── school.db            # Primary SQLite database
│   ├── school.db-wal        # Write-Ahead Log journal
│   └── school.db-shm        # Shared-memory index
├── documents/
│   └── ab/
│       └── cd/
│           └── <sha256>     # Content-addressed student documents
├── backups/                 # Local encrypted database snapshots
├── logs/
│   └── app.log              # Rotating slog operational log files
└── config.yaml              # Local configuration overrides
```

### Environment Overrides:
You can redirect storage paths for custom portable drive installations or automated testing:
- `SMS_DATA_DIR`: Base directory for all application data.
- `SMS_DB_PATH`: Specific absolute path to the SQLite database file.
- `SMS_LOG_LEVEL`: Logging verbosity (`debug`, `info`, `warn`, `error`).

---

## 3. Backup & Disaster Recovery

### 3.1 Consistent Online Database Backups
Because SQLite runs in WAL mode, copying `school.db` directly while the application is running could result in a corrupted snapshot.

Backups must use SQLite's native `VACUUM INTO` command:
```sql
VACUUM INTO 'C:\Users\<User>\AppData\Local\SchoolManagement\backups\backup_2026_10_20.db';
```
This safely creates a verified, defragmented, and consistent point-in-time snapshot of the database while transactions continue processing uninterrupted.

### 3.2 Complete System Backup Archive
A complete backup consists of:
1. The snapshot database file generated via `VACUUM INTO`.
2. The entire `documents/` directory containing content-addressed files.
3. A `manifest.json` file recording the backup timestamp, school identifier, and software version.

---

## 4. Path to LAN & Cloud Deployment

The School Management System is architected from day one to support seamless migration from a standalone desktop tool to a central LAN or cloud server.

### 4.1 Transition from Wails to Chi HTTP
- The application logic in `internal/application/` and the domain models in `internal/domain/` contain no Wails dependencies.
- To run as an HTTP web service, deploy `cmd/server/main.go`, which binds the Chi HTTP router to the exact same application command and query handlers.

### 4.2 Transition from SQLite to PostgreSQL
- The database schema is mirrored in `database/migrations/postgres/`.
- Repository interfaces in `internal/domain/` are implemented for both SQLite and PostgreSQL.
- Switching to PostgreSQL requires changing the database driver configuration in `config.yaml` without modifying any domain or frontend code.

