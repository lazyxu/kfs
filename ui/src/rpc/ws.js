// import ReconnectingWebSocket from 'reconnecting-websocket';
// import WebSocket from 'ws';

import { getConfig } from './config';
import { PathReq, DirItem } from '../pb/fs_pb';

export function list() {
  const rws = new WebSocket(getConfig().wsHost);
  rws.addEventListener('open', () => {
    rws.send(new Uint8Array([1]));
    let req = new PathReq();
    console.log(req);
    req.setBranchname("master");
    req.setPath("/");
    let reqData = req.serializeBinary();
    console.log(reqData);
    rws.send(new Int32Array([reqData.length, 0]));
    rws.send(reqData);
  });
  let state = 0;
  let funcs = [];
  let i = -1;
  let total;
  rws.addEventListener('message', async ({ data }) => {
    i++;
    let ii = i;
    if (i != 0) {
      await new Promise((resolve, reject)=> {
        setInterval(() => {
          if (funcs[ii - 1]) {
            resolve();
          }
        }, 100);
      })
      // console.log('wait', ii, funcs[ii - 1].i, funcs.length)
      await funcs[ii - 1];
      // console.log('wait done', ii, funcs[ii - 1].i, funcs.length)
    }
    let p = async () => {
      // console.log('state', ii, state);
      switch (state) {
        case 0:
          let code = new Int8Array(await data.arrayBuffer())[0];
          console.log('code', code);
          if (code != 0) {
            state = -1;
            return;
          }
          break;
        case 1:
          total = new Int32Array(await data.arrayBuffer())[0];
          console.log('total', total);
          break;
        case 2:
          // console.log('length', new Int32Array(await data.arrayBuffer())[0]);
          break;
        case 3:
          let resp = DirItem.deserializeBinary(await data.arrayBuffer());
          console.log('resp', resp.toObject());
          total--;
          if (total != 0) {
            state -= 2;
          }
          break;
        default:
          break;
      }
      state++;
      // console.log('state done', ii, state);
    }
    p.i=ii;
    funcs.push(p());
  });
}
