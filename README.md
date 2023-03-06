# govtech-onecv-assessment
Technical assessment for GovTech GDS OneCV internship program.

### Developer Guide
This section will guide you through the steps to get the server up and running.

#### Pre-requisites
Before continuing any further, you should have the following installed on your machine:
- [*Go*](https://golang.org/doc/install) 1.20 or above
- [*MySQL*](https://dev.mysql.com/doc/refman/8.0/en/installing.html) 8.0 or above

#### Setup

1. Clone the respository.
```bash
git clone https://github.com/jontmy/govtech-onecv-assessment.git
```
2. Navigate to the project directory.
```bash
cd govtech-onecv-assessment
```
3. Create a `.env` file in the root directory of the project from the example file, and populate it.
```bash
cp .env.example .env
vim .env  # or any other text editor of your choice
```
4. If you do not already have a MySQL database, create one.
```bash
mysql -u <username> -p
```
You will need to enter your password when prompted.
Now, at the `mysql` command prompt, create the database and run the DDL file.
```sql
mysql> CREATE DATABASE <database_name>;
mysql> USE <database_name>;
mysql> SOURCE mysql/db_ddl.sql;
```
Once that's done, exit the `mysql` command prompt.
```sql
mysql> exit
```
5. Build and run the server file.
```bash
go run .
```