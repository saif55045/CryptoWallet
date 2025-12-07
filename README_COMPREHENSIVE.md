# 🔐 Decentralized Cryptocurrency Wallet System

A complete blockchain-based cryptocurrency wallet application with user authentication, wallet management, UTXO transactions, mining, zakat calculation, and admin panel. Built with Go, React, and MongoDB.

## 🎯 Project Overview

This is a full-stack decentralized cryptocurrency wallet system that implements a complete blockchain with proof-of-work mining, UTXO (Unspent Transaction Output) model, and advanced features like Merkle root verification and admin panel.

### ✨ Features Completed

#### Core Modules (7/7) ✅
1. **User Authentication** - Signup, OTP verification, Login with JWT
2. **Wallet System** - Generate wallets, view balance, manage multiple wallets
3. **UTXO Model** - Track unspent outputs, prevent double spending
4. **Transactions** - Send transactions with UTXO inputs/outputs
5. **Blockchain** - Mine blocks with Proof-of-Work, Merkle root verification
6. **Zakat Calculation** - Calculate Islamic zakat automatically
7. **Reports & Logs** - Activity logs, transaction history, system logs

#### Bonus Features (3/3) ✅
1. **Google OAuth** - Sign in with Google account
2. **Merkle Root Verification** - Cryptographic block verification
3. **Admin Panel** - Complete admin dashboard with user/transaction management

#### UI/UX ✅
- Modern, clean light theme with Tailwind CSS
- Responsive design (mobile, tablet, desktop)
- Professional gradient backgrounds and animations
- Intuitive navigation and user flows

---

## 🛠 Technology Stack

### Backend
- **Go 1.21+** - Programming language
- **Gin Framework** - HTTP web framework
- **MongoDB Atlas** - NoSQL database
- **JWT** - Authentication tokens
- **bcrypt** - Password hashing
- **Google OAuth 2.0** - Third-party authentication

### Frontend
- **React 19.2.1** - UI library
- **Tailwind CSS** - Styling framework
- **React Router 7.10.1** - Client-side routing
- **Axios** - HTTP client
- **Google Identity Services** - OAuth integration

### Deployment
- **Render** - Backend hosting (native Go deployment)
- **Vercel** - Frontend hosting (SPA with client-side routing)
- **MongoDB Atlas** - Managed database service

---

## 📁 Project Structure

```
Project/
├── backend/                          # Go backend server (Gin + MongoDB)
│   ├── config/
│   │   └── cors.go
│   ├── controllers/
│   │   ├── auth_controller.go        # User authentication
│   │   ├── wallet_controller.go      # Wallet management
│   │   ├── utxo_controller.go        # UTXO transactions
│   │   └── admin_controller.go       # Admin operations (NEW)
│   ├── crypto/
│   │   ├── encryption.go
│   │   └── keys.go
│   ├── database/
│   │   └── connection.go
│   ├── middleware/
│   │   └── auth_middleware.go
│   ├── models/
│   │   ├── user.go
│   │   ├── wallet.go
│   │   └── utxo.go
│   ├── routes/
│   │   ├── auth_routes.go
│   │   ├── wallet_routes.go
│   │   ├── utxo_routes.go
│   │   └── admin_routes.go           # Admin routes (NEW)
│   ├── utils/
│   │   ├── hash.go
│   │   ├── jwt.go
│   │   └── email.go
│   ├── main.go
│   ├── go.mod
│   ├── Dockerfile                    # Docker configuration (NEW)
│   ├── .dockerignore                 # Docker ignore rules (NEW)
│   └── .env
├── frontend/                         # React frontend (Tailwind CSS)
│   ├── public/
│   │   ├── index.html
│   │   ├── manifest.json
│   │   └── robots.txt
│   ├── src/
│   │   ├── pages/
│   │   │   ├── Login.js              # Login with Google OAuth (NEW)
│   │   │   ├── Signup.js
│   │   │   ├── VerifyOTP.js
│   │   │   ├── Dashboard.js
│   │   │   ├── WalletProfile.js
│   │   │   ├── Beneficiaries.js
│   │   │   └── Admin.js              # Admin panel (NEW)
│   │   ├── components/
│   │   │   └── Navbar.js             # Updated with Admin link
│   │   ├── context/
│   │   │   └── AuthContext.js        # Updated with Google OAuth
│   │   ├── services/
│   │   │   └── api.js                # Updated with admin API
│   │   ├── App.js
│   │   ├── App.css
│   │   └── index.js
│   ├── package.json
│   ├── tailwind.config.js
│   ├── postcss.config.js
│   ├── vercel.json                   # Vercel deployment config (NEW)
│   └── .env
├── render.yaml                       # Render deployment config (NEW)
├── QUICKSTART.md                     # Quick start guide
├── README_SETUP.md                   # Old setup guide
└── README.md                         # This file
```

