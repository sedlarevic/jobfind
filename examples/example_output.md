# Example Output

This example demonstrates the complete JobFind flow:

```text
CV
- Candidate Note
- semantic job retrieval
- AI recommendation
- user selection
- job-specific CV analysis
```

The candidate provides a PDF CV and the following optional preference:

```bash
export CANDIDATE_NOTE="I only want C# or Python jobs, or jobs from company Nordeus."
```

The application is then started from the `ai` directory:

```bash
PYTHONPATH=src uv run python main.py
```

## Job Recommendation

The Job Recommendation Agent first creates a semantic search query from relevant CV evidence and the Candidate Note.

```text
TOOL: search_jobs called with query: Junior software backend engineer Python C# roles, candidate with Python modernization, testing, Docker Kubernetes CI, C# .NET performance profiling, SQL, Go and C systems; prefer Nordeus

TOOL: semantic search returned 50 jobs
```

The retrieved postings are then evaluated against the complete CV.

```text
1. QA Software Engineer Intern @ NORDEUS
Fit: 82/100
Why: Strong entry-level fit combining software engineering, testing, game development, and Belgrade location. The CV demonstrates programming in Python, C#, C, and Go, automated testing with pytest, and game-engine development.
Strengths: Python and C# programming experience, Automated testing with pytest, monkeypatching, and testcontainers, Game-engine project in C and gaming interest, Current master's studies in software engineering and AI
Main gap: The CV does not provide direct QA experience or evidence of functional, regression, or performance testing.
Job link: https://nordeus.com/open-positions/engineering_&_science/8090556/

2. Junior Data Scientist @ NORDEUS
Fit: 74/100
Why: Relevant Nordeus opportunity supported by Python, SQL, Power BI, data warehousing, ETL, and an information-systems degree. The role also aligns with the candidate's ongoing software engineering and AI master's studies.
Strengths: Python and SQL skills, Designed a SQL Server data warehouse and SSIS ETL pipelines, Built Power BI dashboards, Master's studies in software engineering and AI
Main gap: The CV does not provide evidence of statistics, machine learning, R, or A/B testing experience.
Job link: https://nordeus.com/open-positions/engineering_&_science/8104961/

3. Full Stack AI Engineer / Applied AI Engineer @ Xsolla
Fit: 64/100
Why: The role matches the candidate's Python background, AI-focused master's studies, and Docker/Kubernetes exposure. It is an attractive direction, but the CV does not show full-stack or GenAI production experience.
Strengths: Python programming and FastAPI/Flask listed, Master's program in software engineering and AI, Docker and Kubernetes exposure through internship, Backend and CLI development experience
Main gap: No CV evidence of TypeScript, Next.js, frontend development, LLM applications, or vector databases.
Job link: https://jobs.lever.co/xsolla/ebb74747-c739-4247-bafc-a088c10fa643

4. Software Engineer, Backend & Integrations (Remote) @ Supermove
Fit: 58/100
Why: Python/Flask backend, SQL, testing, and integrations align with the candidate's Python internship and data-engineering background. The role is less suitable than the Nordeus options because it expects 2–5 years of relevant experience and AWS/GraphQL exposure.
Strengths: Python experience and Flask listed, Backend CLI and internal-tool development, SQL and data-engineering experience, Automated testing experience
Main gap: The CV does not provide evidence of GraphQL, AWS, or 2–5 years of relevant professional experience.
Job link: https://jobs.lever.co/supermove/a0d545ec-c114-45e5-8707-16e5c1997182

5. Développeur C# / .Net (F/H/NB) @ SFEIR
Fit: 43/100
Why: C# and .NET performance research, profiling, benchmarking, and a relevant master's thesis provide a foundation for this role. However, the position is explicitly senior and requires at least five years of .NET development.
Strengths: C# experience through master's thesis, .NET profiling and BenchmarkDotNet experience, Software engineering master's studies, Docker and Kubernetes listed
Main gap: The CV does not show the required five years of .NET experience or production experience with .NET 8, microservices, Kafka, and cloud platforms.
Job link: https://jobs.lever.co/sfeir/12d3d218-abbc-47f5-814b-0098b3399a59
```

## User Selection

The user selects one of the recommended positions:

```text
Select a job: 4
```

The application takes the corresponding database job ID and passes it to the CV Tailor Agent.

```text
TOOL: get_job_by_id called with id: 527
TOOL: get_job_by_id returned a job.
```

## CV Tailoring Result

The CV Tailor Agent retrieves the selected posting and compares it against the original candidate CV.

```text
CV TAILORING RESULT

Strengths:
- Python backend internship experience involving modernization of 20+ internal libraries, CLI development, refactoring, pytest-based automated testing, Docker, Kubernetes, CI workflows, and code reviews.
- Existing Flask and FastAPI skills align with the role’s Python backend focus.
- Testing experience includes pytest, monkeypatching, testcontainers, and increased coverage, providing relevant unit and integration-testing evidence.
- SQL Server data-warehouse and SSIS experience demonstrates relational data modeling and ETL familiarity.
- Strong systems and performance background through C, Go, and C# projects, including profiling, benchmarking, memory management, and optimization research.
- Relevant academic and professional experience collaborating with senior engineers and maintaining internal services.

Gaps:
- The CV does not provide evidence of GraphQL, Graphene, or Apollo experience, which is central to the role’s backend stack.
- The CV does not provide evidence of PostgreSQL specifically; database experience is shown primarily with Microsoft SQL Server and SSIS.
- The CV does not provide evidence of AWS services such as Lambda, ECS, or DynamoDB, Redis, or third-party integrations involving payments, accounting, or CRM systems.
- The role requests 2–5 years of relevant experience and product-driven ownership; the CV currently shows a short internship and graduate program rather than that level of professional backend experience.
- The CV does not clearly demonstrate scalable production systems serving large user populations, customer-facing product development, architecture documentation, or product scoping.
- The stated English level of B2 may be adequate, but the CV provides limited evidence of customer-facing communication or close collaboration with product and design teams.

CV Feedback:
Position the Wargaming internship first and make its backend relevance explicit: emphasize Python library modernization, CLI ownership, testing, Docker/Kubernetes, CI, code reviews, and work in large internal codebases. Clarify which work was production-facing and describe the impact only where already supported by the CV. Move the most relevant technologies—Python, Flask, SQL, Docker, Kubernetes, Git, and Linux—into a prominent technical summary. Present the Endava data-engineering experience as evidence of relational modeling and data workflows, while clearly labeling SQL Server rather than implying PostgreSQL. Keep the interpreter, parser, and optimization projects as supporting evidence of software design and problem solving, but reduce emphasis on game-specific details unless space permits. Add GraphQL, AWS, Redis, third-party integrations, and customer-facing ownership only if they are genuinely part of the candidate’s experience; otherwise leave them as clear gaps rather than trying to imply equivalence. Ensure dates and current status are unambiguous, particularly the internship ending in February 2026 and the master’s program in progress.
```

## Result

This example demonstrates the intended separation of responsibilities in JobFind:

```text
Semantic retrieval
    |
    V
finds relevant job candidates

Recommendation Agent
    |
    V
evaluates and ranks them against the full CV

User
    |
    V
selects a preferred job

CV Tailor Agent
    |
    V
performs detailed CV-to-job analysis
```

Candidate preferences influence which jobs are recommended, but they are not treated as evidence of skills or experience.

The CV Tailor Agent also avoids inventing missing qualifications. If a requirement is not supported by the CV, it is reported as missing evidence rather than being added to the candidate profile.
