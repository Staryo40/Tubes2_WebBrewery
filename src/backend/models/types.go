package models

type Element struct {
    Name     string     `json:"name"`
    Category string     `json:"category"`
    Recipes  [][]string `json:"recipes"` 
    Image    string     `json:"image"` 
}

type ElementEntry struct {
	Name    string
	Element Element
}

type Node struct {
    Name string
    Ingredient1 string
    Ingredient2 string
}

type RequestPayload struct {
    Target string 		`json:"target"`
    Method string 		`json:"method"`
    PathNumber  int 	`json:"pathNumber"`  
	Bidirectional bool 	`json:"bidirectional"`
}

type ResponsePayload struct {
    Count     int           `json:"count"`   
    ElapsedTime int64       `json:"elapsedTime"` 
    Paths     []CountedPath `json:"paths"`     
}

type CountedPath struct {
    NodeCount int       `json:"nodeCount"`
    Path     []Node     `json:"path"` 
}

type ErrorResponse struct {
	Error   string         `json:"error"`
	Request RequestPayload `json:"request"`
}

