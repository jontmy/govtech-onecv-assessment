# govtech-onecv-assessment
Technical assessment for GovTech GDS OneCV internship program.

## Developer Guide
This section will guide you through the steps to get the server up and running.

### Pre-requisites
You should have the following installed on your machine before you proceed:
- [*Go*](https://golang.org/doc/install) 1.20 or above
- [*MySQL*](https://dev.mysql.com/doc/refman/8.0/en/installing.html) 8.0 or above

### Setup

1. Clone the repository.
```
git clone https://github.com/jontmy/govtech-onecv-assessment.git
```
2. Navigate to the project directory.
```
cd govtech-onecv-assessment
```
3. Create a `.env` file in the root directory of the project from the example file, and populate it.
```
cp .env.example .env
vim .env  # or any other text editor of your choice
```
4. If you do not already have a MySQL database, create one.
```
mysql -u <username> -p
```
You will need to enter your password when prompted.
Now, at the `mysql` command prompt, create the database and run the DDL file.
```
mysql> CREATE DATABASE <database_name>;
mysql> USE <database_name>;
mysql> SOURCE src/database/db_ddl.sql;
```
Once that's done, exit the `mysql` command prompt.
```
mysql> exit
```
5. Build and run the server file. You can then make requests to the server on localhost at the port specified in the `.env` file.
```
go run .
```

## API Reference

For all API endpoints, the following constraints apply:

- If the wrong HTTP method is used, the request will **fail** with status code 405.
- If the content type is not set to `application/json` where required, the request will **fail** with status code 415.

Although not required, the `DELETE /api/reset` endpoint is provided for convenience.

### `POST /api/register`

Registers one or more students to a specified teacher.

- If any of the students or teacher are not already in the database, they will be added.
- A student can be registered to multiple teachers.
- A teacher can have multiple students registered to them.
- The header `Content-Type` must be set to `application/json`.

**Fails** with status code 422 if no teacher *and* no students are specified.

**Fails** with status code 400 if the JSON body is otherwise invalid or missing.

**Succeeds** with status code 204 and an empty response body if the teacher and students are successfully registered.

### `GET /api/commonstudents`

Retrieves the list of students common to a given list of teachers.

- The teachers should be specified via the query parameter `teacher`.
- If no teachers are specified, this will return all students in the database.
- If no students are found, this will return an empty list.

**Succeeds** with status code 200 and a JSON response body containing the list of common students if the request is successful.

### `POST /api/suspend`

Suspends a specified student.

- The header `Content-Type` must be set to `application/json`.
- The student to be suspended should be specified via the JSON body.
- The student must already be registered in the database.
- If the student is already suspended, this request will have no effect.

**Fails** with status code 422 if no student is specified.

**Fails** with status code 404 if the student is not found in the database.

**Fails** with status code 400 if the JSON body is otherwise invalid or missing.

**Succeeds** with status code 204 and an empty response body if the student is successfully suspended.

### `POST /api/retrievefornotifications`

Retrieves the list of students who can receive a given notification.

- The header `Content-Type` must be set to `application/json`.
- The teacher and notification should be specified via the JSON body.
- The notification can contain zero or more @-mentioned students after the message.
- If no teacher is specified, only the @-mentioned students will be returned.

**Fails** with status code 404 if the teacher is specified (i.e. non-empty string) but is not found in the database.

**Fails** with status code 400 if the JSON body is otherwise invalid or missing.

**Succeeds** with status code 200 and a JSON response body containing the list of students who can receive the notification if the request is successful.

### `DELETE /api/reset`

Resets the database, un-registering all teachers and students.

**Succeeds** with status code 204 and an empty response body if the database is successfully reset.