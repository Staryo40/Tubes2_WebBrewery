package utils

import(
	"backend/models"
)

func NodeEqual(a models.Node, b models.Node) bool {
	return a.Name == b.Name && a.Ingredient1 == b.Ingredient2 && a.Ingredient2 == b.Ingredient2
}

func ConvertToElementList(elements map[string]models.Element) []models.ElementEntry {
	list := make([]models.ElementEntry, 0, len(elements))
	for name, element := range elements {
		entry := models.ElementEntry{
			Name:    name,
			Element: element,
		}
		list = append(list, entry)
	}
	return list
}

func ConvertToElementMap(entries []models.ElementEntry) map[string]models.Element {
	elementMap := make(map[string]models.Element, len(entries))
	for _, entry := range entries {
		elementMap[entry.Name] = entry.Element
	}
	return elementMap
}

func PathsEqual(a, b []models.Node) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}