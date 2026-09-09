import random
from agents import function_tool
import requests

@function_tool
def get_active_jobs() -> list[dict]:
    """Retrieve all currently active job postings from the JobFind backend."""
    print("TOOL: get_active_jobs tool called")
    response = requests.get("http://localhost:8081/jobs/active", timeout = 10,)
    response.raise_for_status()

    jobs = response.json()

    print("TOOL: received jobs:", len(jobs))
    
    random.shuffle(jobs)

    jobs = jobs[:100]

    print("TOOL: sending jobs to agent:", len(jobs))

    return [
        {
            "id": job["id"],
            "company": job["companyName"],
            "title": job["title"],
            "description": job["description"][:3000],
        }
        for job in jobs
    ]

