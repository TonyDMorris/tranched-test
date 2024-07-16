## Tranched Test

This repository is structured in a typical go monolith style with extensibility to migrate into microservices

### `./cmd`

contains application entrypoints

### `./internal`

contains application specific repositories and domain types

### `./pkg`

contains various packages required by the application

### `./build`

contains build speciifc configuration and utility tools

### commands

run the services with docker-compose from the root of the project

```bash
docker-compose up --build
```

check bobs assets

```curl
curl --location 'localhost:8080/assets' \
--header 'Authorization: Basic Ym9iOmx1Y2t5'
```

check tracy's assets

```curl
curl --location 'localhost:8080/assets' \
--header 'Authorization: Basic dHJhY3k6YmVsbG93'
```

create an order for bob
bob wants to spend 1200 USD to buy 1000 EUR

```curl
curl --location 'localhost:8080/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic Ym9iOmx1Y2t5' \
--data '{
    "side":"buy",
    "price":1.2,
    "quantity":1200,
    "asset_pair":"EUR-USD"
}'
```

check bobs assets, bob has committed 1200 USD to the order

```curl
curl --location 'localhost:8080/assets' \
--header 'Authorization: Basic Ym9iOmx1Y2t5'
```

create a matching order for tracy with a different side
tracy wants to sell 1000 EUR to buy 1200 USD

```curl
curl --location 'localhost:8080/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic dHJhY3k6YmVsbG93' \
--data '{
    "side":"sell",
    "price":1.2,
    "quantity":1200,
    "asset_pair":"EUR-USD"
}'

```

check tracy's assets, tracy has sold 1000 EUR and gained 1200 USD

```curl
curl --location 'localhost:8080/assets' \
--header 'Authorization: Basic dHJhY3k6YmVsbG93'
```

check bobs assets, bob has been paid 1000 EUR

```curl
curl --location 'localhost:8080/assets' \
--header 'Authorization: Basic Ym9iOmx1Y2t5'
```

get bobs orders

```curl
curl --location 'localhost:8080/orders' \
--header 'Authorization: Basic Ym9iOmx1Y2t5'
```

get traceys orders

```curl
curl --location 'localhost:8080/orders' \
--header 'Authorization: Basic dHJhY3k6YmVsbG93'
```

orders are now settled and filled.
