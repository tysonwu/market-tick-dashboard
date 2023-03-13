import { Redis } from 'ioredis';
import express from 'express';
import { createServer } from 'http';
import { WebSocketServer } from 'ws';
// import { Server } from 'socket.io';
// import { createAdapter } from '@socket.io/redis-adapter';

// const httpServer = createServer();
// const io = new Server(httpServer, {
//     cors: { origin: '*' }
// }
// );

// make ws app
const app = express();
const server = createServer(app);
const wss = new WebSocketServer({ server });

// init redis
const dbConfig = {
    host: 'localhost',
    port: 6379,
    db: 0,
};

// const keyMatch = '(bidAskTicks|ticks):.*:.*';
const client = new Redis(dbConfig);
// const subClient = pubClient.duplicate();

const keyMatch = ['*']; // all ticks and bidAskTicks
client.psubscribe(keyMatch, (err, result) => {
    if (err) {
        console.log(`subscribe error: ${err}`);
    } else {
        console.log(`subscribe result: ${result}`);
    }
});

client.on("pmessage", (_, channel, message) => {
    // const jsonMsg = JSON.parse(message);
    wss.clients.forEach(client => {
        // client.send(`${channel} | ${jsonMsg['Exchange']}, ${jsonMsg['StandardSymbol']}, ${jsonMsg['Time']}`);
        var msg = { channel: channel, message: message };
        client.send(JSON.stringify(msg));
    });
});

// wss.on("connection", (ws) => {
//     console.log('Client connected');
//     ws.on('message',)
// })

server.listen(process.env.PORT || 8999, () => {
    console.log(`Server started on port ${server.address().port} :)`);
});