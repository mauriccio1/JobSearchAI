package config

import (
	"jobsearch/internal/parser"
	"os"
	"sort"
	"strings"
)

type Config struct {
    Name      string
    Contact   string
    Certs     string
    Education []string
	Port      string
	Skills    []parser.Skill
}

func Load() *Config {
	environment := os.Environ()
	return &Config{
        Name:    os.Getenv("NAME"),
        Contact: os.Getenv("CONTACT"),
        Certs:   os.Getenv("CERTS"),
        Education: loadVariableHelper("EDUCATION_", environment),
		Port: os.Getenv("PORT"),
		Skills: loadSkillsHelper(loadVariableHelper("SKILL_", environment)),
    }
}


func loadVariableHelper(field string, environment []string) ([]string){
	var keys []string
	var fieldValueSlice []string
	for _, env := range environment {
		parts := strings.SplitN(env, "=", 2)
		key := strings.ToUpper(parts[0])

		if strings.Contains(key, strings.ToUpper(field)) {	
			keys = append(keys, key)
		}
		
	}

	sort.Strings(keys)

	for _, k := range keys {
		fieldValueSlice =  append(fieldValueSlice, os.Getenv(k))
	}


	return  fieldValueSlice
	
}

func loadSkillsHelper(skills []string) []parser.Skill {
	var skillSlice []parser.Skill 
	for _, s := range skills {
		skill := parser.Skill{}
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			continue
		}

		skill.Label = strings.TrimSpace(parts[0])
		skill.Value = strings.TrimSpace(parts[1])

		skillSlice = append(skillSlice, skill)
	}

	return skillSlice
}


