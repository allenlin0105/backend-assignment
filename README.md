# Dcard Backend Assignment
This is the assignment for the Dcard 2024 backend internship.

## Setup
0. I run the server at my local with MySQL. Therefore, MySQL should be installed in advance. After installation, create a database name `test` and run the two sql files, where `sql/setup.sql` should be run first and `sql/insert.sql` the second.
1. Create `.env` file with the following key-value pair
```
APP_PORT=3000
CONTEXT_TIMEOUT=2
MYSQL_USERNAME=root
MYSQL_PASSWORD=$YOUR_PASSWORD
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_DATABASE=test
```
2. Run `go run ./main.go`
3. Test the API at host `127.0.0.1:3000`

## API Spec
After running `go run ./main.go`, refer to `http://127.0.0.1:3000/swagger/index.html`

## Design
I create 4 tables for `ads`, `genders`, `countries`, and `platforms`, respectively. The schema for each table is shown below.
```
ads
+-----------+--------------+------+-----+---------+----------------+
| Field     | Type         | Null | Key | Default | Extra          |
+-----------+--------------+------+-----+---------+----------------+
| id        | int unsigned | NO   | PRI | NULL    | auto_increment |
| title     | varchar(128) | NO   |     | NULL    |                |
| start_at  | timestamp    | NO   |     | NULL    |                |
| end_at    | timestamp    | NO   |     | NULL    |                |
| age_start | int unsigned | NO   |     | NULL    |                |
| age_end   | int unsigned | NO   |     | NULL    |                |
+-----------+--------------+------+-----+---------+----------------+

genders
+--------+--------------+------+-----+---------+----------------+
| Field  | Type         | Null | Key | Default | Extra          |
+--------+--------------+------+-----+---------+----------------+
| id     | int unsigned | NO   | PRI | NULL    | auto_increment |
| gender | varchar(2)   | NO   | UNI | NULL    |                |
+--------+--------------+------+-----+---------+----------------+

countries
+---------+-------------+------+-----+---------+-------+
| Field   | Type        | Null | Key | Default | Extra |
+---------+-------------+------+-----+---------+-------+
| id      | int         | NO   | PRI | NULL    |       |
| country | varchar(2)  | NO   | UNI | NULL    |       |
| alpha3  | varchar(3)  | NO   | UNI | NULL    |       |
| langCS  | varchar(45) | NO   |     | NULL    |       |
| langDE  | varchar(45) | NO   |     | NULL    |       |
| langEN  | varchar(45) | NO   |     | NULL    |       |
| langES  | varchar(45) | NO   |     | NULL    |       |
| langFR  | varchar(45) | NO   |     | NULL    |       |
| langIT  | varchar(45) | NO   |     | NULL    |       |
| langNL  | varchar(45) | NO   |     | NULL    |       |
+---------+-------------+------+-----+---------+-------+

platforms
+----------+--------------+------+-----+---------+----------------+
| Field    | Type         | Null | Key | Default | Extra          |
+----------+--------------+------+-----+---------+----------------+
| id       | int unsigned | NO   | PRI | NULL    | auto_increment |
| platform | varchar(8)   | NO   | UNI | NULL    |                |
+----------+--------------+------+-----+---------+----------------+
```
Since linking `ads` to `genders`, `countries`, and `platforms` is a M-to-N relationship, I use another 3 tables for linking. The schema is shown below.
```
ad_gender
+-----------+--------------+------+-----+---------+-------+
| Field     | Type         | Null | Key | Default | Extra |
+-----------+--------------+------+-----+---------+-------+
| ad_id     | int unsigned | NO   | MUL | NULL    |       |
| gender_id | int unsigned | NO   | MUL | NULL    |       |
+-----------+--------------+------+-----+---------+-------+

ad_country
+------------+--------------+------+-----+---------+-------+
| Field      | Type         | Null | Key | Default | Extra |
+------------+--------------+------+-----+---------+-------+
| ad_id      | int unsigned | NO   | MUL | NULL    |       |
| country_id | int          | NO   | MUL | NULL    |       |
+------------+--------------+------+-----+---------+-------+

ad_platform
+-------------+--------------+------+-----+---------+-------+
| Field       | Type         | Null | Key | Default | Extra |
+-------------+--------------+------+-----+---------+-------+
| ad_id       | int unsigned | NO   | MUL | NULL    |       |
| platform_id | int unsigned | NO   | MUL | NULL    |       |
+-------------+--------------+------+-----+---------+-------+
```

### Create an ad
When creating a new ad, I append rows to `ads` and the 3 linking tables. 

Age, gender, country, and platform are optional, so I assign "any" value, which corresponds to no restiction. For example, ageStart is set to 1 and ageEnd is set to 100. For gender, country, and platform, the "any" value is "A", "AY", and "any", respectively.

### Get ads
When getting ads, if the condition is not provided, then I don't need to check the corresponding field or table. For example, if the gender condition is not provided, then I pass checking the linking table `ad_gender`. 

For condition gender, country and platform, if they are provided, then I'll add "any" value into query in order to get the ads that do not have restriction on these fields.