
# Microservices

A project exploring microservices architecture through the separation of inventory management and statistics into independent services.

The project focuses on modular backend design, separation of responsibilities, and organizing application functionality into individual services.

## Project Structure

```text
Microservices/
├── inventory_service/
├── statictics_service/
└── README.md
```

## Services

### 1. Inventory Service

The `inventory_service` module is dedicated to inventory-related functionality.

It represents a separate component responsible for the inventory domain, keeping its business logic isolated from other parts of the application.

### 2. Statistics Service

The `statictics_service` module is dedicated to statistics-related functionality.

It separates statistical processing from inventory management, demonstrating how different responsibilities can be organized into individual services.

## Architecture

The project follows a microservices-oriented structure, dividing functionality into separate components.

**Key concepts:**
- **Separation of Concerns** — Each service focuses on a specific responsibility.
- **Modularity** — Functionality is organized into independent modules.
- **Maintainability** — Service-specific logic can be modified without reorganizing the entire application.
- **Extensibility** — The structure supports the addition of new services.

## Project Purpose

This project was developed to explore fundamental backend architecture concepts, including:

- Microservices design principles.
- Separation of business logic.
- Modular application structure.
- Organization of backend services.
