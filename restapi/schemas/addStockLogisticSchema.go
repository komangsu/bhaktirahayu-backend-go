package schemas

type AddStockSchema struct {
	// Institutions string `json:"institutions" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Stock int    `json:"stock" binding:"required"`
	Price int64  `json:"price" binding:"required"`
}
