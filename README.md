# Habr-comments-server

## Info

**Данный сервис реализует систему добавления и чтения постови комментариев с использованием _GraphQL_, аналогичную комментариям к постам на популярных платформах, таких как Хабр или Reddit.**

_Характеристики системы постов:_
- Можно просмотреть список постов.
- Можно просмотреть пост и комментарии под ним.
- Пользователь, написавший пост, может запретить оставление комментариев к своему посту.


_Характеристики системы комментариев к постам:_
- Комментарии организованы иерархически, позволяя вложенность без ограничений.
- Длина текста комментария ограничена до, например, 2000 символов.
- Система пагинации для получения списка комментариев.

*Сервис работает на port 8082*

Запуск при помощи *Makefile*:

```bash
make db-start
#если нужно поднять миграции
make migration-up
make start
```


### Запуск с выбором

```bash
go build -o server ./path/to/main.go
./server -in-memory #если хотите запустить с in-memory
./server #если хотите запустить с PostgreSQL
```

### Запуск с Docker-Compose

```bash
docker-compose -f ./path/to/docker-compose.yml -p habr-comments-server up -d
```

### Дополнительные команды:

- Чтобы остановить контейнеры:

```bash
docker-compose -f ./docker-compose.yml -p habr-comments-server down
```

- Чтобы проверить логи контейнеров:

```bash
docker-compose -f ./docker-compose.yml -p habr-comments-server logs -f
```

### Запуск Dockerfile

```bash
# Сборка Docker-образа
docker build -t habr-comments-server:latest -f ./path/to/Dockerfile .

# Запуск контейнера
docker run -d --name habr-comments-server-container -p 8082:8082 habr-comments-server:latest

# Проверка работы контейнера
docker ps

# Чтобы проверить логи контейнера
docker logs habr-comments-server-container

# Остановка контейнера
docker stop habr-comments-server-container

```

