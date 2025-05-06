import json

# Load the JSON file into a Python variable
with open("recipes.json", "r", encoding="utf-8") as f:
    recipes_data = json.load(f)

# Basic Analysis
total_products = len(recipes_data)
all_product_names = [entry["product"] for entry in recipes_data]
all_categories = list(set(entry["category"] for entry in recipes_data))

print(f"Total products: {total_products}")
print(f"Unique categories: {len(all_categories)}")
# print("First 5 products:", all_product_names[:5])
# print("All products:", all_product_names)

recipe_count = 0
for item in recipes_data:
    for recipe in item.get('recipes', []):
        if isinstance(recipe, list) and len(recipe) == 2:
            recipe_count += 1

print("Total number of valid recipes:", recipe_count)
