# No SQL DB using Go

## Prerequisite

> go 1.15

## Run test

> go test -v ./...

## Env Var

| Variable | Description                      | Default Value | Possible Values             |
| -------- | -------------------------------- | ------------- | --------------------------- |
| DB_DATA  | dir location to store json files | /tmp/         | /User/yourUserName/Desktop/ |

> export DB_DATA="/tmp/db/"

## Build (Generate executable) and Run

> go build

## Get help

> ./go-nosql-db -h

```
Usage: go-nosql-db --table TABLE [--filter-key FILTER-KEY] [--filter-value FILTER-VALUE] [--select SELECT] [--data DATA] [--delete]

Options:
  --table TABLE          --table db table name
  --filter-key FILTER-KEY, -k FILTER-KEY
                         optional: filter condition
  --filter-value FILTER-VALUE, -v FILTER-VALUE
                         optional: filter condition
  --select SELECT, -s SELECT
                         optional: json nodes to select
  --data DATA, -d DATA   optional json body
  --delete               delete record for filter
  --help, -h             display this help and exit
```

## Insert a record

```bash
./go-nosql-db --table="yourtable" --data='{
  "name": "neeraj",
  "lastName": "dubey",
  "address": {
  	"line1": "24 g bank road karta",
  	"state": "uttar-pradesh"
  },
  "dob": "15-05-1986"
}'
```

## Fetch a record

```bash
./go-nosql-db --table="yourtable" --filter-key="name" --filter-value="neeraj" --select="id,name,address"
```

### _$$ Note: Fetch a record By ID - use ony --filter-value_

```bash
./go-nosql-db --table="yourtable" --filter-value="a5bdf0df-432e-4272-69ea-87b13ba06a8a"
```

## _Fetch a nested deep Record & select nested fields_

```bash
./go-nosql-db --table="yourtable" \
--filter-key="address.state" \
--filter-value="uttar-pradesh" \
--select="id,name,address.line1"
```

## Delete record for filter condition

```bash
./go-nosql-db --table="yourtable" \
--filter-key="address.state" \
--filter-value="uttar-pradesh" \
--delete
```
