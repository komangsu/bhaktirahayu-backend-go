package schemas

import (
	"time"
)

type (
	LogisticsSchema struct {
		Id           int64     `json:"id"`
		LogisticType int64     `json:"logistic_type" binding:"required"`
		Name         string    `json:"name" binding:"required,min=3,max=100"`
		Stock        int       `json:"stock"`
		Component    string    `json:"component" binding:"required"`
		MaxPrice     int64     `json:"max_price" binding:"required"`
		Price        int64     `json:"price" binding:"required"`
		ExpiredDate  string    `json:"expired_date" binding:"required"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	UpdateLogisticsSchema struct {
		LogisticType int64     `json:"logistic_type" binding:"required"`
		Name         string    `json:"name" binding:"required,min=3,max=100"`
		Stock        int       `json:"stock" binding:"required"`
		Component    string    `json:"component" binding:"required"`
		MaxPrice     int64     `json:"max_price" binding:"required"`
		Price        int64     `json:"price" binding:"required"`
		ExpiredDate  string    `json:"expired_date" binding:"required"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	ListAllLogistics struct {
		Id           int64     `json:"id"`
		LogisticType int64     `json:"logistic_type"`
		Name         string    `json:"name"`
		Stock        *int      `json:"stock"`
		Component    string    `json:"component"`
		MaxPrice     int64     `json:"max_price"`
		Price        int64     `json:"price"`
		ExpiredDate  string    `json:"expired_date"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	PaginationLogistic struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
		Total   int `json:"total"`
	}

	// logistic services
	LogisticServices struct {
		Id     int64  `json:"id"`
		Name   string `json:"name" binding:"required,min=3,max=100"`
		Margin int64  `json:"margin" binding:"required"`
	}

	ListAllLogisticServices struct {
		Id     int64  `json:"id"`
		Name   string `json:"name"`
		Margin int64  `json:"margin" binding:"required"`
	}

	QueryPage struct {
		Page    string `form:"page" binding:"required"`
		PerPage string `form:"per_page" binding:"required"`
	}
)
