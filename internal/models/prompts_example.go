package models

// prompts_example.go
//
// ============================================================
// SETUP INSTRUCTIONS
// ============================================================
//
// 1. Copy this file to prompts.go in the same directory:
//    cp internal/models/prompts_example.go internal/models/prompts.go
//
// 2. In prompts.go, uncomment all the code (remove the // prefix from every line)
//
// 3. Replace all values marked with [BRACKETS] with your actual resume information
//
// 4. prompts.go is gitignored and will never be committed — your personal
//    resume details stay on your machine only
//
// ============================================================
// TIPS
// ============================================================
//
// EMPLOYERS: List every job exactly as it appears on your resume including
// the exact title, company name, and date range. The LLM will use this
// to prevent hallucination — it won't invent jobs or move bullets between
// employers if you list them explicitly here.
//
// Example:
//   a) Senior Software Engineer | Acme Corp | Jan 2023 – Present
//   b) Software Engineer | Startup Inc | Jun 2021 – Dec 2022
//
// TECHNICAL SKILLS: Copy your skills section exactly from your resume.
// Group them by category the same way your resume does. The LLM will
// reorder categories by relevance to each job description but won't
// add or remove any skills.
//
// Example:
//   Languages: Python (primary), Go, TypeScript
//   Cloud & Platform: AWS (EC2, S3, Lambda, RDS), Docker, Kubernetes
//
// CERTIFICATIONS: List all certifications comma separated exactly as
// they appear on your resume. The LLM will print these verbatim.
//
// Example:
//   AWS Solutions Architect Associate, CompTIA Security+, CKAD
//
// EDUCATION: List each degree or certificate on its own line exactly
// as it appears on your resume.
//
// Example:
//   B.S. Computer Science, University of Florida
//   College Certificate, Cybersecurity, Miami Dade College
//
// NUMBER OF EMPLOYERS: Update the "Keep ALL [NUMBER] employers" line
// to match how many jobs you have listed. If you have 3 jobs write
// "Keep ALL 3 employers", if you have 5 write "Keep ALL 5 employers".
// This helps the LLM know exactly how many entries to expect.
//
// ============================================================

// package models

// const (
// 	resumeRewriteSysPrompt = "You are a resume rewriting engine. You receive a resume and a job description. You rewrite the resume to better match the job description. You ONLY reword and reorder what already exists. You NEVER add or remove information. Every skill, certification, employer, bullet point, and education entry from the original must appear in your output. Your output is always a complete resume with no commentary."

// 	resumeRewritePrompt = `Rewrite my resume to match the job description. Follow these rules exactly.
// RULES:
//
// OUTPUT: Only the final resume. No commentary, no analysis, no tips, no questions.
// STRUCTURE — use this exact section order:
// SUMMARY
// EXPERIENCE
// TECHNICAL SKILLS
// CERTIFICATIONS
// EDUCATION
// EXPERIENCE RULES:
// Keep ALL [NUMBER] employers in this exact order:
// a) [JOB TITLE] | [COMPANY NAME] | [START DATE] – [END DATE]
// b) [JOB TITLE] | [COMPANY NAME] | [START DATE] – [END DATE]
// c) [JOB TITLE] | [COMPANY NAME] | [START DATE] – [END DATE]
// d) [JOB TITLE] | [COMPANY NAME] | [START DATE] – [END DATE]
// Do NOT merge, rename, or remove any employer.
// Do NOT move bullet points between employers.
// You may reword bullets to use keywords from the job description.
// You may reorder bullets within each employer to put the most relevant first.
// SUMMARY: Rewrite using keywords from the job description. Use ONLY skills and experience from my original resume.
// TECHNICAL SKILLS — include ALL of the following, reordered for relevance to the JD:
// [SKILL CATEGORY]: [SKILL 1], [SKILL 2], [SKILL 3].
// [SKILL CATEGORY]: [SKILL 1], [SKILL 2], [SKILL 3].
// [SKILL CATEGORY]: [SKILL 1], [SKILL 2], [SKILL 3].
// [SKILL CATEGORY]: [SKILL 1], [SKILL 2], [SKILL 3].
// [SKILL CATEGORY]: [SKILL 1], [SKILL 2], [SKILL 3].
// [SKILL CATEGORY]: [SKILL 1], [SKILL 2], [SKILL 3].
// CERTIFICATIONS — print exactly:
// [CERT 1], [CERT 2], [CERT 3]
// EDUCATION — print exactly:
// [DEGREE], [INSTITUTION]
// [DEGREE], [INSTITUTION]
// DO NOT ADD anything not in the original resume. No new skills, tools, certifications, degrees, or employers.
// DO NOT REMOVE anything from the original resume. Every item must appear.
// ---RESUME---
// `
// ) 