---

## 🚀 Quick Start Guide

### Prerequisites
- Go 1.21+
- Node.js 18+ and npm
- MongoDB Atlas account (free tier available)
- Google OAuth credentials (optional, for Google Sign-In)

### Quickest Setup (5 minutes)

1. **Clone and navigate to project**
   ```bash
   cd "d:\BSSE Notes\7th Semester\Blockchain\Project"
   ```

2. **Backend (.env configuration)**
   ```bash
   cd backend
   ```
   Create `.env` file:
   ```env
   MONGO_URI=mongodb+srv://YOUR_USERNAME:YOUR_PASSWORD@cluster0.xxxxx.mongodb.net/crypto_wallet?retryWrites=true&w=majority
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   PORT=8080
   ENVIRONMENT=development
   FRONTEND_URL=http://localhost:3000
   ```

3. **Run backend**
   ```bash
   go run main.go
   ```
   Backend runs at: `http://localhost:8080`

4. **Frontend (new terminal)**
   ```bash
   cd frontend
   npm install
   npm start
   ```
   Frontend opens at: `http://localhost:3000`

---

## 📖 Detailed Setup Instructions

### Backend Setup (Go)

#### Step 1: Install Go
Make sure you have Go 1.21+ installed. Check with:
```bash
go version
```

#### Step 2: Navigate to Backend
```bash
cd backend
```

#### Step 3: Create .env File
```env
MONGO_URI=mongodb+srv://YOUR_USERNAME:YOUR_PASSWORD@cluster0.xxxxx.mongodb.net/crypto_wallet?retryWrites=true&w=majority
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000

# Optional: Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Optional: Google OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

#### Step 4: Configure MongoDB Atlas

1. Go to [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
2. Create a free account and cluster
3. Click "Connect" → "Connect your application"
4. Copy the connection string
5. Update `MONGO_URI` in `.env`:

```env
MONGO_URI=mongodb+srv://YOUR_USERNAME:YOUR_PASSWORD@cluster0.xxxxx.mongodb.net/crypto_wallet?retryWrites=true&w=majority
```

**Important**: Whitelist your IP in MongoDB Atlas:
- Go to Network Access → Add IP Address
- For local development: Add your current IP
- For production: Allow 0.0.0.0/0 (or your Render IP)

#### Step 5: Install Dependencies
```bash
go mod download
```

#### Step 6: Run the Backend
```bash
go run main.go
```

Expected output:
```
✅ Connected to MongoDB
🚀 Server listening on 0.0.0.0:8080
```

Backend API: `http://localhost:8080`

---

### Frontend Setup (React)

#### Step 1: Navigate to Frontend
```bash
cd frontend
```

#### Step 2: Install Dependencies
```bash
npm install
```

#### Step 3: Configure Environment
Create `.env` file:
```env
REACT_APP_API_URL=http://localhost:8080/api
```

#### Step 4: Start the Frontend
```bash
npm start
```

Frontend opens at: `http://localhost:3000`

---

## 🧪 Testing All Features

### Authentication Flow

#### 1. Signup
```
1. Navigate to http://localhost:3000
2. Click "Sign Up"
3. Fill in:
   - Full Name: John Doe
   - Email: test@example.com
   - CNIC: 12345-1234567-1
   - Password: password123
4. Click "Sign Up"
```

#### 2. Email Verification (OTP)
```
1. Check backend console for OTP (e.g., 123456)
2. Enter OTP in verification page
3. Click "Verify"
```

#### 3. Login
```
1. Go to Login page
2. Enter email and password
3. Click "Sign In"
4. Redirected to Dashboard
```

