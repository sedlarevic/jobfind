from pydantic import BaseModel
from model.job_recommendation import JobRecommendation

class RecommendationResult(BaseModel):
    recommendations: list[JobRecommendation]
