const WebSocket = require('ws')
const socket = new WebSocket('ws://localhost:8080/ws')
socket.addEventListener('message', event => {
  const data = JSON.parse(event.data)
  console.log(`Received message from ${data.exchange}: ${data.symbol} = ${data.price}`)
})