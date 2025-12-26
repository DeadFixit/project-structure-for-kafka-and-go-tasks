````markdown

Задание 1 — Kafka playground (локально через Docker)

Этот пакет содержит пошаговые инструкции для выполнения практики с Kafka локально через Docker Compose (см. `docker-compose.yml`). Все команды записаны для PowerShell 7.5.4 на Windows 11.

Описание работы

Целью задания является изучение принципов работы с системой публикации-подписки (pub-sub) на основе Apache Kafka. В процессе выполнения задания вы развернете топик Kafka, отправите и получите тестовые сообщения, а также проанализируете метаданные брокера.

Как запустить

1. Установить Docker

Убедитесь, что у вас установлен Docker Desktop для Windows. После установки запустите следующую команду в каталоге проекта для запуска контейнеров:

```powershell
docker compose up -d
````

Это развернет Zookeeper и Kafka в контейнерах.

### 2. Установить переменную окружения

Убедитесь, что переменная окружения `BOOTSTRAP_SERVER` установлена для правильного подключения к Kafka:

```powershell
$Env:BOOTSTRAP_SERVER="localhost:9092"
```

### 3. Выполнить команды из чек-листа

После того как Kafka и Zookeeper запустятся, вы можете выполнять команды из следующего чек-листа для создания топика, отправки сообщений, чтения сообщений и анализа метаданных.

---

## Чек-лист команд (PowerShell 7.5.4)

> В PowerShell команды выполняются на хосте, но бинарники Kafka вызываются через `docker compose exec kafka ...`. Переменные окружения читаются через `$env:VAR`.

### 1. Проверка списка топиков (скриншот №1):

Убедитесь, что Kafka работает и топики доступны:

```powershell
docker compose exec kafka kafka-topics --bootstrap-server $env:BOOTSTRAP_SERVER --list
```

### 2. Создание топика (скриншот №2):

Создайте топик `student-topic` с 2 партициями и репликацией 1:

```powershell
docker compose exec kafka kafka-topics --bootstrap-server $env:BOOTSTRAP_SERVER --create --topic student-topic --partitions 2 --replication-factor 1
```

### 3. Описание топика:

Для того чтобы увидеть структуру топика, выполните команду:

```powershell
docker compose exec kafka kafka-topics --bootstrap-server $env:BOOTSTRAP_SERVER --describe --topic student-topic
```

### 4. Отправка первого сообщения (скриншот №3):

Отправьте тестовое сообщение в топик:

```powershell
echo "Hello from student Иванов И.И." | docker compose exec -T kafka kafka-console-producer --bootstrap-server $env:BOOTSTRAP_SERVER --topic student-topic
```

### 5. Отправка второго сообщения:

Отправьте второе сообщение для демонстрации нескольких записей:

```powershell
echo "Data event: 2025-09-16" | docker compose exec -T kafka kafka-console-producer --bootstrap-server $env:BOOTSTRAP_SERVER --topic student-topic
```

### 6. Чтение всех сообщений с начала:

Для чтения всех сообщений с начала выполните:

```powershell
docker compose exec -T kafka kafka-console-consumer --bootstrap-server $env:BOOTSTRAP_SERVER --topic student-topic --from-beginning
# Нажмите Ctrl+C для выхода
```

### 7. Метаданные брокера (скриншот №4):

Проверьте метаданные брокера с помощью команды:

```powershell
docker compose exec kafka kafka-broker-api-versions --bootstrap-server $env:BOOTSTRAP_SERVER
```

### 8. Проверка переменной окружения:

Если возникают проблемы с переменной окружения, вы можете проверить значение переменной `BOOTSTRAP_SERVER`:

```powershell
echo $env:BOOTSTRAP_SERVER
```

---
