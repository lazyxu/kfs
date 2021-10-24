import { isElectron } from 'common/utils/target';
import localStorageConfig from 'common/kfsConfig/localStorage';

let get = null;
let set = null;
let remove = null;

if (isElectron) {
  const fs = window.require('fs');
  const path = window.require('path');
  const remote = window.require('@electron/remote');
  let configPath;
  if (remote.process.env.ELECTRON_START_URL) {
    const cwd = remote.process.cwd();
    configPath = path.join(cwd, 'public', 'extraResources', 'kfs-config.json');
  } else {
    const cwd = remote.process.resourcesPath;
    configPath = path.join(cwd, 'app', 'public', 'extraResources', 'kfs-config.json');
  }
  get = function () {
    try {
      const config = fs.readFileSync(configPath).toString();
      console.log('getExternalConfig', configPath, config);
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
