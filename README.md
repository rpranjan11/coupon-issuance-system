# Coupon Issuance System

A high-performance system for issuing a limited number of coupons on a first-come-first-served basis at a specific time.

## Overview

This system enables creating coupon campaigns with configurable parameters. Each campaign specifies the number of available coupons and a specific start date and time. The system ensures that exactly the specified number of coupons are issued, with no excess issuance, and each coupon has a unique code.

## Features

- Create campaigns with a specified number of coupons and start time
- Get campaign information including all issued coupon codes
- Issue coupons on a first-come-first-served basis
- Concurrent request handling with data consistency
- Unique coupon code generation with Korean characters and numbers
- Request validation for All APIs 
- 

## Prerequisites

- Go 1.21 or later
- Protocol Buffer Compiler (protoc) 3.15.0 or later

## Installation

### 1. Clone the repository

git clone https://github.com/rpranjan11/coupon-issuance-system.git
cd coupon-issuance-system

### 2. Install the Protocol Buffer Compiler

#### macOS (using Homebrew)

brew install protobuf

#### Linux

apt-get install -y protobuf-compiler

#### Windows

Download the pre-built binary from [GitHub releases](https://github.com/protocolbuffers/protobuf/releases)

### 3. Install Go Protocol Buffer plugins

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

### 4. Generate Protocol Buffer code

protoc --go_out=. --go_opt=paths=source_relative --connect-go_out=. --connect-go_opt=paths=source_relative coupon.proto

### 5. Install dependencies

go mod tidy

### 6. Build the server

go build -o server ./cmd/server

### 7. Run the server

./server


## API Endpoints
### 1. Create Campaign
- **Endpoint**: `/CreateCampaign`
- **Method**: `POST`
- **Request Body**:
```json
{
  "name": "string",
  "start_time": "2025-05-11T00:00:00Z",
  "coupon_count": 100
}
```

- **Response**:
```json
{
  "id": "string",
  "name": "string",
  "totalCoupons": 100,
  "startTime": "2025-05-11T00:00:00Z",
  "createdAt": "2025-05-11T16:23:13.093881Z"
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
    "issuedAt": "2025-05-10T16:25:07.607675Z"
  }
}
```

### 3. Get Campaign Details
- **Endpoint**: `/GetCampaign/`
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
    "startTime": "2025-05-10T00:00:00Z",
    "createdAt": "2025-05-10T16:23:56.405350Z"
  },
  "coupons": [
    {
      "code": "string",
      "campaignId": "string",
      "issuedAt": "2025-05-10T16:25:07.607675Z"
    }
  ]
}
```

### 4. Delete Campaign
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

A Postman collection is provided in the `postman` directory. You can import it into Postman to test the API endpoints.
### I have added a postman collection for testing the API endpoints. You can find it in the `postman` directory. The collection includes requests for creating campaigns, issuing coupons, and retrieving campaign information.