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
- Provides structured recommendations with fit scores, strengths, gaps, and job links
- Allows the user to select a recommended position
- Uses a second AI agent to analyze the CV against the selected job
- Provides job-specific CV strengths, gaps, and tailoring feedback

## Architecture

The project consists of two main parts: a Go backend and a Python AI layer.

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

The Go backend is responsible for collecting and storing job postings, extracting CV text, and exposing the HTTP API used by the AI layer.

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
        ├── cv.py
        └── jobs.py
```

The Python layer contains the AI agents, tools, and structured Pydantic output models.

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

The recommendation agent creates a concise semantic search query based on relevant evidence from the candidate CV and optional career preferences.
The query is converted into an embedding and compared with stored job embeddings using cosine distance in pgvector.
The most relevant postings are then evaluated in detail using the full candidate CV.
Semantic retrieval is therefore responsible for candidate retrieval, while the LLM agent performs the final reasoning and ranking.

## CV Tailoring Flow

After receiving recommendations, the user selects one position.

```text
Recommended jobs
  |
  v
User selects one job
  |
  v
Selected job ID
  |
  v
CV Tailor Agent
  |
  v
get_job_by_id
  |
  v
Selected job posting
  |
  v
CV-to-job analysis
  |
  v
Strengths + gaps + CV feedback
```

The CV Tailor Agent compares the candidate CV with the selected job posting.
The CV remains the source of truth for candidate skills, experience, education, projects, and achievements. The agent does not invent missing experience or technologies.
Its purpose is to identify relevant evidence already present in the CV and suggest how that evidence could be presented more effectively for the selected role.

## Candidate Note

The application supports an optional Candidate Note representing career preferences and desired direction.

Example:

```text
I only want C# or Python jobs, or jobs from company Nordeus.
```

The Candidate Note can influence semantic retrieval and ranking, but it is never treated as evidence of candidate capability.

## Example Output

A complete example of the recommendation and CV tailoring flow is available here:

[View example output](examples/example_output.md)

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

### Refresh Job Postings

```http
POST /jobs/refresh
```

Crawls registered company websites and synchronizes their job postings with the database.

### Get Active Job Postings

```http
GET /jobs/active
```

Returns all currently active job postings.

### Get Job by ID

```http
GET /jobs/{id}
```

Returns a single job posting by its database ID.

This endpoint is used by the CV Tailor Agent after the user selects a recommendation.

### Semantic Job Search

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

### Crawl a Single Company

```http
GET /crawler/{company}
```

Runs the crawler for one registered company without updating the database.

Example:

```http
GET /crawler/nordeus
```

### Extract CV Text

```http
POST /cv/extract
```

Accepts a PDF CV as multipart form data using the field:

```text
cv
```

The uploaded file is temporarily stored, parsed, normalized, and returned as text.

## Database

Job postings are stored in PostgreSQL.

The pgvector extension is used to store embedding vectors and perform semantic similarity search.

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

Embeddings are generated using `text-embedding-3-small`.

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

The project contains utility scripts used to prepare job data for development and demonstration.

### Lever Job Import

```text
ai/scripts/seed_lever.py
```

Imports public job postings from selected Lever career sites into PostgreSQL.

### Job Embeddings

```text
ai/scripts/embed_jobs.py
```

Generates embeddings for job postings that do not yet have an embedding stored in PostgreSQL.
These scripts are data-preparation utilities and are separate from the normal runtime recommendation flow.

## Current Status

Implemented:

- job crawling
- PostgreSQL persistence
- job refresh and deactivation
- PDF CV extraction
- semantic job retrieval using embeddings and pgvector
- LLM-based job recommendation
- Candidate Note preference handling
- structured recommendation output
- interactive job selection
- CV-to-job analysis
- CV tailoring feedback

Possible future improvements:

- additional company crawlers
- automatic embedding generation for newly added or updated jobs
- scheduled job refresh
- improved CLI or web interface
- additional job metadata such as location and work mode
- automated tests for semantic search and AI integration
