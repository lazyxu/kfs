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

export function open(branchName, filePath, onFile, onTotal, onDirItem) {
    console.log('mock.open', branchName, filePath);
    const branch = getBranch(branchName);
    if (!branch) {
        return false;
    }
    let item = listR(branch, filePath.slice());
    if (item.DirItems) {
        let total = item.DirItems.length;
        onTotal?.(total);
        for (let i = 0; i < total; i++) {
            onDirItem?.(item.DirItems[i], i);
        }
        return true;
    }
    onFile(item);
    return false;
}

// returns isDir
function listR(dir, filePath) {
    if (filePath.length === 0) {
        return dir;
    }
    for (let i = 0; i < dir.DirItems.length; i++) {
        let item = dir.DirItems[i];
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

export function list(branchName, filePath, onTotal, onDirItem) {
    console.log('mock.list', branchName, filePath)
    const branch = getBranch(branchName);
    if (!branch) {
        onTotal?.(0);
        return;
    }
    let item = listR(branch, filePath.slice());
    if (item.DirItems) {
        let total = item.DirItems.length;
        onTotal?.(total);
        for (let i = 0; i < total; i++) {
            onDirItem?.(item.DirItems[i], i);
        }
        return;
    }
    // TODO: no such dir.
    onTotal?.(0);
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
        file.DirItems = [];
    } else {
        file.Content = new TextEncoder("utf-8").encode("");
    }
    return file;
}

function addNewFile(item, name, mode) {
    let names = {};
    for (let i = 0; i < item.DirItems.length; i++) {
        names[item.DirItems[i].Name] = true;
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
    item.DirItems.push(newFileItem(tempName, mode));
}

export function newFile(branchName, dirPath, fileName) {
    console.log('mock.newFile', branchName, dirPath, fileName)
    const branch = getBranch(branchName);
    if (!branch) {
        return;
    }
    let item = listR(branch, dirPath.slice());
    if (item.DirItems) {
        item.DirItems.push(newFileItem(fileName, 438));
    }
}

export function newDir(branchName, dirPath, fileName) {
    console.log('mock.newDir', branchName, dirPath, fileName)
    const branch = getBranch(branchName);
    if (!branch) {
        return;
    }
    let item = listR(branch, dirPath.slice());
    // TODO: empty or duplicate
    if (item.DirItems) {
        item.DirItems.push(newFileItem(fileName, 2147484159));
    }
}

export function download(branchName, filePath) {
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
