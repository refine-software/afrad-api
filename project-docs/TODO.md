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
- [ ] We need a better logger, something like zap.

## Database

- [x] Write all tables as structs (models).

## Bugs/Issues

- [x] Session management is broken, the user is not reusing a session and we're not removing them, so the number of session will keep on growing.
- [x] Implement some kind of a tracer, when a database level error happens there is no way to know where this error is coming from.
- [x] In registeration, images are being uploaded to s3 while the database transactions might fail.
- [x] Make sure to rollback like in the resend verification endpoint.
- [x] Scan files sent, they have to be images and have to be under a certain size.
- [x] S3 images aren't accessible in the browser, check bucket permissions.
- [x] DBErrors needs to containe the database error itself too, so we can show it to the programmer.
- [x] The database errors when being mapped are getting a column name, which is a bad why of determining which column is causing the issue, find a better way.
      To reproduce the bad behavior, try to submit a review then submit the review again on the same product, you'll get a conflict,
      Then try to post a review on a non existing product and see the returned message, compare the messages from these 2 cases and you'll see the problem.
- [ ] use result of db.Exec in every update or delete method instead of checking existence.
- [ ] use convStrToInt helper function instead of getting the param manually.
