import {modeIsDir} from "../utils/api";
import moment from "moment";
import {drivers} from './data';

function getDriver(name) {
    for (let i = 0; i < drivers.length; i++) {
        const driver = drivers[i];
        if (driver.name === name) {
            return driver;
        }
    }
    return null;
}

// returns isDir
function listR(dir, filePath) {
    if (filePath.length === 0) {
        return dir;
    }
    for (let i = 0; i < dir.dirItems.length; i++) {
        let item = dir.dirItems[i];
        if (item.Name === filePath[0]) {
            if (modeIsDir(item.Mode)) {
                filePath.shift()
                return listR(item, filePath)
            } else {
                return item;
            }
        }
    }
    return null;
}

export async function listCb(driverName, filePath, onTotal, onDirItem) {
    console.log('mock.listCb', driverName, filePath)
    const driver = getDriver(driverName);
    if (!driver) {
        onTotal?.(0);
        return;
    }
    let item = listR(driver, filePath.slice());
    if (item.dirItems) {
        let total = item.dirItems.length;
        onTotal?.(total);
        for (let i = 0; i < total; i++) {
            onDirItem?.(item.dirItems[i], i);
        }
        return;
    }
    // TODO: no such dir.
    onTotal?.(0);
}

export async function list(driverName, filePath, onTotal, onDirItem) {
    console.log('mock.list', driverName, filePath)
    const driver = getDriver(driverName);
    if (!driver) {
        return [];
    }
    let item = listR(driver, filePath.slice());
    return item.dirItems;
}

function newFileItem(name, mode) {
    let now = moment().valueOf() * 1000 * 1000;
    let file = {
        "Hash": "ee6b3b8a13c0aa770f3b422362aa3e8c57cba7e2c9a5b6635a2eac2fea10dbf8",
        "Name": name,
        "Mode": mode,
        "Size": 0,
        "Count": 1,
        "TotalCount": 1,
        "CreateTime": now,
        "ModifyTime": now,
        "ChangeTime": now,
        "AccessTime": now,
    }
    if (modeIsDir(mode)) {
        file.dirItems = [];
    } else {
        file.Content = new TextEncoder("utf-8").encode("");
    }
    return file;
}

function addNewFile(item, name, mode) {
    let names = {};
    for (let i = 0; i < item.dirItems.length; i++) {
        names[item.dirItems[i].Name] = true;
    }
    let id = 0;
    let tempName;
    while (1) {
        tempName = name;
        if (id !== 0) {
            tempName += " (" + id + ")";
        }
        if (!names[tempName]) {
            break;
        }
        id++;
    }
    item.dirItems.push(newFileItem(tempName, mode));
}

export async function newFile(driverName, dirPath, fileName) {
    console.log('mock.newFile', driverName, dirPath, fileName)
    const driver = getDriver(driverName);
    if (!driver) {
        return;
    }
    let item = listR(driver, dirPath.slice());
    if (item.dirItems) {
        item.dirItems.push(newFileItem(fileName, 438));
    }
}

export async function newDir(driverName, dirPath, fileName) {
    console.log('mock.newDir', driverName, dirPath, fileName)
    const driver = getDriver(driverName);
    if (!driver) {
        return;
    }
    let item = listR(driver, dirPath.slice());
    // TODO: empty or duplicate
    if (item.dirItems) {
        item.dirItems.push(newFileItem(fileName, 2147484159));
    }
}

export async function download(driverName, filePath) {
    console.log('mock.download', driverName, filePath)
    const driver = getDriver(driverName);
    if (!driver) {
        return;
    }
    let item = listR(driver, filePath.slice());
    if (item.Content) {
        return item.Content;
    }
    return null;
}
