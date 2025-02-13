# Бэкенд системы бронирования фильмов

## Обзор
Бэкенд-сервис на Go для системы бронирования фильмов, который обрабатывает аутентификацию пользователей, управление фильмами, расписание сеансов и бронирование мест.

## Фичи

### Аутентификация и авторизация
- Регистрация и вход пользователей
- Аутентификация на основе JWT
- Контроль доступа на основе ролей (админ/юзер)
- Привилегии администратора для управления системой

### Управление фильмами
- Получение списка, создание, обновление и удаление фильмов
- Информация фильма включает в себя название, описание, жанр и постер

### Управление сеансами
- Расписание сеансов фильмов
- Отслеживание вместимости и бронирования мест
- Просмотр доступных мест для каждого сеанса

### Система бронирования
- Бронирование мест на конкретные сеансы
- Отмена бронирования
- Просмотр истории бронирований пользователя
- Предотвращение двойного бронирования мест
- Отслеживание общего количества забронированных мест

### Мониторинг дохода
- Расчет дохода за каждый фильм
- Отслеживание общего дохода системы
- Мониторинг заполняемости мест

## API-эндпоинты

### Аутентификация
- `POST /auth/signup` - Регистрация нового пользователя
- `POST /auth/login` - Вход пользователя

### Фильмы
- `GET /movies` - Список всех фильмов
- `POST /movies/add` - Добавление нового фильма (Администратор)
- `PUT /movies/update/{id}` - Обновление фильма (Администратор)
- `DELETE /movies/delete/{id}` - Удаление фильма (Администратор)

### Сеансы
- `GET /showtimes` - Список всех сеансов
- `POST /showtimes/add` - Добавление нового сеанса (Администратор)
- `PUT /showtimes/update/{id}` - Обновление сеанса (Администратор)
- `DELETE /showtimes/delete/{id}` - Удаление сеанса (Администратор)
- `GET /showtimes/seats/{id}` - Получение доступных мест

### Бронирования
- `POST /reserve/add` - Создание бронирования
- `DELETE /reserve/delete/{id}` - Отмена бронирования
- `GET /reserve` - Получение бронирований пользователя
- `GET /reserve/all` - Получение всех бронирований (Администратор)
- `GET /reserve/movie/{id}` - Получение бронирований по фильму (Администратор)

### Доходы
- `GET /revenue` - Получение статистики общего дохода (Администратор)

## Технологический стек
- Язык: Go
- База данных: PostgreSQL
- Аутентификация: JWT
- Хеширование паролей: bcrypt

## Настройка

1. Установите зависимости:

```bash
go mod tidy
```

2. Создайте файл `.env` в корневой директории со следующими переменными:

```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
```

3. Запустите сервер:

```bash
go run main.go
```

## Схема базы данных
Система использует PostgreSQL с таблицами:
- users (id, username, password_hash, role)
- movies (id, title, description, genre, poster_image)
- showtimes (id, movie_id, start_time, capacity, reserved)
- reservations (id, user_id, movie_id, showtime_id, seats)

## Функции безопасности
- Хеширование паролей с использованием bcrypt
- Аутентификация на основе JWT
- Контроль доступа на основе ролей
- Защищенные маршруты для администратора

## Тестирование
Запустите тесты с помощью:

```bash
go test
```