#### 4. Google OAuth (Bonus)
```
1. Click "Sign in with Google" button
2. Select your Google account
3. Automatically logged in and redirected to Dashboard
```

---

### Wallet Features

#### Create Wallet
```
1. Go to Dashboard
2. Click "Create Wallet"
3. Enter wallet name
4. Click "Create"
5. View wallet address and balance
```

#### Send Transaction
```
1. Go to your wallet profile
2. Click "Send Money"
3. Enter:
   - Beneficiary address: recipient-wallet-address
   - Amount: 10
4. Click "Send"
5. Transaction added to mempool
```

#### Mine Blocks
```
1. Click "Mine Block" button
2. Wait for Proof-of-Work calculation
3. Block mined with:
   - Merkle root of transactions
   - Previous block hash
   - Nonce (proof-of-work)
4. All transactions confirmed
```

#### View Blockchain
```
1. Go to Dashboard → Blocks
2. See all mined blocks
3. Verify Merkle root for each block
4. Click on block to see transaction details
```

---

### Advanced Features

#### Zakat Calculation (Automatic)
```
1. Have wallet with balance ≥ 3000
2. Go to Zakat page
3. System automatically calculates 2.5% zakat due
4. View calculation breakdown
5. Click "Pay Zakat" to send to Zakat fund
```

#### Beneficiaries Management
```
1. Go to Beneficiaries page
2. Click "Add Beneficiary"
3. Enter beneficiary details:
   - Name
   - Email
   - Wallet Address
   - Relationship
4. Save beneficiary
5. Send quick transactions to saved beneficiaries
```

#### Activity Logs
```
1. Go to Dashboard → Reports
2. View complete transaction history
3. Filter by date, amount, status
4. Download transaction reports
```

---

### Admin Features (Bonus)

#### Admin Panel Access
```
1. Login with admin account
2. See "Admin Panel" in navbar
3. Click to access admin dashboard
```

#### Admin Dashboard
```
Tabs available:
1. Dashboard - System statistics, user count, transaction volume
2. Users - List all users, toggle admin status, delete users
3. Transactions - View all transactions with search/filter
4. Blocks - Blockchain explorer showing all blocks and merkle roots
```

#### Admin Actions
```
- View system stats (total users, wallets, transactions)
- Promote/demote users to admin
- Delete users and their data
- Search and filter transactions
- View blockchain with full details
- Monitor system activity logs
```

---

## 📊 API Documentation

### Base URL
```
http://localhost:8080/api
```

### Authentication Endpoints

#### POST /auth/signup
Create new user account
```json
{
  "fullName": "John Doe",
  "email": "john@example.com",
  "cnic": "12345-1234567-1",
  "password": "password123"
}
```

#### POST /auth/verify-otp
Verify email with OTP
```json
{
  "email": "john@example.com",
  "otp": "123456"
}
```

#### POST /auth/login
Login and get JWT token
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

#### POST /auth/google (BONUS)
Google OAuth authentication
```json
{
  "token": "google-id-token"
}
```

#### GET /auth/profile
Get user profile (requires Authorization header)
```
Authorization: Bearer YOUR_JWT_TOKEN
```

---

### Wallet Endpoints

#### GET /wallet/all
Get all user wallets
```
Authorization: Bearer YOUR_JWT_TOKEN
```

#### POST /wallet/create
Create new wallet
```json
{
  "walletName": "My First Wallet"
}
```

#### GET /wallet/:walletId
Get wallet details
```
Authorization: Bearer YOUR_JWT_TOKEN
```

#### GET /wallet/:walletId/balance
Get wallet balance
```
Authorization: Bearer YOUR_JWT_TOKEN
```

---

### Transaction Endpoints

#### POST /utxo/send
Send transaction
```json
{
  "senderWalletId": "wallet-id",
  "recipientAddress": "recipient-address",
  "amount": 10
}
```

#### GET /utxo/transactions
Get all transactions
```
Authorization: Bearer YOUR_JWT_TOKEN
```

#### POST /utxo/mine
Mine block with pending transactions
```
Authorization: Bearer YOUR_JWT_TOKEN
```

#### GET /utxo/blocks
Get all blocks (with Merkle roots)
```
Authorization: Bearer YOUR_JWT_TOKEN
```

---

### Admin Endpoints (BONUS)

