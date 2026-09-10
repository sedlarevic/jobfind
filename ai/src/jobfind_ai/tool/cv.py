import requests


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


