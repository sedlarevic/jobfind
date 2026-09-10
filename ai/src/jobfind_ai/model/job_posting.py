from pydantic import BaseModel


class JobPosting(BaseModel):
    id: int
    title: str
    company: str
    url: str
    description: str
    active: bool
    embedding: list[float] | None = None