All admin endpoints require `AdminRequired` middleware (user must have `isAdmin: true`)

#### GET /admin/stats
Get system statistics
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

#### GET /admin/users
Get all users
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

#### PUT /admin/users/:userId/admin
Toggle user admin status
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

#### DELETE /admin/users/:userId
Delete user
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

#### GET /admin/transactions
Get all transactions
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

#### GET /admin/blocks
Get all blocks
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

---

## 🗄️ MongoDB Collections

### users
```javascript
{
  _id: ObjectId,
  fullName: "John Doe",
  email: "john@example.com",
  cnic: "12345-1234567-1",
  password: "bcrypt-hash",
  walletIds: ["wallet-id-1", "wallet-id-2"],
  publicKey: "public-key",
  privateKey: "encrypted-private-key",
  beneficiaries: [{
    name: "Jane Doe",
    email: "jane@example.com",
    walletAddress: "address",
    relationship: "Sister"
  }],
  isVerified: true,
  isAdmin: false,
  googleId: "google-oauth-id",
  authProvider: "email", // or "google"
  createdAt: ISODate,
  updatedAt: ISODate
}
```

### wallets
```javascript
{
  _id: ObjectId,
  userId: ObjectId,
  walletName: "My First Wallet",
  walletAddress: "unique-address",
  publicKey: "public-key",
  balance: 100,
  createdAt: ISODate
}
```

### transactions
```javascript
{
  _id: ObjectId,
  fromWallet: "sender-address",
  toWallet: "recipient-address",
  amount: 10,
  fee: 1,
  status: "confirmed", // pending or confirmed
  inputs: [/* UTXO inputs */],
  outputs: [/* UTXO outputs */],
  timestamp: ISODate,
  blockHeight: 5
}
```

### blocks
```javascript
{
  _id: ObjectId,
  blockHeight: 5,
  previousHash: "previous-block-hash",
  merkleRoot: "merkle-root-hash",
  timestamp: ISODate,
  nonce: 12345,
  difficulty: 4,
  transactions: [/* transaction IDs */],
  minedBy: "miner-wallet-address",
  reward: 50
}
```

---

## ✅ Features Checklist

### Core Modules
- ✅ User Authentication (signup, OTP, login)
- ✅ Wallet Management (create, view, balance)
- ✅ UTXO Model (inputs, outputs, prevent double spending)
- ✅ Transactions (send, receive, fees)
- ✅ Blockchain (mine, verify, merkle tree)
- ✅ Zakat (calculate, track, pay)
- ✅ Reports & Logs (transaction history, activity logs)

### Bonus Features
- ✅ Google OAuth (sign in with Google)
- ✅ Merkle Root Verification (block verification)
- ✅ Admin Panel (users, transactions, blocks management)

### UI/UX
- ✅ Modern light theme
- ✅ Responsive design
- ✅ Intuitive navigation
- ✅ Error handling
- ✅ Loading states
- ✅ Form validation

---

## 🔒 Security Features

- ✅ Passwords hashed with bcrypt (cost 14)
- ✅ JWT-based authentication (7-day expiry)
- ✅ OAuth 2.0 integration
- ✅ Protected API endpoints
- ✅ Email verification required
- ✅ CORS protection
- ✅ Encrypted private keys
- ✅ MongoDB Atlas encryption at rest

---

## 🚀 Production Deployment

### Deploy Backend to Render

1. Push code to GitHub
2. Create Render account and connect repository
3. Create new Web Service with settings:
   ```
   Build Command: go build -o app
   Start Command: ./app
   Environment Variables:
     - MONGO_URI: your-mongodb-connection-string
     - JWT_SECRET: your-secret-key
     - PORT: 8080
     - ENVIRONMENT: production
   ```
4. Deploy and get Render URL

### Deploy Frontend to Vercel

1. Push code to GitHub
2. Create Vercel account and connect repository
3. Set build command: `npm run build`
4. Set start command: `npm start`
5. Add environment variable:
   ```
   REACT_APP_API_URL=https://your-render-backend-url/api
   ```
6. Deploy and get Vercel URL

### MongoDB Atlas Configuration

1. Create cluster and database
2. Add Network Access (IP whitelist):
   - For Render: Allow 0.0.0.0/0 or add Render's IP range
   - For local: Add your IP
