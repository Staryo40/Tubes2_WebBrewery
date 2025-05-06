package scrape

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"io"
	"path/filepath"
	"encoding/base64"
	"errors"

	"github.com/PuerkitoBio/goquery"
)

type RecipeEntry struct {
	Product  string     `json:"product"`
	Category string     `json:"category"`
	Recipes  [][]string `json:"recipes"`
}

type Element struct {
    Name     string     `json:"name"`
    Category string     `json:"category"`
    Recipes  [][]string `json:"recipes"` 
    Image    string     `json:"image"` 
}

// SCRAPE ONLY RECIPES (USING RECIPEENTRY STRUCT)
func ScrapeRecipe() {
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

// SCRAPE WITH IMAGES (USING ELEMENT STRUCT)
func ScrapeElementsAndImages() {
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

	var results []Element
	doc.Find("h3").Each(func(i int, h3 *goquery.Selection) {
		category := h3.Find("span.mw-headline").Text()
		if category == "" {
			return
		}

		table := h3.NextUntil("table").NextFiltered("table")
		if table.Length() == 0 {
			return
		}

		table.Find("tbody tr").Each(func(i int, tr *goquery.Selection) {
			if i == 0 {
				return
			}

			tds := tr.Find("td")
			if tds.Length() != 2 {
				return
			}

			// Product name
			productTd := tds.Eq(0)
			product := strings.TrimSpace(productTd.Find("a").Last().Text())

			imgTag := productTd.Find("span.icon-hover span[typeof='mw:File'] a.mw-file-description img")
			imgURL, exists := imgTag.Attr("data-src")
			if !exists {
				imgURL, exists = imgTag.Attr("src")
			}
			if !exists || product == "" || strings.HasPrefix(imgURL, "data:") {
				fmt.Printf("Skipping image for: %s\n", product)
				return
			}

			imgURL = cleanImageURL(imgURL)
			ext := filepath.Ext(imgURL)
			imageName := product + ext
			imagePath := "../src/frontend/public/images/" + imageName
			err := downloadImage(imgURL, imagePath)
			if err != nil {
				fmt.Printf("Failed to download image for %s: %v\n", product, err)
			}

			// Recipe parsing
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

			fmt.Printf("Appending: %s\n", product)

			results = append(results, Element{
				Name:     product,
				Category: category,
				Recipes:  recipes,
				Image:    "/images/" + imageName,
			})
		})
	})

	for i, el := range results {
		if i >= 5 {
			break
		}
		fmt.Printf("%+v\n", el)
	}

	file, err := os.Create("elements.json")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	fmt.Println("Scraping complete! Saved to elements.json.")
}


// HELPER FUNCTIONS
func splitAndClean(s string) []string {
	parts := strings.Split(s, "+")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

func cleanImageURL(rawURL string) string {
    i := strings.Index(rawURL, ".svg")
    if i == -1 {
        i = strings.Index(rawURL, ".png")
    }
    if i == -1 {
        return rawURL // fallback if nothing matched
    }
    return rawURL[:i+4] // include the extension
}

func downloadImage(url, filepathStr string) error {
    if strings.HasPrefix(url, "//") {
        url = "https:" + url
    }

    // Handle base64 data URI
    if strings.HasPrefix(url, "data:") {
        return saveBase64Image(url, filepathStr)
    }

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if err := os.MkdirAll(filepath.Dir(filepathStr), 0755); err != nil {
        return err
    }

    file, err := os.Create(filepathStr)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = io.Copy(file, resp.Body)
    return err
}

func saveBase64Image(dataURI, filepathStr string) error {
    // Example: data:image/png;base64,ABCDEF==
    parts := strings.SplitN(dataURI, ",", 2)
    if len(parts) != 2 {
        return errors.New("invalid data URI format")
    }

    data := parts[1]
    decoded, err := base64.StdEncoding.DecodeString(data)
    if err != nil {
        return err
    }

    if err := os.MkdirAll(filepath.Dir(filepathStr), 0755); err != nil {
        return err
    }

    return os.WriteFile(filepathStr, decoded, 0644)
}