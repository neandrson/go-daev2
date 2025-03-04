Для связи tg: @neandrs

# go-daev2
Распределенный вычислитель арифметических выражений.


## Проект

Проект состоит из 2 элементов:

- Orchestrator - сервер, который принимает арифметическое выражение, переводит его в
  набор последовательных задач и обеспечивает порядок их выполнения. 
- Agent - вычислитель, который может получить от "оркестратора" задачу, выполнить его и
  вернуть серверу результат.

<img src="misc/Diagram.png" alt=dia widt="900"/>

## Сборка и запуск

Склонируйте репозиторий
```sh
git clone https://github.com/neandrson/go-daev2.git
```
Перейдите в корневой каталог проекта.

### Docker

Собрать приложение состоящее из оркестратора и трех агентов:
```sh
docker compose build && docker compose up
```

Пожалуйста, дождитесь сообщения о том, что сервер начал работу.
На моей машине этот процесс занимает заметное время.


Для того, чтобы запустить контейнеры в фоновом режиме:
```sh
docker compose up -d
```

Для того, чтобы остановить запущенные контейнеры  в фоновом режиме:
```sh
docker compose down
```
или `CTRC-C`, если контейнеры запущены не в фоне.

Чтобы изменять переменные окружения, модифицируйте файл `.env`


### Linux
  - Собрать и запустить оркестратор командой:
  ```sh
  make orchestrator && source scripts/setenv.sh && ./orchestrator
  ```
  - Собрать и запустить agent командой:
  ```sh
  make agent && source scripts/setenv.sh && ./agent
  ```

Для запуска вам понадобится несколько вкладок терминала.
Вы также можете запустить процессы в фоновом режиме, используя символ `&`.
Потоки вывода следует перенаправить в файл, например `/dev/null`.
Например, чтобы запустить агент в фоновом режиме:
  ```sh
  make agent && source scripts/setenv.sh && ./agent 2>1 >/dev/null &
  ```
В таком случае не забудьте, остановить процессы, командой `kill -s SIGINT proccess_id`.
PID процесса можно узнать, с помощью команды `ps`

Чтобы изменять переменные окружения, модифицируйте файл `scripts/setenv.sh`

## Примеры

Все примеры рассчитаны на то, что приложение собиралось через docker compose.
В ином случае указываете порт, который прослушивает сервер - 8081.

### Чтобы добавить арифметическое выражение на вычисление, используйте запросы вида:

```sh
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
      "expression": "строка_с_выражением"
}'
```

Например:

```sh
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
      "expression": "1 + 2 + 3 * -4 - 20 / 5"
}'
```

Возможные ошибки:
- вы ввели пустое expression
- вы ввели некорректное арифметическое выражение

### Чтобы получить список всех выражений, используйте запрос:

```sh
curl --location 'localhost/api/v1/expressions'
```

Пример:
```
{
    "expressions": [
        {
            "id": "1",
            "status": "Error",
            "result": "",
        },
        {
            "id": "2",
            "status": "Done",
            "result": "3",
        },
        {
            "id": "3",
            "status": "In process",
            "result": ""
        }
    ]
}
```

Всего возможно три статуса:
- `Error`:      найдена синтаксическая или грамматическая ошибка в арифметическом выражении
- `In process`: выражение в процессе вычисления
- `Done`:       получен результат выражения

Стоит отметить, что деление на нуль не является ошибкой и будет получен результат.
Это может быть `+inf`, `-inf` (бесконечность) или `NaN` (not a number), если числитель
также был равен или близок к нулю.

### Чтобы получить выражение по его идентификатору, используйте запрос:

```sh
curl --location 'localhost/api/v1/expressions/{id}'
```

Например,

```sh
curl --location 'localhost/api/v1/expressions/2'
```
```sh
{
    "expression": {
        "id": "2",
        "status": "Done",
        "result": "3"
    }
}
```
или

```sh
curl --location 'localhost/api/v1/expressions/4'
id "4" not found
```

#### Получение задачи для выполения.
 
```
curl --location 'localhost/internal/task'
 
```

Тело ответа

```sh
{
  "task":
    {
      "id": <идентификатор задачи>,
      "arg1": <имя первого аргумента>,
      "arg2": <имя второго аргумента>,
      "operation": <операция>,
      "operation_time": <время выполнения операции>
    }
}
 
```

Например:

```sh
{
  "task":
    {
      "id": <1>,
      "arg1": <2>,
      "arg2": <2>,
      "operation": <+>,
      "operation_time": <0>
    }
}
```
