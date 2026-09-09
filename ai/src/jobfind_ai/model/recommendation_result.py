from pydantic import BaseModel, Field
from jobfind_ai.model.job_recommendation import JobRecommendation

class RecommendationResult(BaseModel):
    recommendations: list[JobRecommendation] = Field(max_length=5)