3. Create database user with username/password
4. Get connection string and update environment variables

---

## 🐛 Troubleshooting

### Backend Issues

**MongoDB Connection Failed**
```
Solution:
1. Check MongoDB Atlas connection string in .env
2. Ensure IP is whitelisted in Network Access
3. Verify username and password are correct
4. For Render: Allow 0.0.0.0/0 in Network Access
```

**Port Already in Use**
```
Solution:
1. Change PORT in .env (e.g., 8081)
2. Or kill process: lsof -ti:8080 | xargs kill
```

**CORS Errors**
```
Solution:
1. Check FRONTEND_URL in .env
2. Ensure frontend URL matches exactly
3. Verify CORS config in backend/config/cors.go
```

### Frontend Issues

**API Calls Failing**
```
Solution:
1. Ensure backend is running: http://localhost:8080/health
2. Check REACT_APP_API_URL in .env
3. Look for CORS errors in browser console
4. Verify backend is serving on correct port
```

**Styles Not Loading**
```
Solution:
1. Reinstall Tailwind: npm install -D tailwindcss
2. Check tailwind.config.js exists
3. Restart dev server: npm start
```

**Google OAuth Not Working**
```
Solution:
1. Get Google credentials from Google Cloud Console
2. Add localhost:3000 to authorized origins
3. Add GOOGLE_CLIENT_ID to frontend environment
```

### Deployment Issues

**Render Build Fails**
```
Solution:
1. Check build logs in Render dashboard
2. Ensure go.mod has all dependencies
3. Verify PORT=8080 is set in environment
4. Check render.yaml configuration
```

**MongoDB Connection Timeout on Render**
```
Solution:
1. Whitelist 0.0.0.0/0 in MongoDB Atlas Network Access
2. Use MONGODB_URI (supports connection pooling)
3. Increase timeout in connection.go to 30s
4. Check MongoDB Atlas is in same region as Render
```

---

## 📝 Environment Variables Reference

### Backend (.env)
```env
# Database
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/db
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/db

# Authentication
JWT_SECRET=your-secret-key-min-32-chars

# Server
PORT=8080
ENVIRONMENT=development

# Frontend
FRONTEND_URL=http://localhost:3000

# Email (Optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=app-password

# Google OAuth (Optional)
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
```

### Frontend (.env)
```env
REACT_APP_API_URL=http://localhost:8080/api
```

---

## 📚 Additional Resources

- [Go Documentation](https://golang.org/doc)
- [React Documentation](https://react.dev)
- [MongoDB Atlas Documentation](https://docs.atlas.mongodb.com)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [JWT Authentication](https://tools.ietf.org/html/rfc7519)

---

## 👨‍💻 Development Commands

### Backend
```bash
# Run server
go run main.go

# Build for production
go build -o app

# Run tests
go test ./...

# Format code
go fmt ./...

# Lint code
go vet ./...

# Install dependencies
go mod download
go mod tidy
```

### Frontend
```bash
# Start development server
npm start

# Build for production
npm run build

# Run tests
npm test

# Eject configuration (irreversible)
npm eject

# Clean and reinstall
rm -rf node_modules package-lock.json && npm install
```

---

## 🎯 Project Milestones

- ✅ **Phase 1**: Authentication & Setup (Module 1)
- ✅ **Phase 2**: Wallet System & Key Generation (Module 2)
- ✅ **Phase 3**: UTXO Model & Transactions (Module 3)
- ✅ **Phase 4**: Blockchain & Mining (Module 4)
- ✅ **Phase 5**: Zakat Calculation (Module 5)
- ✅ **Phase 6**: Reports & Logs (Module 6-7)
- ✅ **Phase 7**: UI/UX Enhancement
- ✅ **Phase 8**: Bonus Features (Google OAuth, Merkle Root, Admin Panel)
- ✅ **Phase 9**: Deployment Configuration

---

## 📄 License

This project is for educational purposes as part of the BSSE 7th Semester Blockchain course.

---

## 📞 Support & Feedback

For issues or questions:
1. Check the troubleshooting section above
2. Review API documentation
3. Check backend console logs
4. Verify .env configuration
5. Ensure MongoDB Atlas is configured correctly

---

**Last Updated**: December 2024
**Status**: ✅ All Features Complete & Ready for Production
