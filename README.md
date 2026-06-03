# Bank Transaction API

A simple and open-source Bank Transaction API for handling basic banking operations such as balance inquiry, create transaction, transaction history, download report, and audit trail.
This project is designed to be easy to understand, easy to integrate, and suitable as a foundation for backend financial services.

## Features

- Balance inquiry  
- Transaction history  
- Create Transaction
- Download report  
- RESTful API endpoints  
- Lightweight and easy to extend  

## Use Cases

- Learning and prototyping banking APIs  
- Backend services for financial applications  
- Internal tools and microservices  

## Getting Started

Check the documentation below for setup instructions and API usage examples.

## Quick Start (Run in 3 Steps)

1. **Clone the repository**
   ```bash
   git clone https://github.com/your-username/bank-transaction-api.git
   cd bank-transaction-api
   
2. **Configure environment variables**
   ```bash
   cp .env.example .env

3. **Run The application**
   ```bash
   go run main.go
4. **The API will be available at:**
   ```bash
   http://localhost:6969


## Docker Quick Start

Make sure **Docker** is installed on your machine.

### 1. Build the Docker image
```bash
docker build -t bank-transaction-api .
```

### 2. Run the container
```bash
docker run -d \
  -p 6969:6969 \
  --env-file .env \
  --name bank-transaction-api \
  bank-transaction-api
```

### 3. The API will be accessible at:
```bash
http://localhost:6969
```

### 4. Optional: Docker Compose
```bash
docker-compose up -d
```
## Sample cURL Request

1. **Inquiry Balance**
   ```bash
   curl --location 'localhost:6969/account/inquiry' \
      --header 'Content-Type: application/json' \
      --header 'Authorization: Bearer YOUR_TOKEN' \
      --data '{
          "username": "your_username",
          "password": "your_password",
          "request": {
              "account_no": "030881899288083"
          }
      }'

2. **Create Transaction**
    ```bash
   curl --location 'localhost:6969/transaction/create' \
   --header 'Content-Type: application/json' \
   --header 'Authorization: Bearer YOUR_TOKEN' \
   --data '{
       "username": "your_username",
       "password": "your_password",
       "request": {
           "rowid_sender": 1,
           "rowid_beneficiary": 1,
           "debet_account": "020601050203045",
           "debet_name": "HENXX XXXMAD",
           "debet_curr": "IDR",
           "credit_account": "020357383928345",
           "credit_name": "JOXX XXE",
           "credit_curr": "IDR",
           "amount": 1500000.75,
           "remark": "Payment invoice February"
       }
   }'
