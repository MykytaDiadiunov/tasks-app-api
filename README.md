# Go Boilerplate

A starter template for building RESTful APIs in Go, designed with a clean architecture and best practices. This boilerplate includes essential tools and configurations to streamline development.

## Features

- **Structured Architecture**:
  - Clear separation of layers: domain, infrastructure, repositories, and services.
- **Logging**:
  - Integrated with `zap` for efficient and structured logging.
- **Database Support**:
  - Pre-configured PostgreSQL driver and database connection.
  - Built-in support for migrations to manage schema changes.
- **Docker Support**:
  - Containerized setup for easier deployment using Docker and Docker Compose.
- **Configuration Management**:
  - Manage settings via a dedicated `config` package with environment variable support.
- **HTTP Server**:
  - Ready-to-use structure for implementing RESTful APIs.

## Installation

```
git clone https://github.com/MykytaDiadiunov/go-boilerplate.git
cd go-boilerplate
```
```
docker-compose up --build
```

## Configure environment variables

```
DB_NAME= {your db name}
DB_HOST= {your db host}
DB_PORT= {your db port}
DB_PORT_EXTERNAL= {your db external port}
DB_USER= {your db user}
DB_PASSWORD= {your db password}
JWT_SECRET= {your jwt secret}
MIGRATE= {your migration version}
MIGRATION_LOCATION=migrations // Path to migrations folder
SMTP_HOST= smtp.gmail.com
SMTP_PORT= 587
WORK_GMAIL= {your sender gmail}
WORK_GMAIL_PASSWORD= {your sender gmail password} 
LOGGER_LEVEL= {your logger level}
CLOUDINARY_NAME_KEY= {your cloudinary name key}
CLOUDINARY_API_KEY= {your cloudinary api key}
CLOUDINARY_SECRET_KEY= {your cloudinary secret key}
```

