import { httpPost as localHttpPost } from "@kfs/common/api/localServer";
import { httpGet, httpPost } from "@kfs/common/api/webServer";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { noteError } from "@kfs/mui/components/Notification";

export async function analyzeMetadata(start, force) {
    console.log('web.analyzeMetadata', start, force);
    return await httpPost("/api/v1/startMetadataAnalysisTask", {
        start, force,
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
        let serverAddr = getSysConfig().socketServer;
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

export async function getMetadata(hash) {
    try {
        console.log('web.getMetadata', hash);
        return await httpGet("/api/v1/metadata", {
            hash,
        });
    } catch (e) {
        noteError("获取文件元数据失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}