# Yet Another Delivery Backend

<!-- ToC start -->
## Содержание

1. [Ссылки](#Ссылки)
1. [Запуск](#Запуск)
1. [Схема БД](#Схема-БД)
1. [Структура проекта](#Структура-проекта)
1. [Правила](#Правила)
<!-- ToC end -->

## Ссылки

:bookmark_tabs: [Доска задач](https://github.com/orgs/MAVIKE/projects/2)

:notebook: [Документация](https://github.com/MAVIKE/yad-docs)

:iphone: [Android](https://github.com/MAVIKE/yad-android)

:phone: [iOS](https://github.com/MAVIKE/yad-ios)

## Запуск

### Локальный запуск

Перед запуском необходимо установить локальные настройки БД в файле _configs/config.yml_,
который генерируется из _configs/config.yml.example_ командой ```make config```

```
make run
```

### Запуск с помощью Docker

```
make docker_build
make docker_run
```

**1. Проверка:**

```
http://localhost:9000/api/v1/ping
```

**2. При первом запуске нужно выполнить скрипт _schema/init.sql_
внутри контейнера с БД.**

**3. Доступные эндпоинты после запуска можно посмотреть по адресу:**

```
http://localhost:9000/swagger/index.html
```

## Схема БД

![](docs/img/db-schema.svg)

## Структура проекта

```
.
├── internal
│   ├── app          // инициализация проекта
│   ├── domain       // основные структуры
│   ├── delivery     // обработчики запросов
│   ├── service      // бизнес-логика
│   └── repository   // взаимодействие с БД
├── cmd              // точка входа в приложение
├── schema           // SQL файлы с миграциями
├── configs          // файлы конфигурации
├── docs             // документация
└── .github          // файлы настройки Github Actions
```

## Правила

Перед тем как коммитить изменения выполните ```make lint```.

### Ветки

Каждый новый тикет (issue) следует выполнять в отдельной ветке с префиксом **fb-N-**,
где **N** - номер тикета. После в названии следует краткая информация о задаче.

Например,
тикет [#1 Проектирование БД](https://github.com/MAVIKE/yad-backend/issues/1),
ветка [fb-1-db-schema](https://github.com/MAVIKE/yad-backend/tree/fb-1-db-schema).

### Коммиты

Коммиты в ветке должны начинаться с **#N**.

Пример для ветки выше - "#1 Update DB schema picture".

### Запросы на слияние

После выполнения задания надо назначить Pull Request (PR) в ветку **develop**.

PR содержит название тикета, в описании указывается
[связь с ним](https://docs.github.com/en/github/managing-your-work-on-github/linking-a-pull-request-to-an-issue).

[Пример PR](https://github.com/MAVIKE/yad-backend/pull/2).
