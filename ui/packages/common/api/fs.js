import { noteError } from "@kfs/common/components/Notification/Notification";
import axios from "axios";
import { getSysConfig } from "../hox/sysConfig";
import { httpGet } from "./webServer";

export async function openFile(driverId, filePath) {
    try {
        console.log('api.openFile', driverId, filePath);
        let resp = await axios.get(`${getSysConfig().webServer}/api/v1/openFile`, {
            params: {
                driverId,
                filePath,
                maxContentSize: getSysConfig().maxContentSize,
            },
            transformResponse: x => x,
        });
        let tooLarge = resp.headers.get("Kfs-Too-Large");
        let content = resp.data;
        return { content, tooLarge };
    } catch (e) {
        noteError("加载文件失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
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
        let resp = await axios.get(`${getSysConfig().webServer}/api/v1/downloadFile`, {
            params: {
                driverId,
                filePath,
            },
            responseType: "arraybuffer",
        });
        downloader(resp.data, filePath[filePath.length - 1]);
    } catch (e) {
        noteError("下载文件失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function downloadByHash(hash, name) {
    try {
        console.log('api.downloadByHash', hash, name);
        let resp = await axios.get(`${getSysConfig().webServer}/api/v1/image?hash=${hash}`, {
            responseType: "arraybuffer",
        });
        downloader(resp.data, name);
    } catch (e) {
        noteError("下载文件失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}


export async function listDriverFileByHash(hash) {
    console.log('api.listDriverFileByHash', hash);
    return await httpGet("/api/v1/listDriverFileByHash", {
        hash
    });
}

export async function getDriverFile(driverId, filePath) {
    console.log('api.getDriverFile', driverId, filePath);
    return await httpGet("/api/v1/driverFile", {
        driverId, filePath
    });
}
