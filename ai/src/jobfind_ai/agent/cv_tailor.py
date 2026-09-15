
from agents import Agent

from jobfind_ai.config.prompts import CV_TAILOR_PROMPT
from jobfind_ai.config.settings import AGENT_MODEL
from jobfind_ai.model.cv_tailor_result import CVTailorResult
from jobfind_ai.tool.jobs import get_job_by_id


def new_cv_tailor_agent() -> Agent:
    return Agent(
            name="CV Tailor Agent",
            model=AGENT_MODEL,
            output_type=CVTailorResult,
            instructions=CV_TAILOR_PROMPT,
            tools=[get_job_by_id]
            )
