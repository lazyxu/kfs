import * as mockApi from "./mock/fs";
import * as webApi from "./web/fs";
import {getSysConfig} from "../hox/sysConfig";

function getFsApi() {
    return getSysConfig().sysConfig.api === "web" ? webApi : mockApi;
}

export async function open(setResourceManager, branchName, filePath) {
    console.log('api.open', branchName, filePath);
    let dirItem = await getFsApi().open(branchName, filePath);
    if (dirItem.dirItems) {
        setResourceManager(prev => {
            return {
                ...prev, branchName, filePath,
                dirItems: dirItem.dirItems,
                file: null, branches: null,
            };
        });
        return;
    }
    if (dirItem.content) {
        setResourceManager(prev => {
            return {
                ...prev, branchName, filePath,
                dirItems: null, branches: null,
                file: dirItem,
            };
        });

    }
}

export async function list(setResourceManager, branchName, filePath) {
    console.log('api.list', branchName, filePath);
    let dirItems = await getFsApi().list(branchName, filePath);
    setResourceManager(prev => {
        return {...prev, branchName, filePath, dirItems, file: null, branches: null};
    });
}

export async function openFile(setResourceManager, branchName, filePath, dirItem) {
    console.log('api.openFile', branchName, filePath);
    let {content, tooLarge} = await getFsApi().openFile(branchName, filePath);
    dirItem.content = content;
    dirItem.tooLarge = tooLarge;
    setResourceManager(prev => {
        return {
            ...prev, branchName, filePath,
            dirItems: null, branches: null,
            file: dirItem,
        };
    });
}

export async function newFile(setResourceManager, branchName, dirPath, fileName) {
    console.log('api.newFile', branchName, dirPath, fileName);
    await getFsApi().newFile(branchName, dirPath, fileName);
    await list(setResourceManager, branchName, dirPath)
}

export async function newDir(setResourceManager, branchName, dirPath, fileName) {
    console.log('api.newDir', branchName, dirPath, fileName);
    await getFsApi().newDir(branchName, dirPath, fileName);
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
    let data = await getFsApi().download(branchName, filePath);
    downloader(data, filePath[filePath.length - 1]);
}
