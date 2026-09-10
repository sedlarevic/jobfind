import random
from agents import function_tool
import requests
from openai import OpenAI

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
        }
        for job in jobs
    ]
