# arbitrage

Crypto arbitrage research dashboard. Displays real-time tick price across two exchanges with a configurable list of tokens.

## Backend

- Run Go backend
```bash
cd backend
go run main.go
```
- Need to have Redis setup:
```bash
# from project root
docker-compose up
```

## Frontend

- Run Next.js dashboard
```bash
cd app
npm install
npm run dev
```

## Specification

- Symbol style example: BTC-USDT (original in Binance: BTCUSDT, Kucoin: BTC-USDT)