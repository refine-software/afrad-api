# Endpoints

| DONE | Method | Endpoint | Description |
| ---- | ------ | -------- | ----------- |
| ❌   | `POST` | `/`      |             |

## Authentication

### JWT

| DONE | Method | Endpoint                       | Description                                 |
| ---- | ------ | ------------------------------ | ------------------------------------------- |
| ✅   | `POST` | `/auth/register`               | Register a new user                         |
| ✅   | `POST` | `/auth/verify-account`         | Verify phone number to activate account     |
| ✅   | `POST` | `/auth/resend-verification`    | Resend verification otp to activate account |
| ✅   | `POST` | `/auth/login`                  | Login and receive JWT                       |
| ✅   | `POST` | `/auth/reset-password`         | Request a password reset                    |
| ✅   | `POST` | `/auth/reset-password/confirm` | Set a new password                          |
| ✅   | `POST` | `/auth/refresh`                | Refresh the access and refresh tokens       |

### Oauth

| DONE | Method | Endpoint                 | Description                                                                                                                  |
| ---- | ------ | ------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| ✅   | `GET`  | `/oauth/google/login`    | Redirects the user to Google's OAuth 2.0 authorization URL with PKCE support                                                 |
| ✅   | `GET`  | `/oauth/google/callback` | Handles Google's redirect with the auth code, exchanges it for tokens, verifies ID token, and signs in or registers the user |

## User

| DONE | Method   | Endpoint                         | Description                         |
| ---- | -------- | -------------------------------- | ----------------------------------- |
| ✅   | `GET`    | `/user`                          | Get user data                       |
| ✅   | `PUT`    | `/user`                          | Update user data                    |
| ✅   | `DELETE` | `/user`                          | Delete self                         |
| ✅   | `GET`    | `/user/reviews`                  | Fetch all reviews of a user         |
| ✅   | `POST`   | `/user/reviews`                  | Review a product                    |
| ✅   | `PUT`    | `/user/reviews/:id`              | Update review                       |
| ✅   | `DELETE` | `/user/reviews/:id`              | Delete review                       |
| ❌   | `PATCH`  | `/user/notificatoin-preferences` | Change notifications preference     |
| ❌   | `GET`    | `/user/notificatoin-preferences` | Fetch user notifications preference |
| ✅   | `POST`   | `/user/logout`                   | Revoke the current session          |
| ✅   | `POST`   | `/user/logout/all`               | Revoke all sessions                 |

## Product

| DONE | Method   | Endpoint             | Description                                          |
| ---- | -------- | -------------------- | ---------------------------------------------------- |
| ✅   | `GET`    | `/products`          | Fetch products with pagination and search and filter |
| ✅   | `GET`    | `/product/:id`       | Fetch product details                                |
| ❌   | `PUT`    | `/admin/product/:id` | Update product (Admin only)                          |
| ❌   | `DELETE` | `/admin/product/:id` | Delete product (Admin only)                          |
| ❌   | `POST`   | `/admin/product`     | Add a product (Admin only)                           |

## Category

| DONE | Method   | Endpoint               | Description                     |
| ---- | -------- | ---------------------- | ------------------------------- |
| ✅   | `GET`    | `/products/categories` | Fetch all categories            |
| ✅   | `POST`   | `/admin/category`      | Add a new category (admin only) |
| ✅   | `PATCH`  | `/admin/category/:id`  | Update category (admin only)    |
| ✅   | `DELETE` | `/admin/category/:id`  | Delete category (admin only)    |

## Cart

| DONE | Method   | Endpoint     | Description                              |
| ---- | -------- | ------------ | ---------------------------------------- |
| ❌   | `GET`    | `/cart`      | Fetch cart details (cart and cart items) |
| ❌   | `POST`   | `/cart/item` | Add item to cart                         |
| ❌   | `PATCH`  | `/cart/:id`  | Update cart item quantity                |
| ❌   | `DELETE` | `/cart/:id`  | Delete cart item quantity                |

## Wishlist

| DONE | Method   | Endpoint        | Description               |
| ---- | -------- | --------------- | ------------------------- |
| ❌   | `GET`    | `/wishlist`     | Fetch the wishlist        |
| ❌   | `POST`   | `/wishlist/:id` | Add item to wishlist      |
| ❌   | `DELETE` | `/wishlist/:id` | Delete item from wishlist |

## Order

| DONE | Method  | Endpoint            | Description                   |
| ---- | ------- | ------------------- | ----------------------------- |
| ❌   | `GET`   | `/admin/orders`     | Fetch all orders (Admin only) |
| ❌   | `GET`   | `/orders`           | Fetch all orders              |
| ❌   | `GET`   | `/order/:id`        | Fetch a specific order        |
| ❌   | `POST`  | `/order`            | Add order (checkout)          |
| ❌   | `PATCH` | `/order/:id/cancel` | Cancel a specific order       |

## Discount

| DONE | Method   | Endpoint                   | Description                                                    |
| ---- | -------- | -------------------------- | -------------------------------------------------------------- |
| ❌   | `POST`   | `/admin/discount/product`  | Add a discount for a product (admin only)                      |
| ❌   | `POST`   | `/admin/discount/variants` | Add a discount for a variant or more of a product (admin only) |
| ❌   | `GET`    | `/admin/discounts`         | Fetch all discounts (admin only)                               |
| ❌   | `PUT`    | `/admin/discount/:id`      | Update discount (admin only)                                   |
| ❌   | `DELETE` | `/admin/discount/:id`      | Delete a discount (admin only)                                 |

## Notification

| DONE | Method | Endpoint | Description |
| ---- | ------ | -------- | ----------- |
| ❌   | `POST` | `/`      |             |

## Messaging

| DONE | Method | Endpoint | Description |
| ---- | ------ | -------- | ----------- |
| ❌   | `POST` | `/`      |             |

## Colors

| DONE | Method   | Endpoint     | Description       |
| ---- | -------- | ------------ | ----------------- |
| ✅   | `GET`    | `/colors`    | Get all colors    |
| ✅   | `POST`   | `/color`     | Create color      |
| ✅   | `PUT`    | `/color/:id` | Update color name |
| ✅   | `DELETE` | `/color/:id` | Delete color      |

## Sizes

| DONE | Method   | Endpoint    | Description      |
| ---- | -------- | ----------- | ---------------- |
| ❌   | `GET`    | `/sizes`    | Get all sizes    |
| ❌   | `POST`   | `/size`     | Create size      |
| ❌   | `PUT`    | `/size/:id` | Update size name |
| ❌   | `DELETE` | `/size/:id` | Delete size      |

---

# INFO

1. Mark implemented endpoints with ✅.
2. Mark not yet implemented endpoints with ❌.
