# Audit Events Module

The Audit Events Module is designed to log every user action within the system. These logs are stored in a structured JSON format and can be retrieved for auditing purposes. The module supports both file-based and database-based logging.

## Features

- **Log User Actions**: Logs various user actions such as login, logout, index deletion, dashboard creation, etc.
- **Flexible Storage**: Supports logging to a file or a database (SQLite).
- **Retrieve Logs**: Allows retrieval of logs based on organization ID and time range.
- **Middleware Support**: Provides HTTP middleware for automatic logging of HTTP requests.

## Installation

To use the Audit Events Module, import it into your Go project.

## Usage

### Initialization

Initialize the audit logger with the desired storage type (file or database).

### Logging Events

Log user actions using the `CreateAuditEvent` function.

### Retrieving Logs

Retrieve logs for a specific organization and time range.

### HTTP Middleware

Use the provided middleware to automatically log HTTP requests.

## Example Log Format

```json
{
    "username": "JohnDoe",
    "actionString": "User logged in",
    "extraMsg": "Login from 192.168.1.1",
    "epochTimestampSec": 1745667898,
    "orgId": 123,
    "metadata": {
        "userAgent": "Mozilla/5.0",
        "ipAddress": "192.168.1.1"
    }
}
```
## Supported Actions

- **Auth**: User login, logout, password reset.
- **Index**: Create, delete, update indices.
- **Organization**: Update organization settings.
- **Dashboard**: Create, update, delete dashboards, toggle favorites.
- **Saved Queries**: Create, update, delete saved queries.
- **Folders**: Create, update, delete folders.
- **Alerts**: Create, update, delete alerts and contact points.
- **Lookup Files**: Create, delete lookup files.

## Configuration

### File Logger

To use the file logger, provide the file path in the configuration.

### Database Logger

To use the database logger, provide the database path.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.
