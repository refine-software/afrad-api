package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/refine-software/afrad-api/internal/database"
	"github.com/refine-software/afrad-api/internal/models"
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

	db := s.DB.Pool()
	product := s.DB.Product()
	products, metadata, err := product.GetAll(c, db, f, &productsFilterOptions)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
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

type productDetailsRes struct {
	Product           database.ProductDetails            `json:"product"`
	ProductVariants   []database.ProductVariantDetails   `json:"productVariants"` // will contain the color and size of each variant
	RatingsAndReviews []database.RatingsAndReviewDetails `json:"ratingsAndReviews"`
	Images            []models.Image                     `json:"images"`
	Discount          []models.Discount                  `json:"discount"`
}

func (s *Server) getProduct(c *gin.Context) {
	productID := convStrToInt(c, c.Param("id"), "product id")
	if productID == 0 {
		return
	}

	db := s.DB.Pool()
	productRepo := s.DB.Product()
	productVariantRepo := s.DB.ProductVariant()
	ratingsAndReviewsRepo := s.DB.RatingReview()
	imageRepo := s.DB.Image()

	p, err := productRepo.GetDetails(c, db, productID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	pvs, err := productVariantRepo.GetAllOfProduct(c, db, p.ID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	rrs, err := ratingsAndReviewsRepo.GetAllOfProduct(c, db, p.ID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	imgs, err := imageRepo.GetAllOfProduct(c, db, p.ID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, productDetailsRes{
		Product:           *p,
		ProductVariants:   pvs,
		RatingsAndReviews: rrs,
		Images:            imgs,
	})
}

type productVariant struct {
	Quantity int   `json:"quantity" binding:"required"`
	Price    int   `json:"price"    binding:"required"`
	ColorID  int32 `json:"colorId"  binding:"required"`
	SizeID   int32 `json:"sizeId"   binding:"required"`
}

type addProductReq struct {
	Name       string           `json:"name"       binding:"required"`
	Details    string           `json:"details"    binding:"required"`
	BrandID    int32            `json:"brandId"    binding:"required"`
	CategoryID int32            `json:"categoryId" binding:"required"`
	Variants   []productVariant `json:"variants"   binding:"required"`
}

func (s *Server) addProduct(c *gin.Context) {
	// get product and variants from the form as json string
	productJSON := c.PostForm("product")

	var req addProductReq
	err := json.Unmarshal([]byte(productJSON), &req)
	if err != nil {
		utils.Fail(c, utils.ErrBadRequest, err)
		return
	}

	// get thumbnail image
	imageUpload, apiErr := getImageFile(c, "thumbnail", 1000<<10) // 1000 << 10 this equals 1MB
	if apiErr != nil {
		utils.Fail(c, apiErr, err)
		return
	}

	if imageUpload == nil {
		utils.Fail(c, &utils.APIError{
			Code:    http.StatusBadRequest,
			Message: "thumbnail image is required",
		}, err)
		return
	}
	thumbnailURL, err := s.S3.UploadImage(c, imageUpload.File, imageUpload.Header)
	if err != nil {
		utils.Fail(
			c,
			&utils.APIError{
				Code:    http.StatusInternalServerError,
				Message: "couldn't upload product thumbnail image",
			},
			err,
		)
		return
	}
	imageUpload.File.Close()

	productRepo := s.DB.Product()
	variantRepo := s.DB.ProductVariant()
	imageRepo := s.DB.Image()
	db, err := s.DB.BeginTx(c)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	committed := false
	defer func() {
		if p := recover(); p != nil {
			_ = db.Rollback(c)
			panic(p)
		} else if !committed {
			_ = db.Rollback(c)
			_ = s.S3.DeleteImageByURL(c, thumbnailURL)
		}
	}()

	p := models.Product{
		Name:            req.Name,
		Details:         pgtype.Text{String: req.Details, Valid: req.Details != ""},
		BrandID:         req.BrandID,
		ProductCategory: req.CategoryID,
		Thumbnail:       thumbnailURL,
	}
	productID, err := productRepo.Create(c, db, &p)
	if err != nil {
		apiErr = utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	// create product variants
	var pv models.ProductVariant
	for _, v := range req.Variants {
		pv = models.ProductVariant{
			Quantity:  v.Quantity,
			Price:     v.Price,
			ColorID:   v.ColorID,
			SizeID:    v.SizeID,
			ProductID: productID,
		}
		err = variantRepo.Create(c, db, &pv)
		if err != nil {
			apiErr = utils.MapDBErrorToAPIError(err)
			utils.Fail(c, apiErr, err)
			return
		}
	}

	// get the product images, required
	uploads, apiErr := getImageFiles(c, "images", 4<<20)
	if apiErr != nil {
		utils.Fail(c, apiErr, errors.New("failed to get image uploads"))
		return
	}

	if len(uploads) == 0 || uploads == nil {
		utils.Fail(c, &utils.APIError{
			Code:    http.StatusBadRequest,
			Message: "product images are required",
		}, errors.New("no images found"))
		return
	}

	defer func() {
		for _, u := range uploads {
			u.File.Close()
		}
	}()

	var (
		imagePair  *ImagePair
		imagePairs []ImagePair
	)
	for _, upload := range uploads {
		imagePair, apiErr = s.UploadImageWithLowRes(c, upload.File, upload.Header, 500)
		if apiErr != nil {
			utils.Fail(c, apiErr, errors.New("couldn't upload image"))
			return
		}

		imagePairs = append(imagePairs, *imagePair)
	}

	var image models.Image
	for _, imgPair := range imagePairs {
		image = models.Image{
			Image:       imgPair.OriginalURL,
			LowResImage: imgPair.LowResURL,
			ProductID:   productID,
		}
		err = imageRepo.Create(c, db, &image)
		if err != nil {
			apiErr = utils.MapDBErrorToAPIError(err)
			utils.Fail(c, apiErr, err)
			return
		}
	}

	err = db.Commit(c)
	if err != nil {
		utils.Fail(c, utils.ErrInternal, err)
		return
	}

	committed = true
	utils.Created(c, "product created successfully")
}

type productReq struct {
	Name       string `json:"name"`
	Details    string `json:"details"`
	BrandID    int32  `json:"brandId"`
	CategoryID int32  `json:"categoryId"`
}

func (s *Server) updateProduct(c *gin.Context) {
	productID := convStrToInt(c, c.Param("id"), "product id")
	if productID == 0 {
		return
	}

	var req productReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		utils.Fail(c, utils.ErrBadRequest, err)
		return
	}

	db := s.DB.Pool()
	productRepo := s.DB.Product()

	p, err := productRepo.Get(c, db, productID)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	if strings.TrimSpace(req.Name) != "" {
		p.Name = req.Name
	}

	if strings.TrimSpace(req.Details) != "" {
		p.Details = pgtype.Text{String: req.Details, Valid: strings.TrimSpace(req.Details) != ""}
	}

	if req.BrandID != 0 {
		p.BrandID = req.BrandID
	}

	if req.CategoryID != 0 {
		p.ProductCategory = req.CategoryID
	}

	err = productRepo.Update(c, db, p)
	if err != nil {
		apiErr := utils.MapDBErrorToAPIError(err)
		utils.Fail(c, apiErr, err)
		return
	}

	utils.Success(c, "product updated successfully")
}
