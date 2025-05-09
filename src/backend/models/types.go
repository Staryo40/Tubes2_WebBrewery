package models

type Element struct {
    Name     string     `json:"name"`
    Category string     `json:"category"`
    Recipes  [][]string `json:"recipes"` 
    Image    string     `json:"image"` 
}

type Node struct {
    Name string
    Ingredient1 string
    Ingredient2 string
}