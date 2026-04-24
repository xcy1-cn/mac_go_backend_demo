// package models

// type Product struct {
// 	ID    int     `json:"id"`
// 	Name  string  `json:"name"`
// 	Price float64 `json:"price"`
// 	Stock int     `json:"stock"`
// }

package models

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"gt=0"`
	Stock int     `json:"stock" binding:"gte=0"`
}
