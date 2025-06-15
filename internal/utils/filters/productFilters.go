package filters

import (
	"fmt"
	"strings"
)

type ProductFilterOptions struct {
	CategoryID int
	BrandID    int
	Search     string
}

func (p *ProductFilterOptions) GetWhereClause() (string, []any) {
	var whereClauses []string
	var args []any
	argIndex := 3

	// Example filters
	if p.CategoryID != 0 {
		whereClauses = append(
			whereClauses,
			fmt.Sprintf("products.product_category = $%d", argIndex),
		)
		args = append(args, p.CategoryID)
		argIndex++
	}

	if p.BrandID != 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("products.brand_id = $%d", argIndex))
		args = append(args, p.BrandID)
		argIndex++
	}

	if p.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("products.name ILIKE $%d", argIndex))
		args = append(args, "%"+p.Search+"%")
		argIndex++
	}

	// Join WHERE clause if needed
	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	fmt.Println(whereSQL)

	return whereSQL, args
}
