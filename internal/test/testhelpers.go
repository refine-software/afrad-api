package test

import (
	"context"

	"github.com/refine-software/afrad-api/internal/database"
)

func truncateAllTables(db database.Service) error {
	ctx := context.Background()
	_, err := db.Pool().Exec(ctx, `
        TRUNCATE TABLE
            order_details,
            orders,
            wishlists,
            cart_items,
            carts,
            images,
            rating_review,
            product_variants,
            products,
            brands,
            categories,
            colors,
            sizes,
            password_resets,
            account_verification_codes,
            sessions,
            oauth,
            local_auth,
            users,
            cities,
            discounts,
            variant_discount
        RESTART IDENTITY CASCADE;
    `)
	return err
}
