const { dialog } = require('@electron/remote');

export function chooseDir() {
  return new Promise((resolve, reject) => {
    dialog.showOpenDialog({
      properties: ['openDirectory', 'showHiddenFiles'],
    }).then(result => {
      if (result.canceled) {
        reject();
        return;
      }
      resolve(result.filePaths[0]);
    }).catch(err => {
      reject(err);
    });
  });
}
