// import ReconnectingWebSocket from 'reconnecting-websocket';
// import WebSocket from 'ws';

import { getConfig } from './config';

var protobuf = require("protobufjs");

let gRoot;

function getRoot() {
  return new Promise((resolve, reject) => {
    if (gRoot) {
      resolve(gRoot);
      return;
    }
    protobuf.load("./fs.proto", function (err, root) {
      if (err)
        throw err;
      gRoot = root;
      resolve(gRoot);
    });
  });
}

async function encode(path, payload) {
  let pbRoot = await getRoot();
  let type = pbRoot.lookupType(path);
  let message = type.create(payload);
  return type.encode(message).finish();
}

async function decode(path, buffer) {
  let pbRoot = await getRoot();
  let type = pbRoot.lookupType(path);
  var message = type.decode(new Uint8Array(buffer));
  return type.toObject(message);
}

async function blobToArrayBuffer(blob) {
  if ('arrayBuffer' in blob) return await blob.arrayBuffer();

  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result);
    reader.onerror = () => reject;
    reader.readAsArrayBuffer(blob);
  });
}

class WebSocketReceiver {
  constructor(ws) {
    this.ws = ws;
    this.nextPromise = this.getNextPromise();
  }
  getNextPromise() {
    let that = this;
    let cur = new Promise((resolve, reject) => {
      this.ws.addEventListener('message', ({ data }) => {
        blobToArrayBuffer(data).then(bytes => {
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

const CommandList = 2;

export function list(onTotal, onDirItem) {
  return new Promise((resolve, reject) => {
    const ws = new WebSocket(getConfig().wsHost);
    ws.addEventListener('open', async () => {
      try {
        ws.send(new Uint8Array([CommandList]));
        let reqData = await encode("PathReq", { branchName: "master", path: '/' });
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
        onTotal && onTotal(total);
        for (let i = 0; i < total; i++) {
          data = await receiver.recv();
          console.log('length', new Int32Array(data)[0]);
          data = await receiver.recv();
          let resp = await decode("DirItem", data);
          onTotal && onTotal(onDirItem);
          console.log('resp', resp);
        }
        data = await receiver.recv();
        code = new Int8Array(data)[0];
        console.log('exit code', code);
        resolve(code);
      } catch (e) {
        reject(e);
      } finally {
        ws.close();
      }
    });
  });
}
