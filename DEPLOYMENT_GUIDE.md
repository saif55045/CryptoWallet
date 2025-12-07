# Deployment Guide - Backend & Frontend

## 🚀 Backend Deployment (Render)

### Option 1: Manual Setup (RECOMMENDED)

1. **Go to Render Dashboard**: https://dashboard.render.com
2. Click **"New +"** → **"Web Service"**
3. Connect your GitHub repository: `saif55045/CryptoWallet`
4. Configure:
   ```
   Name: crypto-wallet-backend
   Region: Oregon (US West)
   Branch: main
   Root Directory: backend
   Runtime: Go
   Build Command: go build -o app
   Start Command: ./app
   ```

5. **Add Environment Variables**:
   ```
   PORT=8080
   GIN_MODE=release
   ENVIRONMENT=production
   MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/crypto_wallet?retryWrites=true&w=majority
   JWT_SECRET=your-jwt-secret-min-32-chars
   FRONTEND_URL=https://your-frontend-url.vercel.app
   ```

6. Click **"Create Web Service"**
7. Wait for deployment (5-10 minutes)
8. Copy your backend URL: `https://your-app.onrender.com`

---

## 🎨 Frontend Deployment (Vercel)

### Step-by-Step Guide

1. **Go to Vercel**: https://vercel.com
2. Click **"Add New"** → **"Project"**
3. **Import Git Repository**: Select `saif55045/CryptoWallet`
4. Configure Project:
   ```
   Framework Preset: Create React App
   Root Directory: frontend
   Build Command: npm run build
   Output Directory: build
   Install Command: npm install
   ```

5. **Add Environment Variables**:
   ```
   REACT_APP_API_URL=https://your-backend.onrender.com/api
   REACT_APP_GOOGLE_CLIENT_ID=862938713766-otpi1uedi44bfm6sh1vht7fm0l66lt7o.apps.googleusercontent.com
   ```

6. Click **"Deploy"**
7. Wait for deployment (2-3 minutes)
8. Your app will be live at: `https://your-app.vercel.app`

---

## 🔄 Update Backend with Frontend URL

After frontend is deployed:

1. Go to Render Dashboard → Your Backend Service
2. Go to **Environment** tab
3. Update `FRONTEND_URL`:
   ```
   FRONTEND_URL=https://your-app.vercel.app
   ```
4. Click **"Save Changes"** (this will redeploy)

---

## ✅ Verify Deployment

### Test Backend
```bash
curl https://your-backend.onrender.com/health
```
Should return: `{"status":"ok"}`

### Test Frontend
1. Open: `https://your-app.vercel.app`
2. Try signup/login
3. Create wallet
4. Send transaction
5. Mine block

---

## 🐛 Common Issues

### Issue 1: CORS Errors
**Fix**: Make sure `FRONTEND_URL` in backend matches your Vercel URL exactly (no trailing slash)

### Issue 2: MongoDB Connection Timeout
**Fix**: In MongoDB Atlas → Network Access → Add IP: `0.0.0.0/0`

### Issue 3: Backend Cold Start (Render Free Tier)
**Note**: Free tier spins down after 15 min inactivity. First request may take 30-60 seconds.

### Issue 4: Frontend 404 on Refresh
**Fix**: `vercel.json` already configured with rewrites to handle this

---

## 📝 Current Configuration

✅ Backend: Go + Gin + MongoDB
✅ Frontend: React + Tailwind CSS
✅ Database: MongoDB Atlas
✅ Deployment: Render (backend) + Vercel (frontend)

---

## 🔗 Quick Links

- **Render Dashboard**: https://dashboard.render.com
- **Vercel Dashboard**: https://vercel.com/dashboard
- **MongoDB Atlas**: https://cloud.mongodb.com
- **GitHub Repo**: https://github.com/saif55045/CryptoWallet

---

**Note**: If you want to deploy frontend on Render instead of Vercel, you'll need to create a SEPARATE Web Service for the frontend with different settings (Node.js runtime, npm build commands, etc.)
