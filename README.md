# warehouse-api

## Examples

[![Golang](https://img.shields.io/badge/Go-v1.22-EEEEEE?logo=go&logoColor=white&labelColor=00ADD8)](https://go.dev/)

<div align="center">
    <h1>Warehouse API</h1>
    <h5>
        A simple design and implementation of JSON API with RPC-like operations as test task
    </h5>
</div>

---

---

## Navigation
* **[Task](#task)**
* **[Installation](#installation)**
* **[Example of requests](#examples-of-requests)**
* **[Additional features](#additional-features)**

---

## Task

* Design and implement API methods for working with `Products` in *many* `Warehouse`. Please note that an API call can be made simultaneously from different systems and they can work with the same products. API methods can be extended by parameters at your discretion. 

*Entity*
1. Warehouse
  * name
  * availability
2. Product
  * name
  * size
  * code
  * quantity

*API*
1. ReserveProducts:
```bash
{
  [
    {
      "code": "12345", // default parameter, possible to set several codes  
      "quantity": 8 // extended parameter by me for reserving any number of products
    }
  ]
}
```
2. ReleaseProducts:
```bash
{
  [
    {
      "reservation_id": "422ab5fa-fbf1-461a-99dc-2c6a49c323f1" // extended parameter by me for working with products that can be located in several warehouses
      "code": "12345", // default parameter, possible to set several codes  
      "quantity": 8 // extended parameter by me for releasing any number of products
    }
  ]
}
```
3. GetProductsByWarehouse:
```bash
{
  "warehouse_id": 3 // default parameter
}
```

| Requirement | Result |
| --- | --- |
| Use go fmt + goimports  | Done (for check: `make lint`) |
| Effective Go | Try to follow |
| JSON-API with RPC-like methods | Done |
| Use `PostgreSQL` or `MySQL`  | Done [PostrgeSQL](https://www.postgresql.org/) |
| `make up` to deploy the service  | Done |
| .http / curl in README.md / Postman collections | Done (curl in README.md + ./examples/*.http) |
| Testing code | Few Unit-tests in Repository's layer |
| Reasoning for choosing packages in go.mod | File `packages.md` |

## Installation
```shell
git clone https://github.com/pintoter/warehouse-api.git
```

---

## Getting started
1. **Create .env file with filename ".env" in the project root and setting up environment your own variables:**
```dotenv
# Database
DB_USER = "user"
DB_PASSWORD = "123456"
DB_HOST = "postgres"
DB_PORT = 5432
DB_NAME = "dbname"
DB_SSLMODE = "disable"

# Local database
LOCAL_DB_PORT = 5432
```
> **Hint:**
if you are running the project using Docker, set `DB_HOST` to "**postgres**" (as the service name of Postgres in the docker-compose).

2. **Compile and run the project:**
```shell
make up
```
3. **Service's structure**
```bash
├── Dockerfile
├── LICENSE
├── Makefile
├── README.md
├── cmd
│   └── warehouse-api
│       └── main.go
├── configs
│   └── main.yml
├── cover.out
├── docker-compose.yml
├── examples
│   ├── releaseproducts.http
│   ├── reserveproducts.http
│   └── showproducts.http
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   └── app.go
│   ├── config
│   │   └── config.go
│   ├── dbutil
│   │   ├── db.go
│   │   └── transaction
│   │       └── txmanager.go
│   ├── migrations
│   │   └── migrations.go
│   ├── repository
│   │   ├── model
│   │   │   └── products.go
│   │   ├── product
│   │   │   ├── repository.go
│   │   │   ├── reservation_create.go
│   │   │   ├── reservation_create_test.go
│   │   │   ├── reservation_get.go
│   │   │   ├── reservation_get_test.go
│   │   │   ├── reservation_update.go
│   │   │   ├── reservation_update_test.go
│   │   │   ├── warehouses_get.go
│   │   │   ├── warehouses_get_test.go
│   │   │   ├── warehouses_update.go
│   │   │   └── warehouses_update_test.go
│   │   └── repository.go
│   ├── server
│   │   └── server.go
│   ├── service
│   │   ├── model
│   │   │   ├── errors.go
│   │   │   ├── product.go
│   │   │   └── req.go
│   │   ├── product
│   │   │   └── service.go
│   │   └── service.go
│   └── transport
│       └── handler.go
├── migrations
│   ├── 20240228134525_init.down.sql
│   └── 20240228134525_init.up.sql
├── packages.md
└── pkg
    ├── database
    │   └── postgres
    │       └── postgres.go
    └── logger
        └── logger.go
```

---

## Examples of requests
### All submitted queries are executed `one by one` on the data that is lifted in the `migration`
#### 1. Reserve products
* Request example:
```shell
curl -X 'POST' \
  'http://localhost:8080/rpc' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "ProductService.ReserveProducts",
    "params": [{"products":[{"code": "12345", "quantity": 5}]}],
    "id": "coola"
    }'
```
* Response example:
```json
{
    "result": {
        "reservation_id": "51ef143e-4c4d-4b8b-a4dc-700ede2832e2",
        "reservation_products_info": [
            {
                "code": "12345",
                "status": "reserved"
            }
        ]
    },
    "error": null,
    "id": "coola"
}
```

* Request example:
```shell
curl -X 'POST' \
  'http://localhost:8080/rpc' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "ProductService.ReserveProducts",
    "params": [{"products":[{"code": "12345", "quantity": 5}]}],
    "id": "coola"
    }'
```
* Response example:
```json
{
    "result": {
        "reservation_id": "e216385f-74f0-4cb7-b3e4-532c11f6c428",
        "reservation_products_info": [
            {
                "code": "12345",
                "status": "rejected: required quantity of products is missing"
            }
        ]
    },
    "error": null,
    "id": "coola"
}
```

* Request example:
```shell
curl -X 'POST' \
  'http://localhost:8080/rpc' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
    -d '{
    "method": "ProductService.ReserveProducts",
    "params": [{"products":[{"code": "12345", "quantity": 5}, {"code": "12346", "quantity": 4}]}],
    "id": "coola"
  }'
```
* Response example:
```json
{
    "result": {
        "reservation_id": "ef99828c-b077-463a-82cc-4bc9abaffab1",
        "reservation_products_info": [
            {
                "code": "12345",
                "status": "rejected: required quantity of products is missing"
            },
            {
                "code": "12346",
                "status": "reserved"
            }
        ]
    },
    "error": null,
    "id": "coola"
}
```

#### 2. Release products
* Request example:
```shell
curl -X 'POST' \
  'http://localhost:8080/rpc' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "ProductService.ReleaseProducts",
    "params": [{"products":[{"reservation_id":"422ab5fa-fbf1-461a-99dc-2c6a49c323f1","code": "12345", "quantity": 2}]}],
    "id": "coola"
    }'
```
* Response example:
```json
{
    "result": {
        "release_products_info": [
            {
                "reservation_id": "422ab5fa-fbf1-461a-99dc-2c6a49c323f1",
                "code": "12345",
                "status": "released"
            }
        ]
    },
    "error": null,
    "id": "coola"
}
```

* Request example:
```shell
curl -X 'POST' \
  'http://localhost:8080/rpc' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "ProductService.ReleaseProducts",
    "params": [{"products":
    [
      {"reservation_id":"422ab5fa-fbf1-461a-99dc-2c6a49c323f1","code": "12345", "quantity": 2},
      {"reservation_id":"965ac486-0451-4e87-be55-2f985cdbf292","code": "12346", "quantity": 2}
    ]
    }],
    "id": "coola"
    }'
```
* Response example:
```json
{
    "result": {
        "release_products_info": [
            {
                "reservation_id": "965ac486-0451-4e87-be55-2f985cdbf292",
                "code": "12346",
                "status": "released"
            },
            {
                "reservation_id": "422ab5fa-fbf1-461a-99dc-2c6a49c323f1",
                "code": "12345",
                "status": "released"
            }
        ]
    },
    "error": null,
    "id": "coola"
}
```

#### 3. Get products by warehouse
* Request example:
```shell
  curl -X 'POST' \
  'http://localhost:8080/rpc' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "ProductService.GetProductsByWarehouse",
    "params": [{"warehouse_id": 1}],
    "id": "coola"
    }'
```
* Response example:
```json
{
    "result": [
        {
            "id": 1,
            "name": "Lacoste T-Shirt",
            "size": "XS",
            "code": "12345",
            "quantity": 2
        },
        {
            "id": 2,
            "name": "Lacoste T-Shirt",
            "size": "S",
            "code": "12346",
            "quantity": 3
        },
        {
            "id": 3,
            "name": "Lacoste T-Shirt",
            "size": "M",
            "code": "12347",
            "quantity": 1
        },
        {
            "id": 4,
            "name": "Lacoste T-Shirt",
            "size": "L",
            "code": "12348",
            "quantity": 2
        }
    ],
    "error": null,
    "id": "coola"
}
```

* Request example:
```shell
  curl -X 'POST' \
  'http://localhost:8080/rpc' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "ProductService.GetProductsByWarehouse",
    "params": [{"warehouse_id": 2}],
    "id": "coola"
    }'
```
* Response example:
```json
{
    "result": [
        {
            "id": 1,
            "name": "Lacoste T-Shirt",
            "size": "XS",
            "code": "12345",
            "quantity": 3
        },
        {
            "id": 5,
            "name": "Lacoste T-Shirt",
            "size": "XL",
            "code": "12349",
            "quantity": 3
        },
        {
            "id": 6,
            "name": "Dads pants",
            "size": "L",
            "code": "1337",
            "quantity": 5
        },
        {
            "id": 7,
            "name": "Dads pants",
            "size": "XL",
            "code": "1338",
            "quantity": 1
        },
        {
            "id": 8,
            "name": "Dads pants",
            "size": "XXL",
            "code": "1339",
            "quantity": 2
        }
    ],
    "error": null,
    "id": "coola"
}
```


* Request example:
```shell
  curl -X 'POST' \
  'http://localhost:8080/rpc' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
    "method": "ProductService.GetProductsByWarehouse",
    "params": [{"warehouse_id": 3}],
    "id": "coola"
    }'
```
* Response example:
```json
{
    "result": [
        {
            "id": 1,
            "name": "Lacoste T-Shirt",
            "size": "XS",
            "code": "12345",
            "quantity": 3
        },
        {
            "id": 9,
            "name": "Adidas Hoodie",
            "size": "L",
            "code": "10101011",
            "quantity": 3
        },
        {
            "id": 10,
            "name": "Adidas Hoodie",
            "size": "XL",
            "code": "10101012",
            "quantity": 5
        },
        {
            "id": 11,
            "name": "Adidas Hoodie",
            "size": "XXL",
            "code": "10101013",
            "quantity": 1
        },
        {
            "id": 12,
            "name": "Nike Longsleeve",
            "size": "S",
            "code": "1111131231",
            "quantity": 2
        }
    ],
    "error": null,
    "id": "coola"
}
```
---

## Additional features
1. **Run tests**
```shell
make test
```
2. **Stop all running containers**
```shell
make stop
```
3. **Run linter**
```shell
make lint
```


