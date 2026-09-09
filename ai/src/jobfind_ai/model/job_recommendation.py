from pydantic import BaseModel, Field

class JobRecommendation(BaseModel):
    job_id: int
    title: str
    company: str
    fit_score: int = Field(ge=0, le=100)
    reason: str
    main_strengths: list[str]
    main_gap: str | None = None

