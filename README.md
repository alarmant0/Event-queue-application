# Event Queue Application

Event Queue Application is a Dockerized Event-Driven Architecture featuring Nginx, Redis, and Go applications.

## Prerequisites

Ensure you have the following installed on your system:

- Docker
- Docker Compose (version 1.29.2)

## Usage Instructions

### 1. Cloning the Repository

Clone the repository to your local machine:

git clone https://github.com/alarmant0/Event-queue-application
cd Event-queue-application

### 2. Building Docker Images, Starting Containers, and Testing the Setup

Execute the following command to build the Docker images, start all containers, and test the setup:

```bash
docker-compose build && docker-compose up -d && curl -I -H "Host: producer" -X POST http://localhost:8080/publish
