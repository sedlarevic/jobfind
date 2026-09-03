# JobFind

JobFind is a Go application for collecting job postings and recommending relevant positions based on the content of a user's CV.

The project combines web scraping, text preprocessing, PostgreSQL persistence, and content-based recommendation using TF-IDF and cosine similarity.

## Features

- Scrapes job postings from company career pages
- Stores and updates job postings in PostgreSQL
- Tracks when job postings first appeared and when they were last seen
- Deactivates expired job postings
- Extracts text from uploaded PDF CVs
- Content-based job recommendation using TF-IDF and cosine similarity
- Provides HTTP endpoints for crawling, refreshing, retrieving jobs, and CV extraction

## Architecture

The application is organized into several layers:

```text
crawler/        Web scraping logic
cv/             PDF text extraction and CV-specific cleanup
preprocessing/  Shared text preprocessing
repository/     PostgreSQL persistence
service/        Application and business logic
handler/        HTTP handlers and routes
model/          Application data models
```

Main application flow:

```text
Company websites
      |
      v
   Crawler
      |
      v
   Service
      |
      v
 Repository
      |
      v
 PostgreSQL
```

CV processing flow:

```text
PDF CV
   |
   v
Text extraction
   |
   v
CV cleanup
   |
   v
Text preprocessing
   |
   v
Recommendation algorithm
```

## Technology Stack

- Go 1.25
- PostgreSQL
- pgxPool
- Colly
- ledongthuc/pdf
- Go standard `net/http` package

## HTTP API

### Refresh job postings

```http
POST /jobs/refresh
```

Crawls all registered company websites and synchronizes the results with the database.

### Get active job postings

```http
GET /jobs/active
```

Returns all currently active job postings.

### Crawl a single company

```http
GET /crawler/{company}
```

Runs the crawler for a specific company without updating the database.

Example:

```http
GET /crawler/nordeus
```

### Extract CV text

```http
POST /cv/extract
```

Accepts a PDF CV as multipart form data using the field name:

```text
cv
```

The uploaded file is temporarily stored, parsed, normalized, and returned as text.

## Database

The application uses the following PostgreSQL table:

```sql
CREATE TABLE job_postings (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    company_name TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    active BOOLEAN NOT NULL DEFAULT TRUE
);
```

The database connection is configured through the `DATABASE_URL` environment variable.

Example:

```text
postgres://postgres:password@localhost:5432/jobfind
```

## Running the Application

Set the database connection:

```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/jobfind"
```

Run the application:

```bash
go run .
```

The HTTP server starts on:

```text
localhost:8081
```

## Current Status

Implemented:

- Crawler
- PostgreSQL persistence
- Job refresh and deactivation logic
- Active job retrieval
- PDF CV text extraction
- Text normalization
- HTTP API

In progress:

- TF-IDF vectorization
- Cosine similarity
- Job ranking

Planned:

- Additional company crawlers
- Scheduled daily job refresh
- Command-line interface
