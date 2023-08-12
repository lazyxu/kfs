import * as mockApi from "./mock/fs";
import * as webApi from "./web/fs";
import {getSysConfig} from "../hox/sysConfig";

function getFsApi() {
    return getSysConfig().sysConfig.api === "web" ? webApi : mockApi;
}

export async function list(setResourceManager, driverName, filePath) {
    console.log('api.list', driverName, filePath);
    let dirItems = await getFsApi().list(driverName, filePath);
    setResourceManager(prev => {
        return {...prev, driverName, filePath, dirItems, file: null, drivers: null};
    });
}

export async function openFile(setResourceManager, driverName, filePath, dirItem) {
    console.log('api.openFile', driverName, filePath);
    let {content, tooLarge} = await getFsApi().openFile(driverName, filePath);
    dirItem.content = content;
    dirItem.tooLarge = tooLarge;
    setResourceManager(prev => {
        return {
            ...prev, driverName, filePath,
            dirItems: null, drivers: null,
            file: dirItem,
        };
    });
}

export async function newFile(setResourceManager, driverName, dirPath, fileName) {
    console.log('api.newFile', driverName, dirPath, fileName);
    await getFsApi().newFile(driverName, dirPath, fileName);
    await list(setResourceManager, driverName, dirPath)
}

export async function newDir(setResourceManager, driverName, dirPath, fileName) {
    console.log('api.newDir', driverName, dirPath, fileName);
    await getFsApi().newDir(driverName, dirPath, fileName);
    await list(setResourceManager, driverName, dirPath)
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

export async function download(driverName, filePath) {
    console.log('api.download', driverName, filePath);
    let data = await getFsApi().download(driverName, filePath);
    downloader(data, filePath[filePath.length - 1]);
}
