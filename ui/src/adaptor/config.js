import { targetBrowser } from './target';

let getConfig = null;
let setConfig = null;
const defaultBrowserConfig = {
  host: 'http://127.0.0.1:9091',
  wsHost: 'ws://127.0.0.1:1323/ws',
  filter: [
  ],
};
const defaultElectronConfig = defaultBrowserConfig;
if (targetBrowser) {
  // browser
  getConfig = function () {
    const item = localStorage.getItem('kfs-config');
    return item ? JSON.parse(item) : defaultBrowserConfig;
  };

  setConfig = function (str) {
    localStorage.setItem('kfs-config', str);
  };
} else {
  const fs = window.require('fs');
  const path = window.require('path');
  console.log(process.cwd(), process.env);
  const { app } = window.require('electron').remote;
  const appPath = app.getPath('userData');
  const configPath = path.join(appPath, 'kfs-config.json');
  console.log('Your App Path: ', app.getPath('userData'), app.getPath('appData'));

  getConfig = function () {
    try {
      const config = fs.readFileSync(configPath).toString();
      console.log('getConfig', config);
      return JSON.parse(config);
    } catch (e) {
      return defaultElectronConfig;
    }
  };

  setConfig = function (str) {
    fs.writeFileSync(configPath, str, { flag: 'w+' });
  };
}

export { getConfig, setConfig };
