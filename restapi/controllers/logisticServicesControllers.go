package controllers

import (
	"math"
	"net/http"
	"restapi/models"
	"restapi/schemas"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func CreateLogisticService(c *gin.Context) {
	var payload schemas.LogisticServices

	// include access token

	// validation
	if err := c.ShouldBindWith(&payload, binding.JSON); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	// check duplicate name
	name := models.CheckLogServiceDuplicateName(payload.Name)
	if name != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "The name has already been taken."})
		return
	}

	// save
	_ = models.CreateLogisticService(payload)
	c.JSON(http.StatusCreated, gin.H{"detail": "Successfully add a new guardian."})
}

func UpdateLogisticService(c *gin.Context) {
	var payload schemas.LogisticServices

	// include access token

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
	id := models.CheckLogServiceById(logistic_id)
	if id != 1 {
		c.JSON(http.StatusNotFound, gin.H{"detail": "logistic_service not found."})
		return
	}

	// check duplicate name
	name := models.CheckLogServiceDuplicateName(payload.Name)
	if name != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "The name has already been taken."})
		return
	}

	// update logistic service
	models.UpdateLogisticService(logistic_id, payload.Name, payload.Margin)
	c.JSON(http.StatusOK, gin.H{"detail": "Successfully update the logistic_service."})
}

func DeleteLogisticService(c *gin.Context) {
	// include access token

	// get id
	param_id := c.Param("id")
	logistic_id, err_convert := strconv.Atoi(param_id)
	if err_convert != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "value is not a valid integer"})
		return
	}

	// check id is exist or not
	id := models.CheckLogServiceById(logistic_id)
	if id != 1 {
		c.JSON(http.StatusNotFound, gin.H{"detail": "logistic_service not found."})
		return
	}

	// delete logistic_service
	models.DeleteLogisticService(logistic_id)
	c.JSON(http.StatusOK, gin.H{"detail": "Successfully delete the guardian."})
}

func GetAllLogisticServices(c *gin.Context) {
	var (
		pagination schemas.PaginationLogistic
		payload    schemas.QueryPage
	)

	// validation
	if err := c.ShouldBind(&payload); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	// get param
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

	logistic_services, total := models.GetAllLogisticServices(q, pagination)

	// pages
	pages := pagesLogServices(pagination.PerPage, total)

	// prev_num
	has_prev := hasPrevLogServices(pagination.Page)
	prev_num := prevNumLogServices(has_prev, pagination.Page)

	// next_num
	has_next := hasNextLogServices(pagination.Page, pages)
	next_num := nextNumLogServices(has_next, pagination.Page)

	// iter_pages
	ledge := 2
	lcurrent := 2
	rcurrent := 5
	redge := 2
	iter := iterPagesLogServices(pagination.Page, pages, ledge, lcurrent, rcurrent, redge)

	c.JSON(http.StatusOK, gin.H{
		"data":       logistic_services,
		"total":      total,
		"prev_num":   prev_num,
		"next_num":   next_num,
		"page":       pagination.Page,
		"iter_pages": iter,
	})

}

func pagesLogServices(per_page int, total int) int {
	var page int

	if per_page == 0 || total == 0 {
		page = 0
	} else {
		page = int(math.Ceil(float64(total) / float64(per_page)))
	}

	return page
}

func prevNumLogServices(has_prev bool, page int) *int {
	if !has_prev {
		return nil
	}
	page = page - 1

	return &page
}

func hasPrevLogServices(page int) bool {
	// true if previous page exists
	return page > 1
}

func hasNextLogServices(page int, pages int) bool {
	// true if next page exists
	return page < pages
}

func nextNumLogServices(has_next bool, page int) *int {
	if !has_next {
		return nil
	}
	page = page + 1

	return &page
}

func iterPagesLogServices(page int, pages int, left_edge int, left_current int, right_current int, right_edge int) []*int {
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
