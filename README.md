# JobSearchAI

A local AI-powered resume optimization pipeline built in Go. Given a job description, it automatically trims the noise, rewrites your base resume to match the role using locally running LLMs, and returns a tailored result — all running on your own machine with no cloud API costs.

Built out of necessity after experiencing firsthand how broken the job application process is.

---

## How It Works

The pipeline runs two LLM inference steps sequentially against a local [Ollama](https://ollama.com) server:

```
Job Description (raw)
        │
        ▼
  [Llama 3.2 3B]          ← strips noise, extracts requirements only
        │
        ▼
  Trimmed JD
        │
        ▼
  [Gemma 2 9B]            ← rewrites base resume to match JD keywords
        │
        ▼
  Tailored Resume (JSON)
```

**Why two models?**
- Llama 3.2 3B is fast and cheap for extraction tasks — no creativity needed, just signal from noise
- Gemma 2 9B has strong instruction-following for the rewrite — it reorganizes and reframes without fabricating experience

**Why local models?**
At scale, cloud API costs for resume rewriting become significant. Running inference locally means zero per-request cost regardless of volume.

---

## Stack

- **Go** — HTTP server, pipeline orchestration, Ollama API client
- **Ollama** — local LLM runtime
- **Llama 3.2 3B** — job description trimmer
- **Gemma 2 9B** — resume rewriter
- **Chrome Extension** *(in progress)* — extracts JD from current browser tab, sends to local server, returns tailored resume for review and upload

---

## Project Structure

```
JobSearchAI/
├── main.go                   # Entry point, signal handling
├── go.mod
├── internal/
│   ├── server/
│   │   └── server.go         # HTTP server, route handlers, graceful shutdown
│   └── models/
│       └── models.go         # Ollama API client, TrimJD, RewriteResume
└── resume/                   # gitignored — put your base_resume.txt here
    └── base_resume.txt
```

---

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Ollama](https://ollama.com) installed and running
- Required models pulled:

```bash
ollama pull llama3.2:3b
ollama pull gemma2:9b
```

---

## Setup

**1. Clone the repo**
```bash
git clone https://github.com/yourusername/JobSearchAI.git
cd JobSearchAI
```

**2. Add your base resume**
```bash
mkdir resume
# paste your resume as plain text
vim resume/base_resume.txt
```

**3. Start Ollama**
```bash
ollama serve
```

**4. Run the server**
```bash
go run main.go
```

Server starts on `localhost:8080`.

---

## API

### `POST /api/resume/rewrite`

Trims a job description and rewrites your base resume to match it.

**Request**
```json
{
  "jd": "paste the full job description text here"
}
```

**Response**
```json
{
  "resume": "rewritten resume text tailored to the job description"
}
```

**Example**
```bash
curl -X POST http://localhost:8080/api/resume/rewrite \
  -H "Content-Type: application/json" \
  -d '{"jd": "your job description here"}'
```

---

## Design Decisions

**Prompt engineering over fine-tuning** — after testing 6+ models (Phi4-Mini, Mistral 7B, Qwen3 4B, Llama 3.2 3B, Qwen3 14B, Gemma 2 9B) with an iterative prompt refinement process, Gemma 2 9B with explicit anti-fabrication guardrails produced the most consistent, accurate rewrites. The key failure modes discovered were fabricated credentials, shuffled work history between employers, and repetition loops on long inputs — all addressed in the final prompt design.

**Two-step pipeline** — passing raw job descriptions directly to the rewriter caused context overflow on longer JDs. Pre-trimming with a fast 3B model keeps the rewriter's input well within the 8K context window regardless of JD length.

**Graceful shutdown** — the server listens for SIGINT/SIGTERM and gives in-flight requests 5 seconds to complete before exiting.

---

## Roadmap

- [ ] Chrome extension — extract JD from current tab, preview tailored resume, upload directly to job page
- [ ] PDF generation — render rewritten resume as a clean, ATS-friendly PDF
- [ ] Application tracker — SQLite-backed dashboard to track applications, statuses, and generated resumes
- [ ] Job crawler — scrape job boards and filter by match score before rewriting
- [ ] Email monitor — Gmail API integration to surface recruiter responses in the dashboard
- [ ] Multi-user support — user profiles, resume management, application history

---

## Why I Built This

After being laid off alongside my team, I started applying manually — tweaking resumes one by one, losing track of applications, getting no responses. A former manager who has deep experience in the industry needed nearly a thousand applications to land one role. The process is broken for everyone.

This tool is the first piece of a larger automated job application platform. The goal is to give candidates the same firepower as enterprise recruiting software, running locally and privately on their own machine.

---

## License

MIT
