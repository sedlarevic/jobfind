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
    cv_path = os.environ.get("CV_PATH")

    extracted_cv = extract_cv(cv_path)
    
    if extracted_cv is None:
        raise ValueError("CV_PATH env variable is not set")

    agent = new_recommender_agent()

    input_text = f"""
    Recommend the five best available jobs for this candidate.

    Candidate CV:
    {extracted_cv}
    """

    candidate_note = os.environ.get("CANDIDATE_NOTE")
    if candidate_note:
        input_text += f"""

        Candidate Note:
        {candidate_note}
        """

    result = Runner.run_sync(
            agent,
            input_text
            )

    print(result.final_output)

if __name__ == "__main__":
    main()
