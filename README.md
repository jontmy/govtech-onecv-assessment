# govtech-onecv-assessment
Technical assessment for GovTech GDS OneCV internship program.

### Developer Guide
This section will guide you through the steps to get the server up and running.

#### Pre-requisites
You should have the following installed on your machine before you proceed:
- [*Go*](https://golang.org/doc/install) 1.20 or above
- [*MySQL*](https://dev.mysql.com/doc/refman/8.0/en/installing.html) 8.0 or above

#### Setup

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
5. Build and run the server file.
```
go run .
```