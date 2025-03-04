# Распределённый вычислитель арифметических выражений
Проект **disperser-calculator** состоит из двух сервисов:
1. **Оркестратор (orchestrator)**:
   - Принимает на вход арифметическое выражение.
   - Разбивает его на задачи (в упрощённом виде — одна задача на всё выражение).
   - Хранит информацию о статусе выражений.
   - Отдаёт задачи агентам и принимает от них результаты.

2. **Агент (agent)**:
   - Периодически опрашивает оркестратора по эндпоинту `GET /internal/task`, получает задачу.
   - Вычисляет задачу (использует парсер арифметических выражений).
   - Отправляет результат обратно на `POST /internal/task`.

- Установка и запуск

Требуется Go 1.18+ (или выше).

```bash
git clone https://github.com/horhhe/disperser-calculator.git
cd go-dispcalc1
```

-- Запуск Оркестратора

```bash
go run ./cmd/orchestrator/...
```
-- Запуск Агента
В другом окне терминала запустите:

```bash
cd go-dispcalc1
go run ./cmd/agent/...
```
- Примеры запросов
- 
*запросы отправляются в третье окно git bash

Добавить новое выражение на вычисление
```bash
curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": "2+2*2"}'
```

Ответ (201):

```
json
{"id": "1"}
```

(Где "1" — некий ID выражения.)

-- Ошибка 422 (некорректное выражение):

Отправка запроса с пустым выражением, что не считается валидным.
```
curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": ""}'
```
-- Ошибка 500 (кодовое слово trigger500):
Отправка запроса с выражением, равным ```"trigger500"```, что инициирует панику в сервере и приводит к ошибке 500.
```
curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"expression": "trigger500"}'
```
-- Получить список всех выражений


 ```curl --location 'http://localhost:8080/api/v1/expressions'```


Ответ (200):

```json

{
  "expressions": [
    {
      "id": "1",
      "status": "done",
      "result": 6
    }
  ]
}
```

- Получить выражение по ID
```curl --location 'http://localhost:8080/api/v1/expressions/1'```
Ответ (200):
```
{
  "expression": {
    "id": "1",
    "status": "done",
    "result": 6
  }
}
```
Если выражения нет, то статус 404.

#### Эндпоинты для агента (внутренние)
```GET /internal/task```
— получить задачу.
``` POST /internal/task ```
— отдать результат.

## Запуск тестов
```bash
go test ./...
```
Будут запущены все тесты из пакета test, а также тесты из других подпакетов (при наличии).
