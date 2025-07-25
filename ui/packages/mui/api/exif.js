import { httpPost as localHttpPost } from "@kfs/common/api/localServer";
import { httpPost } from "@kfs/common/api/webServer";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { noteError } from "@kfs/mui/components/Notification";

export async function analyzeMetadata(start, all) {
    console.log('web.analyzeMetadata', start, all);
    return await httpPost("/api/v1/startMetadataAnalysisTask", {
        start, all,
    });
}

export async function analyzeTranscode(start, all) {
    console.log('web.analyzeTranscode', start, all);
    return await httpPost("/api/v1/startTranscodeAnalysisTask", {
        start, all,
    });
}

export async function startLivePhotoAnalysisTask(start, all) {
    console.log('web.startLivePhotoAnalysisTask', start, all);
    return await httpPost("/api/v1/startLivePhotoAnalysisTask", {
        start, all,
    });
}

export async function truncateLivePhoto() {
    console.log('web.truncateLivePhoto');
    return await httpPost("/api/v1/truncateLivePhoto", {});
}

export async function startBaiduPhotoTask(start, id) {
    try {
        console.log('web.startBaiduPhotoTask', start, id);
        return await httpPost("/api/v1/startBaiduPhotoTask", {
            start, id,
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