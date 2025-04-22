package scrape

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RecipeEntry struct {
	Product  string     `json:"product"`
	Category string     `json:"category"`
	Recipes  [][]string `json:"recipes"`
}

func Scrape() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Non-200 response: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	var results []RecipeEntry

	doc.Find("h3").Each(func(i int, h3 *goquery.Selection) {
		category := h3.Find("span.mw-headline").Text()
		if category == "" {
			return
		}

		// Find the first <table> h3
		table := h3.NextUntil("table").NextFiltered("table")
		if table.Length() == 0 {
			return
		}

		table.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
			if i == 0 {
				return // skip header
			}

			tds := tr.Find("td")
			if tds.Length() != 2 {
				return
			}

			// First TD: product
			productTd := tds.Eq(0)
			links := productTd.Find("a")
			product := ""
			if links.Length() > 0 {
				product = strings.TrimSpace(links.Last().Text())
			} else {
				product = strings.TrimSpace(productTd.Text())
			}

			// Second TD: recipe(s)
			secondTd := tds.Eq(1)
			var recipes [][]string

			if ul := secondTd.Find("ul"); ul.Length() > 0 {
				ul.Find("li").Each(func(_ int, li *goquery.Selection) {
					text := li.Text()
					ingredients := splitAndClean(text)
					recipes = append(recipes, ingredients)
				})
			} else {
				text := strings.TrimSpace(secondTd.Text())
				if text != "" {
					recipes = append(recipes, splitAndClean(text))
				}
			}

			results = append(results, RecipeEntry{
				Product:  product,
				Category: category,
				Recipes:  recipes,
			})
		})
	})

	file, err := os.Create("recipes2.json")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	fmt.Println("Scraping complete! Saved to recipes2.json.")
}

func splitAndClean(s string) []string {
	parts := strings.Split(s, "+")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}
