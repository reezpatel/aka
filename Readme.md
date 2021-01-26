# A.K.A

#### A Dead Simple URL Shortener written in GO

AKA is an in memory url Shortener. It is entirely written in GO with zero dependencies. All data is stored in memory. Data is persisted to a json file at regular interval and on exit.

#### Environment Variables

| Name       | Description                                                                      |
| ---------- | -------------------------------------------------------------------------------- |
| DATA_PATH  | Folder which will be used for persisting data (if not set, auth is not required) |
| AUTH_TOKEN | Secret used with API to modify records                                           |

### Run Application

##### Run Natively

You can directly run application using `go run .`
Application will run on port `3000`

##### Run in Docker

Use docker to run it as a container

```
docker create --name aka -p 3000:3000 reezpatel/aka
```

##### Run with Docker Compose

Use the included `docker-compose.yml` file, and run `docker-compose up -d`

#### API Documentation

1. Handle redirects

```
curl --request GET \
  --url http://localhost:3000/abcs
```

2. Add new entry

```
curl --request POST \
  --url http://localhost:3000/ \
  --header 'Authorization: my-secret-token' \
  --header 'Content-Type: application/json' \
  --data '{
	"url": "https://www.example.com",
	"aka": "xyz"
}'
```

3. Path an entry

```
curl --request PATCH \
  --url http://localhost:3000/ \
  --header 'Authorization: my-secret-token' \
  --header 'Content-Type: application/json' \
  --data '{
	"url": "https://www.google.com",
	"aka": "abc"
}'
```

4. Delete an entry

```
curl --request DELETE \
  --url http://localhost:3000/ \
  --header 'Authorization: my-secret-token' \
  --header 'Content-Type: application/json' \
  --data '{
	"aka": "abc"
}'
```

#### Development

AKA is an dead simple url shortener, it is written in Go and doesn't have any dependencies. To contribute to development clone the repo and run `go run .`
