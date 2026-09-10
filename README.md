# JobFind

JobFind is an agentic job recommendation application that matches candidates with relevant job postings based on their CV, experience, and career preferences.

The application combines a Go backend for job collection, CV extraction, persistence, and semantic search with a Python AI layer built using the OpenAI Agents SDK.

Job postings are represented using embeddings and stored in PostgreSQL with pgvector. Semantic retrieval is used to identify relevant job candidates before an LLM-based recommendation agent evaluates them in more detail.

After receiving recommendations, the user can select a job and receive tailored CV feedback for that specific position.

## Features

- Scrapes job postings from company career pages
- Stores and updates job postings in PostgreSQL
- Tracks when job postings first appeared and when they were last seen
- Deactivates expired job postings
- Extracts and normalizes text from uploaded PDF CVs
- Generates vector embeddings for job postings
- Performs semantic job retrieval using PostgreSQL and pgvector
- Uses an AI agent to rank relevant jobs against a candidate CV
- Supports an optional Candidate Note for career preferences
- Provides structured job recommendations with:
  - fit score
  - explanation
  - candidate strengths
  - main gap
  - original job posting link
- Allows the user to select a recommended position
- Uses a second AI agent to analyze the CV against the selected job
- Provides job-specific CV strengths, gaps, and tailoring feedback

## Architecture

The project is divided into two main parts.

### Go Backend

```text
crawler/        Job posting collection
cv/             PDF CV extraction
preprocessing/  Text normalization
repository/     PostgreSQL persistence and semantic search
service/        Application logic
handler/        HTTP handlers and routes
model/          Backend data models
```

The Go backend is responsible for collecting and storing job postings, extracting CV text, and exposing the data required by the AI layer.

### Python AI Layer

```text
ai/
├── main.py
├── scripts/
│   ├── embed_jobs.py
│   └── seed_lever.py
└── src/jobfind_ai/
    ├── agent/
    │   ├── recommender.py
    │   └── cv_tailor.py
    ├── model/
    └── tool/
```

The Python layer contains the AI agents, their tools, and structured Pydantic output models.

## Recommendation Flow

```text
PDF CV
  |
  v
CV text extraction
  |
  v
Candidate CV + optional Candidate Note
  |
  v
Job Recommendation Agent
  |
  v
Semantic search query
  |
  v
OpenAI embedding
  |
  v
PostgreSQL + pgvector
  |
  v
Top semantically relevant jobs
  |
  v
LLM evaluation and ranking
  |
  v
Up to five recommendations
```

Semantic retrieval is used as a first-stage retrieval mechanism.

The recommendation agent creates a concise search query based on the candidate's CV and preferences. The query is converted into an embedding and compared with stored job embeddings using cosine distance in pgvector.

The most relevant postings are then evaluated by the recommendation agent using the full candidate CV.

## CV Tailoring Flow

After receiving recommendations, the user selects one position.

```text
Selected recommendation
  |
  v
Job ID
  |
  v
CV Tailor Agent
  |
  v
Get selected job posting
  |
  v
Compare job posting with CV
  |
  v
Strengths + gaps + CV feedback
```

The CV Tailor Agent does not invent candidate experience or skills. The existing CV remains the source of truth.

Its purpose is to identify relevant evidence already present in the CV and suggest how that evidence could be presented more effectively for the selected position.

## Candidate Note

The application supports an optional Candidate Note.

Example:

```text
I only want backend Go roles.
```

The Candidate Note represents career preferences and desired direction.

It may influence retrieval and ranking, but it is never treated as evidence that the candidate possesses a particular skill or experience.

## Technology Stack

### Backend

- Go 1.25
- PostgreSQL
- pgvector
- pgx / pgxpool
- Colly
- ledongthuc/pdf
- Go standard `net/http`

### AI

- Python 3.13
- OpenAI Agents SDK
- OpenAI embeddings
- Pydantic
- requests
- psycopg
- uv

## HTTP API

### Refresh job postings

```http
POST /jobs/refresh
```

Crawls registered company websites and synchronizes their job postings with the database.

### Get active job postings

```http
GET /jobs/active
```

Returns all currently active job postings.

### Get job by ID

```http
GET /jobs/{id}
```

Returns one job posting by its database ID.

This endpoint is used by the CV Tailor Agent after the user selects a recommendation.

### Semantic job search

```http
POST /jobs/search
```

Accepts an embedding vector and a result limit.

Example request:

```json
{
  "embedding": [0.01, -0.02, 0.03],
  "limit": 50
}
```

The backend performs cosine-distance search against stored job embeddings using pgvector.

### Crawl a single company

```http
GET /crawler/{company}
```

Runs a registered company crawler without updating the database.

Example:

```http
GET /crawler/nordeus
```

### Extract CV text

```http
POST /cv/extract
```

Accepts a PDF CV as multipart form data using the field:

```text
cv
```

The uploaded PDF is temporarily stored, parsed, normalized, and returned as text.

## Database

Job postings are stored in PostgreSQL.

The pgvector extension is used to store OpenAI embedding vectors and perform semantic similarity search.

```sql
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE job_postings (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    company_name TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    first_seen TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    embedding VECTOR(1536)
);
```

Embeddings are currently generated using `text-embedding-3-small`.

## Running the Backend

Set the database connection:

```bash
export DATABASE_URL="postgres://user@localhost:5432/jobfind"
```

Run the Go backend:

```bash
go run .
```

The HTTP server starts on:

```text
localhost:8081
```

## Running the AI Application

The AI application requires an OpenAI API key and the path to a PDF CV.

```bash
export OPENAI_API_KEY="..."
export CV_PATH="/path/to/cv.pdf"
```

An optional career preference can also be provided:

```bash
export CANDIDATE_NOTE="I only want backend Go roles."
```

From the `ai` directory:

```bash
PYTHONPATH=src uv run python main.py
```

The application recommends suitable jobs, allows the user to select one, and then generates CV feedback for the selected position.

## Demo Data

The project contains two utility scripts used to prepare job data for development and demonstration.

```text
ai/scripts/seed_lever.py
```

Imports public job postings from selected Lever career sites.

```text
ai/scripts/embed_jobs.py
```

Generates embeddings for job postings that do not yet have an embedding stored in PostgreSQL.

These scripts are data preparation utilities and are separate from the normal runtime recommendation flow.

## Current Status

Implemented:

- job crawling
- PostgreSQL persistence
- job refresh and deactivation
- PDF CV extraction
- semantic job retrieval with embeddings and pgvector
- LLM-based job recommendation
- Candidate Note preference handling
- structured recommendation output
- interactive job selection
- CV-to-job analysis
- CV tailoring feedback

Possible future improvements:

- additional company crawlers
- automatic embedding generation when new jobs are inserted
- scheduled job refresh
- improved CLI or web interface
- additional job metadata such as location and work mode
