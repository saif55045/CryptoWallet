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

### Step 2: Navigate to Backend
```bash
cd backend
```

### Step 3: Create .env File
Copy `.env.example` to `.env`:
```bash
copy .env.example .env
```

### Step 4: Configure MongoDB Atlas

1. Go to [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
2. Create a free account and cluster
3. Click "Connect" → "Connect your application"
4. Copy the connection string
5. Update `.env` file:

```env
MONGO_URI=mongodb+srv://YOUR_USERNAME:YOUR_PASSWORD@cluster0.xxxxx.mongodb.net/crypto_wallet?retryWrites=true&w=majority

JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

PORT=8080
ENVIRONMENT=development

# Email Configuration (Optional for testing)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=your-email@gmail.com
SMTP_PASSWORD=your-app-password

FRONTEND_URL=http://localhost:3000
```

**Note**: For development, OTP will be printed in the console if email is not configured.

### Step 5: Install Dependencies
```bash
go mod download
```

### Step 6: Run the Backend
```bash
go run main.go
```

You should see:
```
✅ Connected to MongoDB Atlas
🚀 Server starting on port 8080
```

Backend is now running at: `http://localhost:8080`

---

## 🎨 Frontend Setup (React)

### Step 1: Navigate to Frontend
```bash
cd frontend
```

### Step 2: Install Dependencies
```bash
npm install
```

### Step 3: Configure Environment
The `.env` file is already created with:
```env
REACT_APP_API_URL=http://localhost:8080/api
```

### Step 4: Start the Frontend
```bash
npm start
```

The app will automatically open at: `http://localhost:3000`

---

## 🧪 Testing Module 1

### Test Case 1: User Signup ✅

1. Open `http://localhost:3000`
2. Click "Sign up"
3. Fill in the form:
   - Full Name: `John Doe`
   - Email: `test@example.com`
   - CNIC: `12345-1234567-1`
   - Password: `password123`
   - Confirm Password: `password123`
4. Click "Sign up"
5. You should see: "User created successfully. Please verify your email..."

### Test Case 2: Email Verification (OTP) ✅

1. Check your backend console for the OTP (6-digit code)
   ```
   📧 OTP for test@example.com: 123456
   ```
2. Enter the OTP in the verification page
3. Click "Verify Email"
4. You should be redirected to the dashboard

### Test Case 3: Login ✅

1. Go to login page
2. Enter credentials:
   - Email: `test@example.com`
   - Password: `password123`
3. Click "Sign in"
4. You should be redirected to the dashboard

### Test Case 4: Protected Routes ✅

1. Try accessing `/dashboard` without logging in
2. You should be redirected to `/login`
3. After login, try accessing `/login` or `/signup`
4. You should be redirected to `/dashboard`

### Test Case 5: Logout ✅

1. Click the "Logout" button in the dashboard
2. You should be redirected to login page
3. Token and user data should be cleared

---

## 📊 API Endpoints (Module 1)

### 1. Health Check
```http
GET http://localhost:8080/health
```

### 2. Signup
```http
POST http://localhost:8080/api/auth/signup
Content-Type: application/json

{
  "fullName": "John Doe",
  "email": "test@example.com",
  "cnic": "12345-1234567-1",
  "password": "password123"
}
```

### 3. Verify OTP
```http
POST http://localhost:8080/api/auth/verify-otp
Content-Type: application/json

{
  "email": "test@example.com",
  "otp": "123456"
}
```

### 4. Resend OTP
```http
POST http://localhost:8080/api/auth/resend-otp
Content-Type: application/json

{
  "email": "test@example.com"
}
```

### 5. Login
```http
POST http://localhost:8080/api/auth/login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

### 6. Get Profile (Protected)
```http
GET http://localhost:8080/api/auth/profile
Authorization: Bearer YOUR_JWT_TOKEN
```

---

## 🗄️ MongoDB Collections Created

### users
```javascript
{
  _id: ObjectId,
  fullName: "John Doe",
  email: "test@example.com",
  cnic: "12345-1234567-1",
  password: "$2a$14$...", // bcrypt hashed
  walletId: null,        // Will be generated in Module 2
  publicKey: null,       // Will be generated in Module 2
  privateKey: null,      // Will be generated in Module 2
  beneficiaries: [],
  isVerified: true,
  createdAt: ISODate("2025-12-07T..."),
  updatedAt: ISODate("2025-12-07T...")
}
```

---

## ✅ Module 1 Features Completed

- ✅ User Signup with validation
- ✅ Email OTP verification (10-minute expiry)
- ✅ OTP resend functionality
- ✅ Secure password hashing (bcrypt)
- ✅ JWT authentication
- ✅ Protected routes
- ✅ User login/logout
- ✅ User profile retrieval
- ✅ MongoDB Atlas integration
- ✅ CORS configuration
- ✅ Responsive UI (Tailwind CSS)
- ✅ Error handling
- ✅ Loading states

---

## 🐛 Troubleshooting

### Backend Issues

**Problem**: Cannot connect to MongoDB
```
Solution: Check your MongoDB Atlas connection string in .env
- Ensure your IP is whitelisted in MongoDB Atlas
- Verify username and password are correct
```

**Problem**: Port 8080 already in use
```
Solution: Change PORT in .env file to another port (e.g., 8081)
```

### Frontend Issues

**Problem**: API calls failing
```
Solution: 
- Ensure backend is running on port 8080
- Check REACT_APP_API_URL in .env
- Check browser console for CORS errors
```

**Problem**: Tailwind styles not working
```
Solution:
- Make sure tailwindcss is installed: npm install -D tailwindcss
- Check that tailwind.config.js exists
- Restart the dev server: npm start
```

---

## 🎯 Next Steps

After testing Module 1 successfully, we will proceed to:

### Module 2: Wallet System & Key Generation
- RSA/ECDSA key pair generation
- Wallet ID generation from public key
- Encrypted private key storage
- Wallet profile page

---

## 📝 Notes

- OTPs are valid for 10 minutes
- Passwords must be at least 6 characters
- JWT tokens expire in 7 days
- CNIC and Email must be unique

---

## 🔒 Security Features

- Passwords hashed with bcrypt (cost: 14)
- JWT-based authentication
- Protected API routes
- Email verification required
- CORS protection
- Environment variable configuration

---

## ✉️ Contact & Support

If you encounter any issues during setup, check:
1. Go version (1.21+)
2. Node version (14+)
3. MongoDB Atlas connection
4. Environment variables in .env files

---

**🎉 Congratulations! Module 1 is complete. Please test all features and confirm before moving to Module 2.**
