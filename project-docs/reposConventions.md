# Repository Convention

## Methods

On every repo we'll have the following methods:

1. Get
2. Create
3. Update
4. Delete
5. Any number of needed methods

These methods will take the model of the repo as a parameter and make the changes

## Comments

1. Create a comment for every method inside the repo interface (to make the comment visible for the user).
2. Make sure to tell the user 2 things, what will be changed and by what, for example:

```go
	// this method will update the following columns:
	// first_name, last_name, image, role.
	// based on the user_id.
```

## Method naming

- Don't use the repo name again in the the method name, for example:

```go
// Wrong naming, we're already using the userRepo, so ofc we'll get the user no shit
func (r *userRepo) GetUserRole() {}
// Correct naming
func (r *userRepo) GetRole() {}
```
