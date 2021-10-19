import { targetBrowser } from '../common/utils/target';

let getExternalConfig = null;
let setExternalConfig = null;
let resetExternalConfig = null;

const defaultBrowserConfig = {
  theme: 'dark',
  username: '17161951517',
  refreshToken: '96246b97eb994fcaa4e8abb553d502bb',
  downloadPath: '',
};

const defaultElectronConfig = defaultBrowserConfig;
if (targetBrowser) {
  // browser
  getExternalConfig = function () {
    const item = localStorage.getItem('kfs-config');
    const obj = item ? JSON.parse(item) : defaultBrowserConfig;
    console.log(obj);
    return obj;
  };

  setExternalConfig = function (str) {
    localStorage.setItem('kfs-config', str);
  };
  resetExternalConfig = function (cb) {
    localStorage.removeItem('kfs-config');
    cb(getExternalConfig());
  };
} else {
  const fs = window.require('fs');
  const path = window.require('path');
  console.log(process.cwd(), process.env);
  const { app } = window.require('electron').remote;
  const appPath = app.getPath('userData');
  const configPath = path.join(appPath, 'kfs-config.json');
  console.log('Your App Path: ', app.getPath('userData'), app.getPath('appData'));

  getExternalConfig = function () {
    try {
      const config = fs.readFileSync(configPath).toString();
      console.log('getExternalConfig', config);
      const obj = JSON.parse(config);
      // const newObj = bindData(obj);
      return obj;
    } catch (e) {
      return defaultElectronConfig;
    }
  };

  resetExternalConfig = function (cb) {
    fs.unlinkSync(configPath);
    cb(getExternalConfig());
  };
}

export { getExternalConfig, setExternalConfig, resetExternalConfig };
