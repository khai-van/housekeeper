# housekeeper

This project is a microservices-based system for connecting housekeepers (cleaning service providers) with customers who require home cleaning services. When a customer posts a cleaning job on the system, the job is sent to a housekeeper.

## Architecture

The system is built using a microservices architecture, consisting of the following services:

*   **Booking Service:**  Accepts job postings from customers, calls the Pricing Service to determine the job price, saves job information to a MongoDB database, and triggers the send process by invoking the Send Service.
*   **Pricing Service:**  Determines the price of a job based on the job date and other factors (e.g., start date, hour per job). The pricing logic is configurable and can be updated.
*   **Send Service:**  Exposes a gRPC API for send jobs. Publishes job IDs to a RabbitMQ queue for asynchronous processing.
*   **Send Worker Service:**  Consumes job from the RabbitMQ queue and send the jobs to suitable housekeepers.

## Technologies

*   **Programming Language:** Golang
*   **Microservice Communication:** gRPC
*   **Message Queue:** RabbitMQ
*   **Database:** MongoDB
*   **REST API:** Echo (golang)
*   **Containerization:** Docker
*   **Orchestration:** Docker Compose 

## Getting Started (Local Development)

### Prerequisites

*   Go (version 1.21 or later)
*   Docker
*   Docker Compose
*  Protoc

### Installation

1.  **Clone the repository:**

    ```bash
    git clone <repository_url>
    cd housekeeper
    ```

2.  **Start the services using Docker Compose:**

    ```bash
    docker-compose up -d --build
    ```


### Accessing the Services

*   **Booking Service:** `http://localhost:8080`
*   **Pricing Service:** gRPC on port 8081
*   **Send Service:** gRPC on port 8082

