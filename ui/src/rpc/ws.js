// import ReconnectingWebSocket from 'reconnecting-websocket';
// import WebSocket from 'ws';

import { getConfig } from './config';

var protobuf = require("protobufjs");

let pbRoot;

function encode(path, payload) {
  let type = pbRoot.lookupType(path);
  let message = type.create(payload);
  return type.encode(message).finish();
}

function decode(path, buffer) {
  let type = pbRoot.lookupType(path);
  var message = type.decode(new Uint8Array(buffer));
  return type.toObject(message);
}

protobuf.load("./fs.proto", function(err, root) {
  if (err)
      throw err;
  pbRoot = root;
  console.log(root);
});

class WebSocketReceiver {
  constructor(ws) {
    this.ws = ws;
    this.nextPromise = this.getNextPromise();
  }
  getNextPromise() {
    let that = this;
    let cur = new Promise((resolve, reject)=> {
      this.ws.addEventListener('message', ({ data }) => {
        data.arrayBuffer().then(bytes => {
          resolve(bytes);
        });
        cur.next = that.getNextPromise();
      }, { once: true });
    });
    return cur;
  }
  async recv() {
    let data = await this.nextPromise;
    this.nextPromise = this.nextPromise.next;
    return data;
  }
}

export function list() {
  const ws = new WebSocket(getConfig().wsHost);
  ws.addEventListener('open', async () => {
    ws.send(new Uint8Array([1]));
    let reqData = encode("PathReq", {branchName: "master", path: '/'});
    console.log('reqData', reqData);
    ws.send(new Int32Array([reqData.length, 0]));
    ws.send(reqData);

    let receiver = new WebSocketReceiver(ws);
    let data = await receiver.recv();
    let code = new Int8Array(data)[0];
    console.log('code', data, code);
    data = await receiver.recv();
    let total = new Int32Array(data)[0];
    console.log('total', data, total);
    for (let i = 0; i < total; i++) {
      data = await receiver.recv();
      console.log('length', new Int32Array(data)[0]);
      data = await receiver.recv();
      let resp = decode("DirItem", data);
      console.log('resp', resp);
    }
    data = await receiver.recv();
    code = new Int8Array(data)[0];
    console.log('exit code', code);
  });
}
