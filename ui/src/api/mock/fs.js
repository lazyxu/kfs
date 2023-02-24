import {modeIsDir} from "../utils/api";
import moment from "moment";
import {branches} from './data';

function getBranch(name) {
    for (let i = 0; i < branches.length; i++) {
        const branch = branches[i];
        if (branch.name === name) {
            return branch;
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

export async function listCb(branchName, filePath, onTotal, onDirItem) {
    console.log('mock.listCb', branchName, filePath)
    const branch = getBranch(branchName);
    if (!branch) {
        onTotal?.(0);
        return;
    }
    let item = listR(branch, filePath.slice());
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

export async function list(branchName, filePath, onTotal, onDirItem) {
    console.log('mock.list', branchName, filePath)
    const branch = getBranch(branchName);
    if (!branch) {
        return [];
    }
    let item = listR(branch, filePath.slice());
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

export async function newFile(branchName, dirPath, fileName) {
    console.log('mock.newFile', branchName, dirPath, fileName)
    const branch = getBranch(branchName);
    if (!branch) {
        return;
    }
    let item = listR(branch, dirPath.slice());
    if (item.dirItems) {
        item.dirItems.push(newFileItem(fileName, 438));
    }
}

export async function newDir(branchName, dirPath, fileName) {
    console.log('mock.newDir', branchName, dirPath, fileName)
    const branch = getBranch(branchName);
    if (!branch) {
        return;
    }
    let item = listR(branch, dirPath.slice());
    // TODO: empty or duplicate
    if (item.dirItems) {
        item.dirItems.push(newFileItem(fileName, 2147484159));
    }
}

export async function download(branchName, filePath) {
    console.log('mock.download', branchName, filePath)
    const branch = getBranch(branchName);
    if (!branch) {
        return;
    }
    let item = listR(branch, filePath.slice());
    if (item.Content) {
        return item.Content;
    }
    return null;
}
