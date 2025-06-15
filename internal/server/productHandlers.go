package server

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/utils"
	"github.com/refine-software/afrad-api/internal/utils/filters"
	"github.com/refine-software/afrad-api/internal/utils/validator"
)

type productsRes struct {
	Metadata filters.Metadata   `json:"metadata"`
	Products []database.Product `json:"products"`
}

func (s *Server) getAllProducts(c *gin.Context) {
	// Initialize a new Validator instance.
	v := validator.New()

	page := getRequiredQueryInt(c, "page")
	if page == 0 {
		return
	}

	pageSize := getRequiredQueryInt(c, "page_size")
	if pageSize == 0 {
		return
	}

	sort := c.Query("sort")
	if sort == "" {
		sort = "id"
	}

	var f filters.Filters
	f.Page = int(page)
	f.PageSize = int(pageSize)
	f.Sort = sort

	// Add the supported sort value for this endpoint to the sort safelist.
	f.SortSafeList = []string{
		// ascending sort values
		"id", "price", "name", "rating",
		// descending sort values
		"-id", "-price", "-name", "-rating",
	}

	// Execute the validation checks on the Filters struct and send a response
	// containing the errors if necessary.
	if filters.ValidateFilters(v, f); !v.Valid() {
		utils.Fail(c, utils.ErrBadRequest, errors.New("bad filter options"))
		return
	}

	categoryIDStr := c.Query("category_id")
	categoryID, _ := strconv.Atoi(categoryIDStr)

	brandIDStr := c.Query("brand_id")
	brandID, _ := strconv.Atoi(brandIDStr)

	search := c.Query("search")

	productsFilterOptions := filters.ProductFilterOptions{
		CategoryID: categoryID,
		BrandID:    brandID,
		Search:     search,
	}

	fmt.Println(productsFilterOptions)

	db := s.db.Pool()
	product := s.db.Product()
	products, metadata, dbErr := product.GetAll(c, db, f, &productsFilterOptions)
	if dbErr != nil {
		apiErr := utils.MapDBErrorToAPIError(dbErr, "product")
		utils.Fail(c, apiErr, dbErr)
		return
	}

	if len(products) == 0 || products == nil {
		utils.NoContent(c)
		return
	}

	utils.Success(c, productsRes{
		Metadata: metadata,
		Products: products,
	})
}
