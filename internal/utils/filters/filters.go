package filters

import (
	"math"
	"strings"

	"github.com/refine-software/afrad-api/internal/utils/validator"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafeList []string
}

// Metadata holds pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	FirstPage    int `json:"firstPage,omitempty"`
	LastPage     int `json:"lastPage,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
}

// calculateMetadata calculates the appropriate pagination metadata values given the total number
// of records, current page, and page size values. Note, the last page value is calculated using the
// math.Ceil() function, which rounds up a float to the nearest integer. So, for example, if there
// were 13 records in total and a page size of 5, the last page value would be math.Ceil(13/5) = 3.
func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{} // return an empty Metadata struct if there are no records
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

// ValidateFilters runs validation checks on the Filters type.
func ValidateFilters(v *validator.Validator, f Filters) {
	// Check that page and page_size parameters contain sensible values.
	v.Check(f.Page > 0, "page", "must be greater than 0")
	v.Check(f.Page <= 10_000_0000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than 0")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.In(f.Sort, f.SortSafeList...), "sort", "invalid sort value")
}

// sortColumn checks that the client-provided Sort field matches one of the entries in our
// SortSafeList and if it does, it extracts the column name from the Sort field by stripping the
// leading hyphen character (if one exists).
func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafeList {
		if f.Sort == safeValue {
			column := strings.TrimPrefix(f.Sort, "-")

			switch column {
			case "price":
				return "MIN(product_variants.price)"
			case "rating":
				return "COALESCE(ROUND(AVG(DISTINCT rating_review.rating)::numeric, 2), 0.00)"
			default:
				return "products." + column // assume default columns are in `products`
			}
		}
	}

	// The panic below should technically not happen because the Sort value should have already
	// been checked when calling the ValidateFilters helper function. However, this is a sensible
	// failsafe to help stop a SQL injection attach from occurring.
	panic("unsafe sort parameter:" + f.Sort)
}

// sortDirection returns the sort direction ("ASC" or "DESC") depending on the prefix character
// of the Sort field.
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}
