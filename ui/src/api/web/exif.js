import { httpGet, httpPost } from "./webServer";

export async function analyzeMetadata(start) {
    console.log('web.analyzeMetadata', start);
    return await httpPost("/api/v1/startMetadataAnalysisTask", {
        start,
    });
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