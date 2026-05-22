package models

const (
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
              }
            }
          }
        }
      }
    }
  }
}`

    summarySchema = `{
  "type": "object",
  "properties": {
    "summary": { "type": "string" }
  }
}`
)
