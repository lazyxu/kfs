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

export async function uploadFile(driverId, filePath, file) {
    try {
        const formData = new FormData();
        formData.append("driverId", driverId);
        formData.append("filePath", JSON.stringify(filePath));
        formData.append("file", file);
        await axios.post(`${getSysConfig().webServer}/api/v1/uploadFile`, formData, {
            headers: {
                "Content-Type": "multipart/form-data"
            }
        });
    } catch (e) {
        noteError("上传文件失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function checkFileExists(dirPath, name, sha256, md5) {
    console.log('api.checkFileExists', dirPath, name, sha256, md5);
    try {
        const resp = await axios.get(`${getSysConfig().webServer}/api/v1/checkFileExists`, {
            params: {
                dirPath, name,
                sha256, md5,
            }
        });
        return resp.data.exists;
    } catch (e) {
        noteError("检查文件失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function createFolder(dirPath, folderName) {
    try {
        await axios.post(`${getSysConfig().webServer}/api/v1/createFolder`, {
            dirPath,
            folderName
        });
    } catch (e) {
        noteError("创建文件夹失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function mergeFileChunks( dirPath, hash, totalChunks, fileMeta) {
    try {
        await axios.post(`${getSysConfig().webServer}/api/v1/mergeFileChunks`, {
            dirPath,
            hash,
            totalChunks,
            fileMeta,
        });
    } catch (e) {
        noteError("合并分块失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
