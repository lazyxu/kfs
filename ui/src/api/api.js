import * as mockApi from "./mock/api";

export async function open(setResourceManager, branchName, filePath) {
    console.log('api.open', branchName, filePath);
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

export async function list(setResourceManager, branchName, filePath) {
    console.log('api.list', branchName, filePath);
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

export async function newFile(setResourceManager, branchName, dirPath, fileName) {
    console.log('api.newFile', branchName, dirPath, fileName);
    await mockApi.newFile(branchName, dirPath, fileName);
    await list(setResourceManager, branchName, dirPath)
}

export async function newDir(setResourceManager, branchName, dirPath, fileName) {
    console.log('api.newDir', branchName, dirPath, fileName);
    await mockApi.newDir(branchName, dirPath, fileName);
    await list(setResourceManager, branchName, dirPath)
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

export async function download(branchName, filePath) {
    console.log('api.download', branchName, filePath);
    let data = await mockApi.download(branchName, filePath);
    downloader(data, filePath[filePath.length - 1]);
}
