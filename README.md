# JobSearchAI

A local AI-powered resume tailoring pipeline built in Go. Paste a job description, get back a tailored PDF — all running on your own machine with no cloud API costs.

Built out of necessity after experiencing firsthand how broken the job application process is.

---

## How It Works

```
Job Description (raw page text)
           │
           ▼
    [Llama 3.2 3B]              ← strips noise, extracts requirements only
           │
           ▼
      Trimmed JD
           │
           ▼
    [Gemma 2 9B]                ← extracts structured experience from base resume
           │
           ├─────────────────────────────────┐
           ▼                                 ▼
    RankBullets (Gemma 2 9B)      RewordBullets + RewriteSummary (Gemma 2 9B)
    returns ranked indices         returns reworded bullets and tailored summary
           │                                 │
           └──────────────┬──────────────────┘
                          ▼
               Go Assembly Layer              ← deterministic, no LLM involved
               reorders bullets by rank,
               merges reworded text
                          │
                          ▼
                    PDF Generation
```

**Why decomposed?**
A single LLM call asked to simultaneously rank, reword, preserve, and format all experience bullets under one prompt was too much — the model dropped bullets, ignored ordering, and played it safe. Breaking it into focused single-task calls with deterministic Go assembly in between produces significantly more consistent output.

**Why two models?**
- Llama 3.2 3B is fast for extraction — strip noise from a raw JD with no creativity needed
- Gemma 2 9B handles the judgment tasks — rewording, ranking, summary framing

**Why local?**
Zero per-request cost. Your resume and job data never leave your machine.

---

## Stack

- **Go** — HTTP server, pipeline orchestration, Ollama API client, PDF generation
- **Cobra** — CLI with `serve` and `setup` commands
- **Ollama** — local LLM runtime
- **Llama 3.2 3B** — job description trimmer
- **Gemma 2 9B** — experience extractor, bullet ranker, bullet rewriter, summary writer
- **Chromium Extension** — extracts JD from current browser tab, sends to local server, returns tailored PDF for preview, download, or upload

---

## Project Structure

```
JobSearchAI/
├── main.go                        # one-liner entry point
├── go.mod
├── command/
│   ├── root.go                    # Cobra root command
│   ├── serve.go                   # go run main.go serve
│   └── setup.go                   # go run main.go setup
├── internal/
│   ├── server/
│   │   └── server.go              # HTTP server, route handlers, graceful shutdown
│   ├── models/
│   │   ├── models.go              # pipeline orchestration, Ollama client
│   │   ├── prompts.go             # all LLM prompts
│   │   └── schemas.go             # JSON schemas for structured outputs
│   ├── config/
│   │   └── config.go              # loads .env, dynamic SKILL_N / EDUCATION_N
│   └── parser/
│       └── parser.go              # Resume, Job, JobSection, Skill types
├── extension/                     # Chromium extension (Chrome + Edge)
└── resume/                        # gitignored — your personal resume files live here
    └── base_resume.txt
```

---

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Ollama](https://ollama.com) installed and running
- Required models:

```bash
ollama pull llama3.2:3b
ollama pull gemma2:9b
```

---

## Setup

**1. Clone the repo**
```bash
git clone https://github.com/mauriccio1/JobSearchAI.git
cd JobSearchAI
```

**2. Add your base resume**
```bash
mkdir resume
vim resume/base_resume.txt   # paste your resume as plain text
```

**3. Run setup**

The setup command reads your base resume, parses it with an LLM, and writes a `.env` file with your name, contact, skills, education, and certifications.

```bash
go run main.go setup
```

Review the generated `.env` before confirming. Setup takes ~30–40 seconds.

**4. Start Ollama**
```bash
ollama serve
```

**5. Start the server**
```bash
go run main.go serve
```

Server starts on the port defined in your `.env` (default `8080`).

---

## Extension

**Install**

1. Open Chrome or Edge and navigate to `chrome://extensions`
2. Enable **Developer mode** (top right toggle)
3. Click **Load unpacked** and select the `extension/` folder

**Usage**

1. Navigate to any job posting page in your browser
2. Click the JobSearchAI extension icon
3. Verify the server status shows **CONNECTION ESTABLISHED** — if not, make sure `go run main.go serve` is running
4. Click **GENERATE RESUME**
5. The extension reads the job description from the current page, sends it to your local server, and runs the full pipeline
6. When complete, you can:
   - **PREVIEW** — opens the tailored PDF in a new tab
   - **DOWNLOAD** — saves the PDF with your chosen filename
   - **UPLOAD TO PAGE** — attempts to inject the PDF into a file upload field on the current page (works on some ATS platforms, not all)

**Note:** The extension cannot read certain page types (PDFs, browser internal pages, some authenticated portals). If the server is running but generation fails, try copying the job description manually and using the API directly.

---

## API

### `POST /api/resume/rewrite`

Trims a job description and runs the full tailoring pipeline against your base resume.

**Request**
```json
{ "jd": "full job description text" }
```

**Response**

Returns a PDF binary (`application/pdf`).

```bash
curl -X POST http://localhost:8080/api/resume/rewrite \
  -H "Content-Type: application/json" \
  -d '{"jd": "your job description here"}' \
  --output resume.pdf
```

---

## Design Decisions

**Decomposed pipeline over monolithic rewrite** — splitting rewrite into three concurrent calls (RankBullets, RewordBullets, RewriteSummary) with Go assembly in between gives each model a single focused task. Bullet preservation is handled deterministically in code, not by the LLM. If any model output fails validation (wrong index count, out-of-bounds index, missing company), the pipeline falls back gracefully to the original order.

**Anti-fabrication guardrails** — after testing 6+ models and iterating through multiple prompt architectures, the final prompts distinguish between framing (describing real work in JD language) and fabrication (inventing tools, metrics, or employers). The reword prompt uses concrete before/after examples to teach this distinction. The summary prompt enforces exact numbers from the resume rather than calculated percentages.

**Two-step JD trimming** — passing raw job description page text directly to the rewriter caused context overflow and noise. Pre-trimming with Llama 3.2 3B keeps the rewriter's input clean and within the context window regardless of JD length.

**Structured JSON output** — all model calls use Ollama's `format` field with explicit JSON schemas. Wrapper structs (`expWrapper`, `rankWrapper`) normalize model output before it reaches the assembly layer.

---

## Roadmap

### V1 — complete
- [x] Decomposed LLM pipeline (trim → extract → rank + reword → assemble)
- [x] PDF generation
- [x] Cobra CLI (`serve`, `setup`)
- [x] LLM-powered onboarding writes `.env` from base resume
- [x] Chromium extension (Chrome + Edge) with preview, download, and upload
- [x] Anti-hallucination prompt architecture

### V2 — planned
- [ ] Pluggable provider interface — `--provider=ollama|claude|openai`, `--model=<any>`
- [ ] Cloud API provider (Claude/OpenAI) for higher quality output on demand
- [ ] Cloud deployment — GCP + Kubernetes, GPU-backed inference nodes
- [ ] Fine-tuning pipeline — generate high-quality resume pairs at scale, fine-tune a small model on curated data

---

## Why I Built This

After being laid off alongside my team, I started applying manually — tweaking resumes one by one, losing track of applications, getting no responses. A former manager who has deep experience in the industry needed nearly a thousand applications to land one role. The process is broken for everyone.

This tool is the first piece of a larger automated job application platform. The goal is to give candidates the same firepower as enterprise recruiting software, running locally and privately on their own machine.

---

## License

MIT
