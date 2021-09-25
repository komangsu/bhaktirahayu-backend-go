package controllers

import (
	"log"
	"math"
	"net/http"
	"restapi/models"
	"restapi/schemas"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateLogistics(c *gin.Context) {
	var payload schemas.LogisticsSchema

	// validation
	if err := c.ShouldBindWith(&payload, binding.JSON); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	// check logistic type
	log_type := models.CheckLogisticType(payload.LogisticType)
	if log_type != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "type logistic not found."})
		return
	}

	// check logistic name
	name := models.CheckLogisticDuplicateName(payload.Name)
	if name != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "The name has already been taken."})
		return
	}

	// check the price must be below the maximum price
	if payload.Price > payload.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "The price cannot be more than the maximum price."})
		return
	}

	// check expired date format
	layout := "02-01-2006"
	_, err := time.Parse(layout, payload.ExpiredDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid date format"})
		return
	}

	// save logistics
	models.CreateLogistic(payload)
	c.JSON(http.StatusCreated, gin.H{"detail": "Successfully add a new logistic."})

}

func UpdateLogistics(c *gin.Context) {
	var payload schemas.UpdateLogisticsSchema

	// validation
	if err := c.ShouldBindWith(&payload, binding.JSON); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	// get id & convert to integer
	param_id := c.Param("id")
	logistic_id, err_convert := strconv.Atoi(param_id)
	if err_convert != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "value is not a valid integer"})
		return
	}

	// check id is exist or not
	id := models.CheckLogisticById(logistic_id)
	if id != 1 {
		c.JSON(http.StatusNotFound, gin.H{"detail": "logistics not found."})
		return
	}

	// check logistic type
	log_type := models.CheckLogisticType(payload.LogisticType)
	if log_type != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "type logistic not found."})
		return
	}

	// check duplicate name
	name := models.CheckLogisticDuplicateName(payload.Name)
	if name != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "The name has already been taken."})
		return
	}

	// check the price must be below the maximum price
	if payload.Price > payload.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "The price cannot be more than the maximum price."})
		return
	}

	// check expired date format
	layout := "02-01-2006"
	_, err := time.Parse(layout, payload.ExpiredDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid date format"})
		return
	}

	// update logistics
	log.Println(payload)
	models.UpdateLogistic(logistic_id, payload)
	c.JSON(http.StatusOK, gin.H{"detail": "Successfully update the logistics."})
}

func DeleteLogistics(c *gin.Context) {
	// get id & convert to integer
	param_id := c.Param("id")
	logistics_id, err_convert := strconv.Atoi(param_id)
	if err_convert != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "value is not a valid integer"})
		return
	}

	// check id is exist or not
	id := models.CheckLogisticById(logistics_id)
	if id != 1 {
		c.JSON(http.StatusNotFound, gin.H{"detail": "logistics not found."})
		return
	}

	// delete logistics
	models.DeleteLogistic(logistics_id)
	c.JSON(http.StatusOK, gin.H{"detail": "Successfully delete the logistics."})
}

func GetAllLogistics(c *gin.Context) {
	var (
		pagination schemas.PaginationLogistic
		payload    schemas.QueryPage
	)

	// validation
	if err := c.ShouldBind(&payload); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	// get query
	q := c.Query("q")
	query_page := c.Query("page")
	query_ppage := c.Query("per_page")

	q = strings.ToLower(q)

	page_value, err_page := strconv.Atoi(query_page)
	pagination.Page = page_value
	if err_page != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "value is not a valid integer"})
		return
	}

	ppage_value, err_ppage := strconv.Atoi(query_ppage)
	pagination.PerPage = ppage_value
	if err_ppage != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "value is not a valid integer"})
		return
	}

	logistics, total := models.GetAllLogistic(q, pagination)

	// pages
	pages := pagesLogistics(pagination.PerPage, total)

	// prev_num
	has_prev := hasPrevLogistics(pagination.Page)
	prev_num := prevNumLogistics(has_prev, pagination.Page)

	// next_num
	has_next := hasNextLogistics(pagination.Page, pages)
	next_num := nextNumLogistics(has_next, pagination.Page)

	// iter_pages
	ledge := 2
	lcurrent := 2
	rcurrent := 5
	redge := 2
	iter := iterPagesLogistics(pagination.Page, pages, ledge, lcurrent, rcurrent, redge)

	c.JSON(http.StatusOK, gin.H{
		"data":       logistics,
		"total":      total,
		"prev_num":   prev_num,
		"next_num":   next_num,
		"page":       pagination.Page,
		"iter_pages": iter,
	})
}

func pagesLogistics(per_page int, total int) int {
	var page int

	if per_page == 0 || total == 0 {
		page = 0
	} else {
		page = int(math.Ceil(float64(total) / float64(per_page)))
	}

	return page
}

func prevNumLogistics(has_prev bool, page int) *int {
	if !has_prev {
		return nil
	}
	page = page - 1

	return &page
}

func hasPrevLogistics(page int) bool {
	// true if previous page exists
	return page > 1
}

func hasNextLogistics(page int, pages int) bool {
	// true if next page exists
	return page < pages
}

func nextNumLogistics(has_next bool, page int) *int {
	if !has_next {
		return nil
	}
	page = page + 1

	return &page
}

func iterPagesLogistics(page int, pages int, left_edge int, left_current int, right_current int, right_edge int) []*int {
	last := 0

	var p = []*int{}

	for i := 1; i <= pages; i++ {
		if i <= left_edge || (i > page-left_current-1 && i < page+right_current) || i > pages-right_edge {
			if last+1 != i {
				p = append(p, nil)
			}
			last = i
			a := &last
			b := *a
			p = append(p, &b)
		}
	}
	return p
}
