import requests
from bs4 import BeautifulSoup
import json

url = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
response = requests.get(url)
soup = BeautifulSoup(response.content, "html.parser")

data = []

# Find all h3 headers
for h3 in soup.find_all("h3"):
    span = h3.find("span", class_="mw-headline")
    if not span:
        continue

    category = span.get_text(strip=True)

    # Find the first table after this h3
    table = h3.find_next("table")
    if not table:
        continue

    tbody = table.find("tbody")
    if not tbody:
        continue
    
    # Header
    rows = tbody.find_all("tr")[1:]  

    for row in rows:
        tds = row.find_all("td")
        if len(tds) != 2:
            continue

        # First TD: product
        product_td = tds[0]
        a_tags = product_td.find_all("a")
        product = a_tags[-1].get_text(strip=True) if a_tags else product_td.get_text(strip=True)

        # Second TD: recipe(s)
        second_td = tds[1]
        ul = second_td.find("ul")

        if ul:
            recipes = []
            for li in ul.find_all("li"):
                # Split ingredients by "+" and strip whitespace
                ingredients = [ing.strip() for ing in li.get_text().split("+")]
                recipes.append(ingredients)
        else:
            # Handle plain text recipe (not a list)
            text = second_td.get_text(strip=True)
            if text:
                ingredients = [ing.strip() for ing in text.split("+")]
                recipes = [ingredients]
            else:
                recipes = []

        data.append({
            "product": product,
            "category": category,
            "recipes": recipes
        })

# Output to JSON
with open("recipes.json", "w", encoding="utf-8") as f:
    json.dump(data, f, ensure_ascii=False, indent=2)

print("Scraping complete! Saved to recipes.json.")
