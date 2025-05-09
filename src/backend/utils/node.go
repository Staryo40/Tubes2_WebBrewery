package utils

import(
	"backend/models"
)

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