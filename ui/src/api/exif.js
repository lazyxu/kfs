import { noteError } from "components/Notification/Notification";
import { getSysConfig } from "hox/sysConfig";
import { httpPost as localHttpPost } from "./localServer";
import { httpGet, httpPost } from "./webServer";

export async function analyzeMetadata(start) {
    console.log('web.analyzeMetadata', start);
    return await httpPost("/api/v1/startMetadataAnalysisTask", {
        start,
    });
}

export async function startBaiduPhotoTask(start, driverId) {
    try {
        console.log('web.startBaiduPhotoTask', start, driverId);
        return await httpPost("/api/v1/startBaiduPhotoTask", {
            start, driverId,
        });
    } catch (e) {
        noteError("备份一刻相册失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function startDriverLocalFile(start, driverId, srcPath, ignores, encoder) {
    try {
        let serverAddr = getSysConfig().sysConfig.socketServer;
        console.log('web.startDriverLocalFile', start, driverId, serverAddr, srcPath, ignores, encoder);
        return await localHttpPost("/api/v1/startDriverLocalFile", {
            start, driverId, serverAddr, srcPath, ignores, encoder
        });
    } catch (e) {
        noteError("备份本地文件失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function startDriverLocalFileFilter(start, driverId, srcPath, ignores) {
    try {
        console.log('web.startDriverLocalFileFilter', start, driverId, srcPath, ignores);
        return await localHttpPost("/api/v1/startDriverLocalFileFilter", {
            start, driverId, srcPath, ignores,
        });
    } catch (e) {
        noteError("测试过滤规则设置失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function analysisExif(start) {
    console.log('web.analysisExif', start);
    return await httpPost("/api/v1/analysisExif", {
        start,
    });
}

export async function exifStatus() {
    console.log('web.exifStatus');
    return await httpGet("/api/v1/analysisExif");
}

export async function listExif() {
    console.log('web.exif');
    return await httpGet("/api/v1/exif");
}

export async function getMetadata(hash) {
    console.log('web.getMetadata', hash);
    return await httpGet("/api/v1/metadata", {
        hash,
    });
}