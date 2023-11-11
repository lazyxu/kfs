import { noteError } from "components/Notification/Notification";
import { getSysConfig } from "../hox/sysConfig";
import * as mockApi from "./mock/fs";
import * as webApi from "./web/fs";

function getFsApi() {
    return getSysConfig().sysConfig.api === "web" ? webApi : mockApi;
}

export async function list(setResourceManager, driverId, driverName, filePath) {
    try {
        console.log('api.list', driverId, filePath);
        let dirItems = await getFsApi().list(driverId, filePath);
        setResourceManager(prev => {
            return { ...prev, driverId, driverName, filePath, dirItems, file: null, drivers: null };
        });
    } catch (e) {
        noteError(e.response.data ? e.response.data : e.message);
    }
}

export async function openFile(setResourceManager, driverId, filePath, dirItem) {
    console.log('api.openFile', driverId, filePath);
    let { content, tooLarge } = await getFsApi().openFile(driverId, filePath);
    dirItem.content = content;
    dirItem.tooLarge = tooLarge;
    setResourceManager(prev => {
        return {
            ...prev, driverId, filePath,
            dirItems: null, drivers: null,
            file: dirItem,
        };
    });
}

export async function newFile(setResourceManager, driverId, driverName, dirPath, fileName) {
    console.log('api.newFile', driverId, dirPath, fileName);
    await getFsApi().newFile(driverId, dirPath, fileName);
    await list(setResourceManager, driverId, dirPath);
}

export async function newDir(setResourceManager, driverId, dirPath, fileName) {
    console.log('api.newDir', driverId, dirPath, fileName);
    await getFsApi().newDir(driverId, dirPath, fileName);
    await list(setResourceManager, driverId, dirPath);
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

export async function download(driverId, filePath) {
    console.log('api.download', driverId, filePath);
    let data = await getFsApi().download(driverId, filePath);
    downloader(data, filePath[filePath.length - 1]);
}
