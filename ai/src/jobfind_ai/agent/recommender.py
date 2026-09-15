from agents import Agent

from jobfind_ai.config.prompts import RECOMMENDER_PROMPT
from jobfind_ai.config.settings import AGENT_MODEL
from jobfind_ai.model.recommendation_result import RecommendationResult
from jobfind_ai.tool.jobs import search_jobs

def new_recommender_agent() -> Agent:
    return Agent(
        name="Job Recommendation Agent",
        model=AGENT_MODEL,
        output_type=RecommendationResult,
        instructions = RECOMMENDER_PROMPT,
        tools=[search_jobs],
    )
