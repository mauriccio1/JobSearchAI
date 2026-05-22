package parser

type Resume struct {
    Name       string  `json:"name"`
    Contact    string  `json:"contact"`
    Summary    string  `json:"summary"`
    Experience []Job   `json:"experience"`
    Skills     []Skill `json:"skills"`
    Certs      string  `json:"certs"`
    Education  []string `json:"education"`
}

type Job struct {
    Title    string       `json:"title"`
    Company  string       `json:"company"`
    Dates    string       `json:"dates"`
    Intro    string       `json:"intro"`
    Sections []JobSection `json:"sections"`
}

type JobSection struct {
    Header  string   `json:"header"`
    Bullets []string `json:"bullets"`
	BulletOrder  []int  `json:"order"`
}

type Skill struct {
    Label string `json:"label"`
    Value string `json:"value"`
}
