import os
import re
import psycopg
import requests


DATABASE_URL = os.environ["DATABASE_URL"]

LEVER_SITES = {
    "Xsolla": "xsolla",
    "Collectly": "CollectlyInc",
    "Filevine": "filevine",
    "Level AI": "levelai",
    "RIVR": "rivr",
    "Aleph": "aleph",
    "Unlimit": "unlimit",
    "Legend": "Legend",
    "SFEIR": "sfeir",
    "Supermove": "supermove",
}
def strip_html(value: str) -> str:
    return " ".join(
        re.sub(r"<[^>]+>", " ", value).split()
    )


def build_description(job: dict) -> str:
    parts = []

    description = job.get("descriptionPlain", "")
    if description:
        parts.append(description)

    for section in job.get("lists", []):
        title = section.get("text", "")
        content = strip_html(section.get("content", ""))

        if title:
            parts.append(title)

        if content:
            parts.append(content)

    additional = job.get("additionalPlain", "")
    if additional:
        parts.append(additional)

    return "\n".join(parts)

def fetch_jobs(company_name: str, site: str) -> list[dict]:
    url = f"https://api.lever.co/v0/postings/{site}"

    response = requests.get(
        url,
        params={
            "mode": "json",
            "limit": 100,
        },
        timeout=20,
    )
    response.raise_for_status()

    jobs = response.json()

    print(f"{company_name}: received {len(jobs)} jobs")

    return jobs


def insert_jobs(
    connection: psycopg.Connection,
    company_name: str,
    jobs: list[dict],
):
    for job in jobs:
        title = job["text"]
        url = job["hostedUrl"]
        description = build_description(job)
        connection.execute(
            """
            INSERT INTO job_postings (
                company_name,
                title,
                url,
                description,
                first_seen,
                last_seen,
                active
            )
            VALUES (
                %s,
                %s,
                %s,
                %s,
                NOW(),
                NOW(),
                TRUE
            )
            ON CONFLICT (url)
            DO UPDATE SET
                title = EXCLUDED.title,
                description = EXCLUDED.description,
                last_seen = NOW(),
                active = TRUE
            """,
            (
                company_name,
                title,
                url,
                description,
            ),
        )


def main():
    with psycopg.connect(DATABASE_URL) as connection:
        total = 0

        for company_name, site in LEVER_SITES.items():
            try:
                jobs = fetch_jobs(company_name, site)

                insert_jobs(
                    connection,
                    company_name,
                    jobs,
                )

                total += len(jobs)

            except Exception as e:
                print(f"{company_name}: failed: {e}")

        connection.commit()

    print(f"Done. Imported {total} jobs.")


if __name__ == "__main__":
    main()
