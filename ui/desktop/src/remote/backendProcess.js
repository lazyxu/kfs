import axios from 'axios';
import https from 'https';

const remote = require('@electron/remote');

const httpsAgent = new https.Agent({
  ca: axios('./extraResources/rootCA.pem'),
  cert: axios('./extraResources/localhost.pem'),
  key: axios('./extraResources/localhost-key.pem'),
});

function updateStatus(port) {
  return new Promise((resolve, reject) => {
    const instance = axios.create({
      baseURL: `https://localhost:${port}`,
      httpsAgent,
    });
    instance.get('/api/clientID').then(function (response) {
      if (response.status === 200 && response.data === window.clientID) {
        resolve();
      } else {
        reject();
      }
    }).catch(reject);
  });
}

let statusInterval;

export function backendProcess(port) {
  return new Promise((resolve, reject) => {
    try {
      remote.require('./backendProcess')(port);
    } catch (e) {
      console.error(e);
      reject(e);
      return;
    }
    if (statusInterval) {
      clearInterval(statusInterval);
    }
    let i = 0;
    statusInterval = setInterval(() => {
      updateStatus(port).then(() => {
        clearInterval(statusInterval);
        resolve();
      }).catch(() => {
        i++;
        if (i === 10) {
          clearInterval(statusInterval);
        }
        reject();
      });
    }, 500);
  });
}
