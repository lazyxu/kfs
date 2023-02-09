let ws = require("../rpc/ws");

const testRootDir = [
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
                "Size": 913,
                "Count": 1,
                "TotalCount": 1,
                "CreateTime": 1661133306379099400,
                "ModifyTime": 1661133306379099400,
                "ChangeTime": 1661133306379099400,
                "AccessTime": 1661133306379099400,
                "Content": "console.log(\"in a.js\\n\")",
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
        "Content": "console.log(\"hello, world\\n\")",
    },
]

export function open(sysConfig, branchName, path, onFile, onTotal, onDirItem) {
    console.log('test.open', branchName, path, path.join('/'))
    return listR(testRootDir, path.slice(), onFile, onTotal, onDirItem)
}

function modeIsDir(mode) {
    return mode >= 2147483648;
}

function listR(dirItems, path, onFile, onTotal, onDirItem) {
    if (path.length === 0) {
        let total = dirItems.length;
        onTotal && onTotal(total);
        for (let i = 0; i < total; i++) {
            onDirItem && onDirItem(dirItems[i], i);
        }
        return true;
    }
    for (let i = 0; i < dirItems.length; i++) {
        if (dirItems[i].Name === path[0]) {
            if (modeIsDir(dirItems[i].Mode)) {
                path.shift()
                listR(dirItems[i].DirItems, path, onFile, onTotal, onDirItem)
                return true;
            } else {
                let enc = new TextEncoder("utf-8")
                onFile(enc.encode(dirItems[i].Content));
                return false;
            }
        }
    }
    onTotal && onTotal(0);
    return true;
}

export function list(sysConfig, branchName, path, onTotal, onDirItem) {
    console.log('test.list', branchName, path, path.join('/'))
    return listR(testRootDir, path.slice(), () => {
        onTotal(0)
    }, onTotal, onDirItem)
}
