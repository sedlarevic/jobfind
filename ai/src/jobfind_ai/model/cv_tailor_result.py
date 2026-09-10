from pydantic import BaseModel

class CVTailorResult(BaseModel):
    job_id: int
    strengths: list[str]
    gaps: list[str]
    cv_feedback: str
