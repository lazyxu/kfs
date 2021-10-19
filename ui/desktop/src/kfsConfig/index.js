import { isElectron } from 'common/utils/target';
import localStorageConfig from 'common/kfsConfig/localStorage';

let get = null;
let set = null;
let remove = null;

if (isElectron) {
  const fs = window.require('fs');
  const path = window.require('path');
  console.log(process.cwd(), process.env);
  const { app } = window.require('@electron/remote');
  const appPath = app.getPath('userData');
  const configPath = path.join(appPath, 'kfs-config.json');
  console.log('Your App Path: ', app.getPath('userData'), app.getPath('appData'));

  get = function () {
    try {
      const config = fs.readFileSync(configPath).toString();
      console.log('getExternalConfig', config);
      const obj = JSON.parse(config);
      return obj;
    } catch (e) {
      return undefined;
    }
  };

  set = function (json) {
    fs.writeFileSync(configPath, JSON.stringify(json, undefined, 2), { flag: 'w+' });
  };

  remove = function () {
    fs.unlinkSync(configPath);
  };
} else {
  ({ get, set, remove } = localStorageConfig);
}

export default { get, set, remove };
