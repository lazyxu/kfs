import {modeIsDir} from "../utils/api";
import moment from "moment";

export const testRootDir = {
    "DirItems": [
        {
            "Hash": "03b756cb01cf80426459ec0dbc1b75bfce640d874fa6797f42a89d216512c224",
            "Name": "lib",
            "Mode": 2147484159,
            "Size": 200762606,
            "Count": 2,
            "TotalCount": 2,
            "CreateTime": 1661863700248713200,
            "ModifyTime": 1661863700248713200,
            "ChangeTime": 1661863700248713200,
            "AccessTime": 1661863700248713200,
            "DirItems": [
                {
                    "Hash": "ee6b3b8a13c0aa770f3b422362aa3e8c57cba7e2c9a5b6635a2eac2fea10dbf8",
                    "Name": "a.js",
                    "Mode": 438,
                    "Size": 101 * 1024 * 1024,
                    "Count": 1,
                    "TotalCount": 1,
                    "CreateTime": 1661133306379099400,
                    "ModifyTime": 1661133306379099400,
                    "ChangeTime": 1661133306379099400,
                    "AccessTime": 1661133306379099400,
                    "Content": new TextEncoder("utf-8").encode("console.log(\"in a.js\\n\")"),
                }
            ],
        },
        {
            "Hash": "ee6b3b8a13c0aa770f3b422362aa3e8c57cba7e2c9a5b6635a2eac2fea10dbf8",
            "Name": "index.js",
            "Mode": 438,
            "Size": 913,
            "Count": 1,
            "TotalCount": 1,
            "CreateTime": 1661133306379099400,
            "ModifyTime": 1661133306379099400,
            "ChangeTime": 1661133306379099400,
            "AccessTime": 1661133306379099400,
            "Content": new TextEncoder("utf-8").encode("console.log(\"hello, world\\n\")"),
        },
    ]
}

export function open(branchName, filePath, onFile, onTotal, onDirItem) {
    console.log('mock.open', branchName, filePath)
    let item = listR(testRootDir, filePath.slice());
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
    let item = listR(testRootDir, filePath.slice());
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
    let item = listR(testRootDir, dirPath.slice());
    if (item.DirItems) {
        item.DirItems.push(newFileItem(fileName, 438));
    }
}

export function newDir(branchName, dirPath, fileName) {
    console.log('mock.newDir', branchName, dirPath, fileName)
    let item = listR(testRootDir, dirPath.slice());
    // TODO: empty or duplicate
    if (item.DirItems) {
        item.DirItems.push(newFileItem(fileName, 2147484159));
    }
}

export function download(branchName, filePath) {
    console.log('mock.download', branchName, filePath)
    let item = listR(testRootDir, filePath.slice());
    if (item.Content) {
        return item.Content;
    }
    return null;
}
