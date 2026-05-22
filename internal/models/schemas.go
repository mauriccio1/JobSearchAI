package models

const (
	// Used by extractExp and RewordBullets — wraps []parser.Job in {"experience": [...]}
	// NOTE: unmarshal code in models.go must use a wrapper struct, e.g.:
	//   type expWrapper struct { Experience []parser.Job `json:"experience"` }
	//   var w expWrapper
	//   json.Unmarshal([]byte(response), &w)
	//   return w.Experience, nil
	experienceSchema = `{
  "type": "object",
  "properties": {
    "experience": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "title":   { "type": "string" },
          "company": { "type": "string" },
          "dates":   { "type": "string" },
          "intro":   { "type": "string" },
          "sections": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "header":  { "type": "string" },
                "bullets": { "type": "array", "items": { "type": "string" } }
              },
              "required": ["header", "bullets"]
            }
          }
        },
        "required": ["title", "company", "dates", "sections"]
      }
    }
  },
  "required": ["experience"]
}`

	// Used by RewordBullets — same structure as experienceSchema
	rewordBulletSchema = experienceSchema

	// Used by RankBullets — wraps []SectionRanking in {"rankings": [...]}
	// NOTE: unmarshal code in models.go must use a wrapper struct, e.g.:
	//   type rankWrapper struct { Rankings []SectionRanking `json:"rankings"` }
	//   var w rankWrapper
	//   json.Unmarshal([]byte(response), &w)
	//   return w.Rankings, nil
	rankSchema = `{
  "type": "object",
  "properties": {
    "rankings": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "company": { "type": "string" },
          "section": { "type": "string" },
          "order":   { "type": "array", "items": { "type": "integer" } }
        },
        "required": ["company", "section", "order"]
      }
    }
  },
  "required": ["rankings"]
}`

	summarySchema = `{
  "type": "object",
  "properties": {
    "summary": { "type": "string" }
  }
}`
)
