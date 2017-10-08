// var sleep = require('sleep');
const WebSocket = require('ws');
const ws = new WebSocket('ws://127.0.0.1:2706/aep');

ws.on('open', function open() {
    console.log('connected');
    while (true) {
        ws.send("{\"username\": \"raycad\", \"email\": \"seedotech@gmail.com\", \"result\": \"\"}");
        // sleep.sleep(1);
    }
});

ws.on('close', function close(code, reason) {
    console.log('disconnected code = %d, reason = %d', code, reason);
});

ws.on('message', function incoming(data) {
    console.log(data);
    setTimeout(function timeout() {
        ws.send("{\"username\": \"raycad\", \"email\": \"seedotech@gmail.com\", \"result\": \"\"}");
        // let msgBatch = 10000;
        // for (let i = 0; i < msgBatch; ++i)
        //     ws.send("{\"username\": \"raycad\", \"email\": \"seedotech@gmail.com\", \"result\": \"\"}");
    }, 2000);
});