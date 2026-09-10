from agents import function_tool
import requests
from openai import OpenAI

from jobfind_ai.model.job_posting import JobPosting

client = OpenAI()

@function_tool
def search_jobs(query: str, limit: int = 50) -> list[dict]:
    """
    Search active job postings by semantic similarity.

    Args:
        query: Description of the type of job that should be retrieved.
        limit: Maximum number of job postings to return.
    """

    print(f"TOOL: search_jobs called with query: {query}")

    response = client.embeddings.create(
        model="text-embedding-3-small",
        input=query,
        encoding_format="float",
    )

    embedding = response.data[0].embedding

    search_response = requests.post(
        "http://localhost:8081/jobs/search",
        json={
            "embedding": embedding,
            "limit": limit,
        },
        timeout=20,
    )

    search_response.raise_for_status()

    jobs = search_response.json()

    print(f"TOOL: semantic search returned {len(jobs)} jobs")

    return [
        {
            "id": job["id"],
            "company": job["companyName"],
            "title": job["title"],
            "description": job["description"],
            "url": job["url"],
        }
        for job in jobs
    ]

@function_tool
def get_job_by_id(id: int) -> JobPosting:
    """
    Get job posting based on id.

    Args:
        id: job posting id
    """

    print(f"TOOL: get_job_by_id called with id: {id}")


    search_response = requests.post(
        f"http://localhost:8081/jobs/{id}",
        timeout=20,
    )

    search_response.raise_for_status()

    job = search_response.json()

    print(f"TOOL: get_job_by_id returned a job: {job}")

    return JobPosting(
            id=job["id"],
            company=job["companyName"],
            title=job["title"],
            description=job["description"],
            url=job["url"],
            active=job["active"]
            )
