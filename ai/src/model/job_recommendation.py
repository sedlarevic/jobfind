from pydantic import BaseModel

class JobRecommendation(BaseModel):
    job_id: int
    title: str
    company: str
    fit_score: int
    reason: str
    main_strengths: list[str]
    main_gap: str | None

