import os
import requests
from agents import Agent, Runner, function_tool

from agent.recommender import new_recommender_agent 


def extract_cv(path: str | None):
    if path is None:
        return None

    with open(path, "rb") as f:
        response = requests.post(
                "http://localhost:8081/cv/extract",
                files={"cv": f},
                )
    response.raise_for_status()
        
    return response.json()["text"]

def main():
    agent = new_recommender_agent()
    cv_path = os.environ.get("CV_PATH")
    extracted_cv = extract_cv(cv_path)

if __name__ == "__main__":
    main()
