const fs = window.require('fs');
const path = window.require('path');
const remote = window.require('@electron/remote');

let configFilename = 'kfs-config.json';
let configPath = path.join(remote.process.resourcesPath, configFilename);

console.log(configPath);

function get() {
    try {
        const config = fs.readFileSync(configPath).toString();
        console.log('getExternalConfig', configPath, config);
        const obj = JSON.parse(config);
        return obj;
    } catch (e) {
        return undefined;
    }
};

function set(json) {
    fs.writeFileSync(configPath, JSON.stringify(json, undefined, 2), { flag: 'w+' });
};

function remove() {
    fs.unlinkSync(configPath);
};

export default { get, set, remove };
