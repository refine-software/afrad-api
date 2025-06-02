# TODO

## Project setup

- [x] Implement schema as migrations.
- [x] Setup postgres DB using "pgx.pool" instead of "database/sql" package.
- [x] Setup viper for environment variables.
- [x] Setup S3 for image storing.
  - [ ] Create a bucket using aws and add the required info.
- [x] Write the project endpoints docs.
- [ ] Database and models structure, we need to decide what structure will we use, then implement it.
  - [ ] Make sure to expose custom database errors, to make error handling easier in the handlers.
- [ ] Setup an error package for custom and predefined API errors.
- [ ] We'll have 2 layers a database layer and an API layer, those layers will be used to separate the programmer concerns.
  - [ ] We'll expose functions through the Service interface
- [ ] Write all expected routes in the routes.go file

## Endpoints
