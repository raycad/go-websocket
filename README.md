A realtime distributed messaging platform<br />
https://github.com/nsqio/nsq<br />

**1. Install packages**<br />
```
    $ go get github.com/gorilla/websocket
    $ go get github.com/shirou/gopsutil/mem
```

**2. Run Go Websocket server**<br />
```
    $ cd go-ws-server
    $ go run 
```

**3. Start NodeJS Websocket client**<br />
```
    $ cd node-ws-client
    $ npm install
    $ node main.js        
```    

**4. Get server statistics information**<br />
```
http://127.0.0.1:2706/stats
``` 