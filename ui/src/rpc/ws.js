// import ReconnectingWebSocket from 'reconnecting-websocket';
// import WebSocket from 'ws';

import { getConfig } from './config';
import { PathReq, DirItem } from '../pb/fs_pb';

async function wait(fn) {
  return await new Promise((resolve, reject) => {
    setInterval(() => {
      if (fn()) {
        resolve();
      }
    }, 100);
  });
}

export function list() {
  let p = {};
  const rws = new WebSocket(getConfig().wsHost);
  rws.addEventListener('open', async () => {
    rws.send(new Uint8Array([1]));
    let req = new PathReq();
    console.log(req);
    req.setBranchname("master");
    req.setPath("/");
    let reqData = req.serializeBinary();
    console.log(reqData);
    rws.send(new Int32Array([reqData.length, 0]));
    rws.send(reqData);
    await wait(() => { return p.p; });
    let data = await p.p;
    p = p.p;
    let code = new Int8Array(data)[0];
    console.log('code', code);
    data = await p.p;
    p = p.p;
    let total = new Int32Array(data)[0];
    console.log('total', total);
    for (let i = 0; i < total; i++) {
      data = await p.p;
      p = p.p;
      console.log('length', new Int32Array(data)[0]);
      data = await p.p;
      p = p.p;
      let resp = DirItem.deserializeBinary(data);
      console.log('resp', resp.toObject());
    }
  });
  let lastP = p;
  rws.addEventListener('message', ({ data }) => {
    lastP.p = (async () => {
      return await data.arrayBuffer();
    })()
    lastP = lastP.p;
  });
}
