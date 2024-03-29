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

export async function getMetadata(hash) {
    try {
        console.log('web.getMetadata', hash);
        return await httpGet("/api/v1/metadata", {
            hash,
        });
    } catch (e) {
        window.noteError("获取文件元数据失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function listDCIMMetadataTime() {
    try {
        console.log('api.listDCIMMetadataTime');
        return await httpGet("/api/v1/DCIMMetadataTime");
    } catch (e) {
        window.noteError("获取图片信息失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}


export async function listDCIMMediaType() {
    try {
        console.log('api.listDCIMMediaType');
        return await httpGet("/api/v1/listDCIMMediaType");
    } catch (e) {
        window.noteError("获取媒体类型失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
