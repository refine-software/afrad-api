# TODO

## Project setup

- [x] Implement schema as migrations.
- [x] Setup postgres DB using "pgx.pool" instead of "database/sql" package.
- [x] Setup viper for environment variables.
- [x] Setup S3 for image storing.
  - [x] Create a bucket using aws and add the required info.
- [x] Write the project endpoints docs.
- [x] Database and models structure, we need to decide what structure will we use, then implement it.
  - [x] Make sure to expose custom database errors, to make error handling easier in the handlers.
- [x] Error handling
  - [x] Setup an error package for custom and predefined API errors.
  - [x] Setup A predefined and exposed errors in the database package.
- [x] Write all expected routes in the routes.go file

## API Features

- [ ] We need pagination on a lot of endpoints, try to make a reusable function or a struct to make implementing pagination more simple.

## Database

- [x] Write all tables as structs (models).

## Bugs/Issues

- [ ] session management is broken, the user is not reusing a session and we're not removing them, so the number of session will keep on growing.
- [ ]
