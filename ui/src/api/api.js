let mockApi = require("./mock/api");

export async function open(sysConfig, setResourceManager, branchName, filePath) {
    console.log('api.open', branchName, filePath, filePath.join('/'));
    let dirItems;
    let isDir = await mockApi.open(branchName, filePath, (file) => {
        setResourceManager(prev => {
            return {
                ...prev, branchName, filePath,
                dirItems: null,
                file,
            };
        });
    }, (total) => {
        dirItems = new Array(total);
    }, (dirItem, i) => {
        dirItems[i] = dirItem;
    });
    if (isDir) {
        setResourceManager(prev => {
            return {
                ...prev, branchName, filePath,
                dirItems: dirItems ? dirItems : prev.dirItems,
                file: null,
            };
        });
    }
}

export async function list(sysConfig, setResourceManager, branchName, filePath) {
    console.log('api.list', branchName, filePath, filePath.join('/'));
    let dirItems;
    await mockApi.list(branchName, filePath, (total) => {
        dirItems = new Array(total);
    }, (dirItem, i) => {
        dirItems[i] = dirItem;
    });
    setResourceManager(prev => {
        return {...prev, branchName, filePath, dirItems, file: null};
    });
}

export async function newFile(sysConfig, setResourceManager, branchName, filePath) {
    console.log('api.newFile', branchName, filePath, filePath.join('/'));
    await mockApi.newFile(branchName, filePath);
    await list(sysConfig, setResourceManager, branchName, filePath)
}

export async function newDir(sysConfig, setResourceManager, branchName, filePath) {
    console.log('api.newDir', branchName, filePath, filePath.join('/'));
    await mockApi.newDir(branchName, filePath);
    await list(sysConfig, setResourceManager, branchName, filePath)
}

function downloadURI(uri, name) {
    let link = document.createElement("a");
    link.download = name;
    link.href = uri;
    link.click();
}

function downloader(data, name) {
    let blob = new Blob([data]);
    let url = window.URL.createObjectURL(blob);
    downloadURI(url, name);
    window.URL.revokeObjectURL(url);
}

export async function download(sysConfig, branchName, filePath) {
    console.log('api.download', branchName, filePath, filePath.join('/'));
    let data = await mockApi.download(branchName, filePath);
    downloader(data, filePath[filePath.length - 1]);
}
