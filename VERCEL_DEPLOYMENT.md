# Frontend Deployment to Vercel

Your backend is already running at: https://cryptowallet-rsf1.onrender.com ✅

Now deploy the frontend to Vercel:

## Step 1: Install Vercel CLI (Optional)
```bash
npm install -g vercel
```

## Step 2: Deploy via Vercel Dashboard (Recommended)

1. Go to https://vercel.com
2. Sign in with GitHub
3. Click "Add New" → "Project"
4. Import your repository: `saif55045/CryptoWallet`
5. Configure project:
   - **Framework Preset**: Create React App
   - **Root Directory**: `frontend`
   - **Build Command**: `npm run build`
   - **Output Directory**: `build`
   - **Install Command**: `npm install`

6. Add Environment Variables:
   ```
   REACT_APP_API_URL=https://cryptowallet-rsf1.onrender.com/api
   REACT_APP_GOOGLE_CLIENT_ID=862938713766-otpi1uedi44bfm6sh1vht7fm0l66lt7o.apps.googleusercontent.com
   ```

7. Click "Deploy"

## Step 3: Update Google OAuth Authorized Origins

After deployment, you'll get a Vercel URL like: `https://crypto-wallet.vercel.app`

Add this to Google Cloud Console:
1. Go to https://console.cloud.google.com/apis/credentials
2. Click on your OAuth 2.0 Client ID
3. Add to "Authorized JavaScript origins":
   - `https://your-vercel-url.vercel.app`
4. Add to "Authorized redirect URIs":
   - `https://your-vercel-url.vercel.app`
5. Save

## Step 4: Update Backend CORS

Update your backend `.env` on Render:
```
FRONTEND_URL=https://your-vercel-url.vercel.app
```

Then restart the backend service on Render.

---

## Alternative: Deploy via CLI

```bash
cd frontend
vercel
```

Follow the prompts and add environment variables when asked.

---

## Troubleshooting

**Issue**: "Module not found" errors
**Solution**: Ensure `package.json` is in the `frontend` directory

**Issue**: API calls fail
**Solution**: Check REACT_APP_API_URL points to your Render backend

**Issue**: Blank page after deployment
**Solution**: Check browser console for errors, verify build completed successfully

---

## Summary

- ✅ Backend on Render: https://cryptowallet-rsf1.onrender.com
- 🚀 Frontend on Vercel: [Your URL after deployment]
- 📦 Root Directory: `frontend`
- 🔧 Environment: `REACT_APP_API_URL` pointing to Render backend
