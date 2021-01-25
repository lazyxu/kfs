import ReconnectingWebSocket from 'reconnecting-websocket';
import { v4 as uuidv4 } from 'uuid';
import { EventEmitter } from 'events';
import prettyBytes from 'pretty-bytes';

import { getConfig } from './config';

const rws = new ReconnectingWebSocket(getConfig().wsHost);
const emitter = new EventEmitter();

rws.addEventListener('open', () => {
  rws.send(JSON.stringify({ method: 'hello, world' }));
  rws.send(JSON.stringify({ method: 'echo' }));
});
rws.addEventListener('message', ({ data }) => {
  try {
    data = JSON.parse(data);
  } catch (e) {
    e; // ignore error
    return;
  }
  console.log('message', data.id, data.result);
  if (data.id) {
    emitter.emit(data.id, data.result);
  }
});

export function invoke(method, params, cb) {
  const id = uuidv4();
  emitter.addListener(id, cb);
  rws.send(JSON.stringify({ method, id, params }));
}

// let fileCnt = 0;
// let fileSize = 0;
invoke('backup', {}, result => {
  // const files = result.files.filter(f => f.type === 'file');
  // fileCnt += files.length;
  // const size = files.reduce((prev, curr) => prev += curr.size, 0);
  // fileSize += size;
  console.log('backup result', result);
});

export default rws;
