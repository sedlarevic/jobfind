from pydantic import BaseModel

from jobfind_ai.model.job_posting import JobPosting


class CVTailorResult(BaseModel):
    job_id: int
    strengths: list[str]
    gaps: list[str]
    cv_feedback: str
