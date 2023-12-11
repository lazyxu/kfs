// import ReconnectingWebSocket from 'reconnecting-websocket';
// import WebSocket from 'ws';

import { getConfig } from './config';

let protobuf = require("protobufjs");

let gRoot;

function getRoot() {
  return new Promise((resolve, reject) => {
    if (gRoot) {
      resolve(gRoot);
      return;
    }
    protobuf.load("./fs.proto", (err, root) => {
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
  let message = type.decode(new Uint8Array(buffer));
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

const CommandPing = 0;
const CommandReset = 1;
const CommandList = 2;
const CommandOpen = 3;

export function list(sysConfig, driverName, path, onTotal, onDirItem) {
  console.log('list', driverName, path, path.join('/'))
  return new Promise((resolve, reject) => {
    const ws = new WebSocket(`ws://${sysConfig.socketServer}/ws`);
    ws.addEventListener('open', async () => {
      try {
        ws.send(new Uint8Array([CommandList]));
        let reqData = await encode("PathReq", { driverName, path: path.join('/') });
        // console.log('reqData', reqData);
        ws.send(new Int32Array([reqData.length, 0]));
        ws.send(reqData);

        let receiver = new WebSocketReceiver(ws);
        let data = await receiver.recv();
        let code = new Int8Array(data)[0];
        // console.log('code', data, code);
        data = await receiver.recv();
        let total = new Int32Array(data)[0];
        // console.log('total', data, total);
        onTotal && onTotal(total);
        for (let i = 0; i < total; i++) {
          data = await receiver.recv();
          // console.log('length', new Int32Array(data)[0]);
          data = await receiver.recv();
          let resp = await decode("DirItem", data);
          onDirItem && onDirItem(resp, i);
          // console.log('resp', resp);
        }
        data = await receiver.recv();
        code = new Int8Array(data)[0];
        // console.log('exit code', code);
        resolve(code);
      } catch (e) {
        reject(e);
      } finally {
        ws.close();
      }
    });
  });
}

export function open(sysConfig, driverName, path, onFile, onTotal, onDirItem) {
  console.log('open', driverName, path, path.join('/'))
  return new Promise((resolve, reject) => {
    const ws = new WebSocket(`ws://${sysConfig.socketServer}/ws`);
    ws.addEventListener('open', async () => {
      try {
        ws.send(new Uint8Array([CommandOpen]));
        let reqData = await encode("PathReq", { driverName, path: path.join('/') });
        // console.log('reqData', reqData);
        ws.send(new Int32Array([reqData.length, 0]));
        ws.send(reqData);

        let receiver = new WebSocketReceiver(ws);
        let data = await receiver.recv();
        let code = new Int8Array(data)[0];
        // console.log('code', data, code);
        data = await receiver.recv();
        let mode = new Int32Array(data)[0];
        console.log('mode', data, mode);
        if (mode >= 0) {
          data = await receiver.recv();
          let size = new Int32Array(data)[0];
          console.log('size', data, size);
          if (size !== 0) {
            let buf = new Uint8Array(size);
            let offset = 0;
            while (offset < size) {
              data = await receiver.recv();
              console.log('file', buf, data, offset, offset + data.byteLength, size);
              buf.set(new Uint8Array(data), offset);
              console.log('buf', buf);
              offset += data.byteLength;
            }
            onFile && onFile(buf);
          }
          resolve(0);
          return;
        }
        data = await receiver.recv();
        let total = new Int32Array(data)[0];
        // console.log('total', data, total);
        onTotal && onTotal(total);
        for (let i = 0; i < total; i++) {
          data = await receiver.recv();
          // console.log('length', new Int32Array(data)[0]);
          data = await receiver.recv();
          let resp = await decode("DirItem", data);
          onDirItem && onDirItem(resp, i);
          // console.log('resp', resp);
        }
        data = await receiver.recv();
        code = new Int8Array(data)[0];
        // console.log('exit code', code);
        resolve(1);
      } catch (e) {
        reject(e);
      } finally {
        ws.close();
      }
    });
  });
}
