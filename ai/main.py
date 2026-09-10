import os
import requests
from agents import Runner

from jobfind_ai.agent.recommender import new_recommender_agent
from jobfind_ai.model.job_posting import JobPosting 


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
    print("1. main started")
    cv_path = os.environ.get("CV_PATH")
    print("1. cv path:", cv_path)

    extracted_cv = extract_cv(cv_path)
    print("3. CV extracted:", len(extracted_cv) if extracted_cv else None)

    if extracted_cv is None:
        raise ValueError("CV_PATH env variable is not set")

    agent = new_recommender_agent()
    print("4. agent created")

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

    print("5. starting agent")

    result = Runner.run_sync(
            agent,
            input_text
            )

    print("6. agent finished")

    output = result.final_output

    print(output)
    print(output.model_dump_json(indent=2))

    for recommendation in output.recommendations:
        print(f"\n{recommendation.title} @ {recommendation.company}")
        print(f"Fit: {recommendation.fit_score}")
        print(f"Reason: {recommendation.reason}")
        print(f"Strengths: {recommendation.main_strengths}")
        print(f"Gap: {recommendation.main_gap}")

if __name__ == "__main__":
    main()
