import os
from openai import OpenAI
import openai
import psycopg
from psycopg.rows import dict_row
import requests

from jobfind_ai.model.job_posting import JobPosting

DATABASE_URL = os.environ["DATABASE_URL"]

def fetch_jobs_missing_embedding(connection: psycopg.Connection) -> list[JobPosting]:
    with connection.cursor(row_factory=dict_row) as cur:
        
        cur.execute("SELECT * FROM job_postings WHERE embedding IS NULL;")
        
        results: list[dict] = cur.fetchall()

    jobs: list[JobPosting] = []
    for result in results:
        job = JobPosting(
                 id=result["id"],
                 title=result["title"],
                 company=result["company_name"],
                 url=result["url"],
                 description=result["description"],
                 active=result["active"],
                 embedding=None
                )

        jobs.append(job)

    return jobs
    

def set_embedding(connection: psycopg.Connection, jobs: list[JobPosting]):
    client = OpenAI()
    with connection.cursor() as cur:
        for job in jobs:
            text = f"{job.title}\n{job.description}"
            response = openai.embeddings.create(
                input=text,
                model="text-embedding-3-small",
                encoding_format="float")

            embedding = response.data[0].embedding

            cur.execute(
                    """
                    UPDATE job_postings
                    SET embedding = %s
                    WHERE id = %s;
                    """, (embedding, job.id)
                    )

            print(f"Embedded job {job.id}: {job.title}")

        connection.commit()
        
def main():
    with psycopg.connect(DATABASE_URL) as connection:
        jobs = fetch_jobs_missing_embedding(connection)
        set_embedding(connection, jobs)

if __name__ == "__main__":
    main()
