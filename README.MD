# Weather Service

Сервис на Go, который по названию города возвращает текущую температуру, используя [Open-Meteo API](https://open-meteo.com/).

---

## Быстрый запуск через Docker Compose

1. Клонируй репозиторий:

   ```bash
   git clone https://github.com/<your-username>/weather-service.git
   cd weather-service
   ```

2. Подними окружение (PostgreSQL с миграциями):

   ```bash
   docker compose up -d
   ```

   Это создаст базу `weather` и применит SQL-миграции из `./migrations/`.

3. Запусти сервис:

   ```bash
   go run cmd/server/main.go
   ```

4. Проверь:

   ```bash
   curl http://localhost:3000/moscow
   ```

   Если данных пока нет, подожди - фоновая задача обновляет температуру каждые 10 секунд.

---

## Что делает

* Периодически получает координаты и температуру через Open-Meteo API.
* Сохраняет данные в PostgreSQL.
* Отдаёт последние данные по запросу `/город`.

---

## Лицензия

MIT © 2025
