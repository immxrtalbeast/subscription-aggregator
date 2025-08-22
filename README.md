#  Subscription aggregator
[![wakatime](https://wakatime.com/badge/user/42cf6868-b638-4d34-9e52-ec8f63476139/project/efb823af-3887-466c-983b-ad1f883e7614.svg)](https://wakatime.com/badge/user/42cf6868-b638-4d34-9e52-ec8f63476139/project/efb823af-3887-466c-983b-ad1f883e7614)

### Используемые технологии:
Swagger, gin, docker

### База данных:
GORM, PostgreSQL, Миграции на SQL.

## Launch methods

Скрипт запуска проекта

```bash
go run cmd/main.go --config=./config/local.yaml
```

# Run with docker
Скопируйте себе docker compose файл и запустите

```bash
docker compose up
```

### Swagger:
Swagger доступен по адресу http://localhost:8080/api/v1/swagger/index.html#/default/post_create


### Что можно добавить?
Трассировку с Jaeger

## Architecture
```
├── cmd 
|    ├── docs(файлы swagger`а)
│    └── main.go
├── config
|    ├── dev.yaml
│    └── local.yaml
├── internal
|    ├── config
|    |     └── config.go
│    ├── controller
│    │    └── subscription_contoller.go
│    ├── domain
│    │    ├── month_year.go
│    │    └── subscription.go
│    ├── lib
│    │    ├── sl 
│    │    └── slogpretty 
│    └── service
│    │    └── subscription
│    │      └── subscription_interactor.go
│    └── storage
│         └── psql
│           └── subscriptionRepo.go
└── migrations
     ├── 001_init.down.sql
     └── 001_init.up.sql
```
