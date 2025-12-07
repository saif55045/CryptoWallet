# Crypto Wallet Backend (Go)

A decentralized cryptocurrency wallet backend built with Go, featuring blockchain implementation, UTXO model, and Proof-of-Work mining.

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or higher
- MongoDB Atlas account

### Installation

1. Clone the repository and navigate to backend:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Create `.env` file from `.env.example`:
```bash
cp .env.example .env
```

4. Update `.env` with your MongoDB Atlas URI and other credentials

### Running the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Building for Production

```bash
go build -o crypto-wallet-server
```

## 📚 API Documentation

### Module 1: Authentication

#### POST `/api/auth/signup`
Create a new user account
```json
{
  "fullName": "John Doe",
  "email": "john@example.com",
  "cnic": "12345-1234567-1",
  "password": "securePassword123"
}
```

#### POST `/api/auth/verify-otp`
Verify email with OTP
```json
{
  "email": "john@example.com",
  "otp": "123456"
}
```

#### POST `/api/auth/resend-otp`
Resend OTP to email
```json
{
  "email": "john@example.com"
}
```

#### POST `/api/auth/login`
Login to account
```json
{
  "email": "john@example.com",
  "password": "securePassword123"
}
```

#### GET `/api/auth/profile`
Get user profile (Protected)
- Headers: `Authorization: Bearer <token>`

## 🗄️ Database Schema

### Users Collection
```javascript
{
  _id: ObjectId,
  fullName: String,
  email: String (unique),
  cnic: String (unique),
  password: String (hashed),
  walletId: String,
  publicKey: String,
  privateKey: String (encrypted),
  beneficiaries: [String],
  isVerified: Boolean,
  otp: String,
  otpExpiry: Date,
  createdAt: Date,
  updatedAt: Date
}
```

## 🔒 Security Features
- Password hashing with bcrypt
- JWT authentication
- Email verification with OTP
- CORS protection
- Environment variable configuration

## 📦 Project Structure
```
backend/
├── config/          # Configuration files
├── controllers/     # Request handlers
├── database/        # Database connection
├── middleware/      # Custom middleware
├── models/          # Data models
├── routes/          # API routes
├── utils/           # Utility functions
├── main.go          # Entry point
└── .env             # Environment variables
```
