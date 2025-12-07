# Quick Start Scripts

## Backend Start
```bash
cd backend
go run main.go
```

## Frontend Start  
```bash
cd frontend
npm start
```

## MongoDB Setup Checklist
- [ ] Created MongoDB Atlas account
- [ ] Created a free cluster
- [ ] Whitelisted IP address (0.0.0.0/0 for testing)
- [ ] Created database user
- [ ] Copied connection string to backend/.env
- [ ] Replaced <password> with actual password

## Testing Flow
1. Start Backend (terminal 1)
2. Start Frontend (terminal 2)
3. Open http://localhost:3000
4. Sign up with test credentials
5. Check backend console for OTP
6. Verify OTP
7. Login and access dashboard

## Quick Test Credentials
- Email: test@example.com
- Password: password123
- CNIC: 12345-1234567-1
- Full Name: John Doe
