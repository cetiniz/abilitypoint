package model

type Skill struct {
	Resources   []interface{}
	Images      []interface{}
	Name        string
	Description string
}

type Edge struct {
	From Skill
	Name string
	To   Skill
}
