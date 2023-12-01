import axios from "axios";
import { noteError } from "components/Notification/Notification";
import { getSysConfig } from "../hox/sysConfig";
import * as mockApi from "./mock/fs";
import * as webApi from "./web/fs";
import { httpGet } from "./web/webServer";

function getFsApi() {
    return getSysConfig().sysConfig.api === "web" ? webApi : mockApi;
}

export async function list(setResourceManager, driverId, driverName, filePath) {
    try {
        console.log('api.list', driverId, filePath);
        let dirItems = await httpGet("/api/v1/list", { driverId, filePath });
        setResourceManager(prev => {
            return { ...prev, driverId, driverName, filePath, dirItems, file: null, drivers: null };
        });
    } catch (e) {
        noteError("获取文件列表失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function openDir(setResourceManager, driverId, driverName, filePath) {
    setResourceManager({ driverId, driverName, filePath });
}

export async function openFile(setResourceManager, driverId, filePath, dirItem) {
    try {
        console.log('api.openFile', driverId, filePath);
        let resp = await axios.get(`${getSysConfig().sysConfig.webServer}/api/v1/openFile`, {
            params: {
                driverId,
                filePath,
                maxContentSize: getSysConfig().sysConfig.maxContentSize,
            },
            transformResponse: x => x,
        });
        let tooLarge = resp.headers.get("Kfs-Too-Large");
        let content = resp.data;
        dirItem.content = content;
        dirItem.tooLarge = tooLarge;
        setResourceManager(prev => {
            return {
                ...prev, driverId, filePath,
                dirItems: null, drivers: null,
                file: dirItem,
            };
        });
    } catch (e) {
        noteError("打开文件失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
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
    try {
        console.log('api.download', driverId, filePath);
        let resp = await axios.get(`${getSysConfig().sysConfig.webServer}/api/v1/downloadFile`, {
            params: {
                driverId,
                filePath,
            },
            responseType: "arraybuffer",
        });
        downloader(resp.data, filePath[filePath.length - 1]);
    } catch (e) {
        noteError("打开文件失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
