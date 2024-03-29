import { downloadURI } from "@kfs/common/api/web";
import { httpGet } from "@kfs/common/api/webServer";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { noteError } from "@kfs/mui/components/Notification";
import axios from "axios";

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


export async function getDriverFile(driverId, filePath) {
    console.log('api.getDriverFile', driverId, filePath);
    return await httpGet("/api/v1/driverFile", {
        driverId, filePath
    });
}
