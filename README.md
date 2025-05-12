# Coupon Issuance System

A high-performance system for issuing a limited number of coupons on a first-come-first-served basis at a specific time.

## Overview

This system enables creating coupon campaigns with configurable parameters. Each campaign specifies the number of available coupons and a specific start date and time. The system ensures that exactly the specified number of coupons are issued, with no excess issuance, and each coupon has a unique code.

## Features

- Create campaigns with a specified number of coupons and start time
- Get campaign information including all issued coupon codes
- Issue coupons on a first-come-first-served basis
- Delete campaigns and all associated coupons
- Request validation and error handling
- Generate only the specified number of coupons
- Unique coupon code generation with Korean characters and numbers
- ConnectRPC for efficient communication
- Concurrent request handling with data consistency
- High Traffic Handling: The system is designed to handle high traffic efficiently, ensuring that multiple requests can be processed concurrently without performance degradation.
- Edge Case Handling: The system includes checks to prevent issuing more coupons than available, ensuring that the number of issued coupons does not exceed the total number of coupons in a campaign.
- Error Handling: The system provides clear error messages for invalid requests, such as missing fields or invalid campaign IDs, making it easier to debug and resolve issues.
- Postman Collection: A Postman collection is included for easy testing of the API endpoints, allowing developers to quickly test and validate the functionality of the system.
- Testing Scripts for Concurrency: The system includes testing tools to simulate high traffic and concurrent requests, allowing developers to test the performance and reliability of the system under load.

## Prerequisites

- Go 1.21 or later
- Protocol Buffer Compiler (protoc) 3.15.0 or later

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/rpranjan11/coupon-issuance-system.git

cd coupon-issuance-system
```

### 2. Install the Protocol Buffer Compiler

#### macOS (using Homebrew)

```bash
brew install protobuf
```

#### Linux

```bash
apt-get install -y protobuf-compiler
```

#### Windows

Download from [GitHub releases](https://go.dev/doc/install)

### 3. Install Go Protocol Buffer plugins

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest
```

### 4. Generate Protocol Buffer code

```bash
protoc --go_out=. --go_opt=paths=source_relative --connect-go_out=. --connect-go_opt=paths=source_relative coupon.proto
```

### 5. Install dependencies

```bash
go mod tidy
```

## Server

### 1. Build the server

```bash
go build -o server ./cmd/server
```

### 2. Run the server

```bash
./server
```

## Client

The client is a command-line tool for interacting with the coupon issuance system. It allows you to create campaigns, issue coupons, and retrieve campaign information.

### 1. Build the client

```bash
go build -o client ./cmd/client
```

### 2. Create a campaign

```bash
./client -command=create -name="Test Campaign" -total=100 -start-in=30s
```

### 3. Issue a coupon

```bash
./client -command=issue -campaign-id=<CAMPAIGN_ID>
```

### 4. Get campaign details

```bash
./client -command=get -campaign-id=<CAMPAIGN_ID>
```

### 5. Delete a campaign and all issued coupons

```bash
./client -command=delete -campaign-id=<CAMPAIGN_ID>
```

## Load Testing

To test the performance of the system under high traffic, you can use the `/test/load/main.go` file. This file contains a simple load testing implementation that simulates multiple concurrent requests to the API endpoints.

### 1. Build the testing tool

```bash
go build -o loadtest ./test/load
```

### 2. Run to create a load testing campaign 

```bash
./client -command=create -name="Load Test Campaign" -total=1000 -start-in=30s
```
#### Note the campaign ID from the output.

### 3. Run the load test

```bash
./loadtest -campaign-id=<CAMPAIGN_ID> -concurrency=50 -rate=500 -duration=10s
```


## API Endpoints
### 1. Create Campaign
- **Endpoint**: `/CreateCampaign`
- **Method**: `POST`
- **Request Body**:
```json
{
  "name": "string",
  "start_time": "2025-05-11T00:00:00+09:00",
  "coupon_count": 100
}
```

- **Response**:
```json
{
  "id": "string",
  "name": "string",
  "totalCoupons": 100,
  "startTime": "2025-05-11T00:00:00+09:00",
  "createdAt": "2025-05-11T16:23:13.093881+09:00"
}
```

### 2. Issue Coupon
- **Endpoint**: `/IssueCoupon`
- **Method**: `POST`
- **Request Body**:
```json
{
  "campaign_id": "string"
}
```
- **Response**:
```json
{
  "success": true,
  "coupon": {
    "code": "string",
    "campaignId": "string",
    "issuedAt": "2025-05-10T16:25:07.607675+09:00"
  }
}
```

### 3. Get Campaign Details
- **Endpoint**: `/GetCampaign`
- **Method**: `POST`
- **Request Body**:
```json
{
  "campaign_id": "string"
}
```
- **Response**:
```json
{
  "campaign": {
    "id": "string",
    "name": "string",
    "totalCoupons": 100,
    "issuedCoupons": 1,
    "startTime": "2025-05-10T00:00:00+09:00",
    "createdAt": "2025-05-10T16:23:56.405350+09:00"
  },
  "coupons": [
    {
      "code": "string",
      "campaignId": "string",
      "issuedAt": "2025-05-10T16:25:07.607675+09:00"
    }
  ]
}
```

### 4. Delete Campaign and all issued coupons
- **Endpoint**: `/DeleteCampaign`
- **Method**: `POST`
- **Request Body**:
```json
{
  "campaign_id": "string"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Campaign deleted successfully"
}
```

## Postman Collection

A Postman collection is provided in the `postman` directory. You can import it into Postman to test the API endpoints. The collection includes requests for creating campaigns, issuing coupons, retrieving campaign information, and deleting campaign along with its all issued coupons.

