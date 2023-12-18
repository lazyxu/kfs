import { httpGet, httpPost } from "@kfs/common/api/webServer";

export async function analysisExif(start) {
    try {
        console.log('api.analysisExif', start);
        return await httpPost("/api/v1/analysisExif", {
            start,
        });
    } catch (e) {
        window.noteError("分析图片信息失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function exifStatus() {
    try {
        console.log('api.exifStatus');
        return await httpGet("/api/v1/analysisExif");
    } catch (e) {
        window.noteError("获取图片信息分析状态失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function listExif() {
    try {
        console.log('api.exif');
        return await httpGet("/api/v1/exif");
    } catch (e) {
        window.noteError("获取图片信息失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
