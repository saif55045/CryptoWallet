# 🔐 Crypto Wallet System - Complete Setup Guide

## Module 1: User Authentication & Project Setup ✅

This guide will help you set up and test Module 1 of the Decentralized Cryptocurrency Wallet System.

---

## 📁 Project Structure

```
Project/
├── backend/           # Go backend server
│   ├── config/
│   ├── controllers/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   ├── utils/
│   ├── main.go
│   ├── go.mod
│   └── .env
├── frontend/          # React frontend
│   ├── src/
│   │   ├── pages/
│   │   ├── context/
│   │   ├── services/
│   │   └── App.js
│   └── package.json
└── README_SETUP.md
```

---

## 🚀 Backend Setup (Go)

### Step 1: Install Go
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
