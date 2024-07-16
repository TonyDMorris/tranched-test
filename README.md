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

check bobs assets

```curl
curl --location 'localhost:8080/assets' \
--header 'Authorization: Basic Ym9iOmx1Y2t5'
`
```

check tracy's assets

```curl
curl --location 'localhost:8080/assets' \
--header 'Authorization: Basic dHJhY3k6YmVsbG93'
```

create an order for bob

```curl
curl --location 'localhost:8080/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic Ym9iOmx1Y2t5' \
--data '{
    "side":"buy",
    "price":1.2,
    "quantity":1000,
    "asset_pair":"EUR-USD"
}'
```

create a matching order for tracy with a different side

```curl
curl --location 'localhost:8080/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic dHJhY3k6YmVsbG93' \
--data '{
    "side":"sell",
    "price":1.2,
    "quantity":1000,
    "asset_pair":"EUR-USD"
}'

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
