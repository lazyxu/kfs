import { noteError } from "@kfs/mui/components/Notification";

Object.defineProperty(window, 'kfsEnv', {
    value: import.meta.env,
    enumerable: false,
    configurable: false,
    writable: false,
});

const fs = window.require('fs');
const path = window.require('path');
const remote = window.require('@electron/remote');
const { app } = window.require('@electron/remote');

const configFilename = 'kfs-config.json';
const userData = app.getPath('userData');
console.log('userData', userData);
const configPath = path.join(userData, configFilename);

console.log(configPath);

Object.defineProperty(window, "kfsConfig", {
    get() {
        try {
            const config = fs.readFileSync(configPath).toString();
            console.log('getExternalConfig', configPath, config);
            const obj = JSON.parse(config);
            return obj;
        } catch (e) {
            return undefined;
        }
    },
    set(json) {
        fs.writeFileSync(configPath, JSON.stringify(json, undefined, 2), { flag: 'w+' });
    },
});

window.noteError = noteError;
