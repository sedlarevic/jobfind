from agents import function_tool
import requests

@function_tool
def get_active_jobs() -> list[dict]:
    response = requests.get("http://localhost:8081/jobs/active", timeout = 10,)
    response.raise_for_status()

    return response.json()

