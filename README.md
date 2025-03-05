# 🚀 Bidding & Auction Services

This project consists of two services: **Bidding Service** and **Auction Service**. The Bidding Service responds to ad requests with bid prices, while the Auction Service selects the highest bid from multiple bidding services.

🔗 **GitHub Repository:** [Bidding_Auction_Go](https://github.com/Yhrone/Bidding_Auction_Go)

## 📌 Overview

### 🎯 Bidding Service
- Receives an **AdRequest** via an HTTP request containing an `AdPlacementId` (a unique string identifying an ad slot or ad spot).
- Responds with an **AdObject** containing:
  - 🆔 `AdID` (Unique identifier for the ad).
  - 💰 `bidPrice` (Randomly generated price for the `AdPlacementId` in USD).
- If the service **chooses not to bid**, it returns **HTTP 204 (No Content)** randomly.
- If the service **bids**, it responds with an **AdObject** and **HTTP 200 (OK)**.

### ⚡ Auction Service
- Calls multiple **Bidding Services** simultaneously (header bidding mechanism).
- Accepts an **AdPlacementId** via an external API.
- Collects all bids (valid responses with **HTTP 200** status code).
- Selects the bid with the **highest bidPrice** for the given `AdPlacementId`.
- If no valid bids are received, it returns **HTTP 204 (No Content)**.
- Implements a **safety circuit** to handle slow or non-responsive bidding services:
  - ⏳ If a bidding service takes **more than 200ms**, its bid is ignored.
  - ⚡ The auction service itself **must always respond within 200ms**.

## 📡 API Endpoints

### 🎲 Bidding Service API
#### 📩 Request
```
POST /bid
Content-Type: application/json
{
  "AdPlacementId": "string"
}
```

#### ✅ Responses
- **200 OK**
```json
{
  "AdID": "unique-ad-id",
  "bidPrice": 2.50
}
```
- **204 No Content** (if no bid is placed)

### 🏆 Auction Service API
#### 📩 Request
```
POST /auction
Content-Type: application/json
{
  "AdPlacementId": "string"
}
```

#### ✅ Responses
- **200 OK** (Highest bid selected)
```json
{
  "AdID": "winning-ad-id",
  "bidPrice": 3.75
}
```
- **204 No Content** (No valid bids received)

## ⚙️ Implementation Notes
- Use **parallel requests** to reach multiple bidding services efficiently.
- Implement **timeout handling** (ignore any bid taking longer than 200ms).
- The auction service must always return within **200ms**.
- 🎲 Randomization is used in bidding to simulate different pricing strategies.

## 🚀 Running the Services
1. Start the **Bidding Service** on multiple instances.
2. Start the **Auction Service**, ensuring it can communicate with the Bidding Services.
3. Test the Auction Service by making requests to its exposed API.

## 🛠️ Testing
- Send multiple **AdRequests** to the Bidding Service to observe varied responses.
- Call the **Auction Service** API with an `AdPlacementId` to verify bid selection.
- Introduce artificial delays in a bidding service to ensure the auction circuit breaker works correctly.


