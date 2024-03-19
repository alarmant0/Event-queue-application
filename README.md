# Event Queue Application

Event Queue Application is a Dockerized Event-Driven Architecture featuring Nginx, Redis, and Go applications.

## Prerequisites

Ensure you have the following installed on your system:

- Docker
- Docker Compose (version 1.29.2)

## Usage Instructions

### 1. Cloning the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/alarmant0/Event-queue-application
cd Event-queue-application
```

### 2. Building Docker Images, Starting Containers, and Testing the Setup

Execute the following command to build the Docker images, start all containers, and test the setup:

```bash
docker-compose build && docker-compose up -d && curl -I -H "Host: producer" -X POST http://localhost:8080/publish
```

### Publishing Messages

To publish 100 messages to the producer, use the following command:

```bash
curl -H "Host: producer" -X POST http://localhost:8080/publish/100
````

### Viewing Logs

To view the logs of the container, use the following command:

```bash
docker-compose logs producer -f
docker-compose logs consumer -f
```

### Monitoring Message Consumption

To monitor message consumption by the consumer, use the following command:

```bash
watch -n 0.1 docker exec -it redis redis-cli LLEN esilv
