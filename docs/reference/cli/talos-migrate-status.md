---
id: talos-migrate-status
title: talos migrate status
description: talos migrate status
---

<!--
This file is auto-generated.

To improve this file please make your change against the appropriate "./cmd/*.go" file.
-->

## talos migrate status

Show migration status

### Synopsis

Display the current database migration status.

Shows:

- Current migration version
- Latest migration version available in this binary
- Whether the database is in a dirty state

With --block, the command polls the database every second and only returns once all bundled
migrations have been applied. Use this to gate a rollout on another process (such as the primary
cluster) finishing 'migrate up'.

```
talos migrate status [flags]
```

### Examples

```
  # Show the migration status
  talos migrate status --database "sqlite3://./data/talos.db"

  # Block until all migrations have been applied
  talos migrate status --block --database "sqlite3://./data/talos.db"
```

### Options

```
      --block             Block until all migrations have been applied
      --database string   database DSN (overrides DB_DSN)
  -h, --help              help for status
```

### Options inherited from parent commands

```
      --config string     config file (default is $HOME/.talos.yaml or ./config.yaml)
  -e, --endpoint string   HTTP server base URL including scheme, e.g. http://host:port (for client commands) (default "http://localhost:4420")
```

### See also

- [talos migrate](talos-migrate) Database migration tools
