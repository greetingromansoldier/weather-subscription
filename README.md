# weather-subscription

get regular weather updates for your city.

this app gives you an api to:
1.  get current weather for a city.
2.  subscribe to get weather updates (hourly or daily) by email.

it uses weatherapi.com for live weather and an sqlite database for subscriptions.

## what it does

*   **get weather:** see current weather for any city.
*   **subscribe:** sign up with your email, city, and how often you want updates (hourly/daily).
*   **email check:** new subscriptions need you to click a link (we just show the link in the console for now).
*   **unsubscribe:** stop getting updates using a special link.
*   **docker ready:** has a `dockerfile` and `docker-compose.yml` to run it in a container.

## how the files are set up

```
.
├── cmd/application/      # main app code
│   └── main.go
├── handlers/             # handles web requests
│   ├── subscription_handler.go
│   └── weather_handler.go
├── services/             # talks to outside services (like weatherapi)
│   └── weather_service.go
├── models/               # how data looks (weather, subscription)
│   ├── subscription.go
│   └── weather.go
├── storage/              # database stuff (sqlite)
│   └── sqlite.go
├── go.mod                # go project file
├── go.sum                # go project file
├── .env                  # api keys, port (not in git)
├── Dockerfile            # builds the docker image
├── docker-compose.yml    # runs with docker compose
└── README.md             # this file
```

## what you need

*   go (version 1.24.1 or similar)
*   api key from [weatherapi.com](https://www.weatherapi.com/)
*   docker & docker compose (if you want to use containers)

## how to run it

### 1. api key & port

*   make a file named `.env` in the main project folder.
*   put your weatherapi.com key in it:
    ```env
    WEATHERAPI_KEY="YOUR_WEATHERAPI_COM_API_KEY"
    PORT="8080" # optional, uses 8080 if you don't set it
    DB_PATH="./subscriptions.db"
    ```
*   **important:** `.env` is ignored by git. don't share it.

### 2. run it on your computer (no docker)

```bash
# get needed code parts
go mod tidy

# start the app
go run ./cmd/application/main.go
```
the server should start (usually on port 8080).

### 3. run it with docker compose (easier)

uses `dockerfile` and `docker-compose.yml`.

```bash
# build and run the docker container
docker-compose up --build
```
to run it in the background:
```bash
docker-compose up -d --build
```
to stop it:
```bash
docker-compose down
```

## api links

*   **`get /weather?city=<city_name>`**
    *   gets current weather.
    *   example: `http://localhost:8080/weather?city=london`

*   **`post /subscribe`**
    *   starts a new subscription.
    *   send this json data:
        ```json
        {
          "email": "user@example.com",
          "city": "paris",
          "frequency": "daily" // or "hourly"
        }
        ```
    *   you'll get a message saying a confirmation email was "sent" (link shows in console).

*   **`get /confirm/<confirmation_token>`**
    *   confirms your subscription using the token from the "email."
    *   example: `http://localhost:8080/confirm/some-long-confirmation-token`

*   **`get /unsubscribe/<unsubscribe_token>`**
    *   stops your subscription using the token you get after confirming.
    *   example: `http://localhost:8080/unsubscribe/some-long-unsubscribe-token`

## database

*   uses an sqlite database (`subscriptions.db`) for subscription info.
*   it's made in the main folder if you run locally, or in a docker volume with docker compose.

## quick notes

*   email sending is just faked for now: links are shown in the server console.
*   sending actual weather update emails isn't done yet.
*   the `weather` data currently has a `city` field for the `/weather` link. this might change later as city info is mostly for subscriptions.
