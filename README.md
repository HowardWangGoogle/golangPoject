# Deal Project 

## Overview
This is a basic Go programming project that I developed to practice my coding skills. In the future, I plan to build a decentralized exchange (DEX) for digital currencies based on the concepts explored here.

### Technologies Used
This project incorporates the following tools and technologies:
- **Swagger**: For API documentation and testing.
- **JWT**: To ensure user authentication and reduce unnecessary resource consumption.
- **Docker**: To containerize each component, allowing for independent updates without affecting others.
- **gRPC**: Utilizes Protobuf for efficient message handling between services.
- **AWS** and **Kubernetes (K8S)**: For deployment and scaling.
- **Redis**: Handles priority queues. In future projects, I plan to experiment with alternatives like Kafka or RabbitMQ.
- **PostgreSQL**: As the database.

---

### Project Features
- **Containerization**: Docker is used to separate components, ensuring flexibility and scalability. 
- **Swagger**: Simplifies API testing by displaying methods, parameters, and endpoints directly within the UI. While tools like Postman are also useful, Swagger provides a more integrated approach.
- **gRPC & RESTful APIs**: gRPC is used for high-performance communication, while RESTful APIs run on a separate port for compatibility with other tools.
- **JWT Integration**: Ensures only authenticated users can access certain functionalities, improving efficiency.
- **Redis for Queuing**: Manages priority queues. Plans for integration with Kafka or RabbitMQ are in the pipeline for future iterations.

---

### Getting Started

#### Prerequisites
1. Ensure you have Go installed. Run:
   ```bash
   go mod tidy
   ```

2. Use the provided `Makefile` to simplify commands. To run the application locally:
   ```bash
   go run .
   ```

#### Recommended Setup
For a more robust setup, use Docker:
1. Build and start the application:
   ```bash
   docker-compose -f docker-compose.yaml up --build
   ```

2. To stop the application:
   ```bash
   docker-compose -f docker-compose.yaml down --build
   
