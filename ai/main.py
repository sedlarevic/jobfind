import os
import requests
from agents import Runner

from jobfind_ai.agent.cv_tailor import new_cv_tailor_agent
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

    output = result.final_output

    for i, recommendation in enumerate(output.recommendations, start=1):
        print(f"\n{i}. {recommendation.title} @ {recommendation.company}")
        print(f"Fit: {recommendation.fit_score}/100")
        print(f"Why: {recommendation.reason}")
        print(f"Strengths: {', '.join(recommendation.main_strengths)}")
        print(f"Main gap: {recommendation.main_gap or 'None'}")
        print(f"Job link: {recommendation.url}")

    while True:
        try:
            choice = int(input("\nSelect a job: "))

            if 1 <= choice <= len(output.recommendations):
                break

            print("Invalid selection.")
        except ValueError:
            print("Please enter a number.")

    selected_recommendation = output.recommendations[choice - 1]
    selected_job_id = selected_recommendation.job_id

    input_text = f"""
    Analyze this CV for job ID {selected_job_id}.

    Candidate CV:
    {extracted_cv}
    """

    tailor_agent = new_cv_tailor_agent()

    result = Runner.run_sync(
        tailor_agent,
        input_text,
    )
if __name__ == "__main__":
    main()
