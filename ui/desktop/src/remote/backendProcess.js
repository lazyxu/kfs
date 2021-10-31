const remote = require('@electron/remote');

function updateStatus(port) {
  return new Promise((resolve, reject) => {
    window.goBackendInstance.get('/api/clientID').then(function (response) {
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
