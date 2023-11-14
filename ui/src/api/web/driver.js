import { noteError } from "components/Notification/Notification";
import { httpDelete, httpGet, httpPost } from "./webServer";

export async function listDriver() {
    console.log('web.listDriver');
    return await httpGet("/api/v1/drivers");
}

export async function newDriver(name, description, type, code) {
    console.log('web.newDriver', name, description, type, code);
    return await httpPost("/api/v1/drivers", { name, description, type, code });
}

export async function getDriverSync(id) {
    console.log('web.getDriverSync', id);
    return await httpGet("/api/v1/getDriverSync", { id });
}

export async function getDriverLocalFile(driverId) {
    try {
        console.log('web.getDriverLocalFile', driverId);
        return await httpGet("/api/v1/getDriverLocalFile", { driverId });
    } catch (e) {
        noteError("获取本地文件备份配置失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function updateDriverSync(driverId, sync, h, m) {
    try {
        console.log('web.updateDriverSync', driverId, sync, h, m);
        return await httpGet("/api/v1/updateDriverSync", { driverId, sync, h, m });
    } catch (e) {
        noteError("设置定时同步失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function deleteDriver(id) {
    console.log('web.deleteDriver', id);
    return await httpDelete("/api/v1/drivers", { id });
}

export async function getDriversFileSize(id) {
    console.log('web.getDriversFileSize', id);
    return await httpGet("/api/v1/drivers/fileSize", { id });
}


export async function getDriversFileCount(id) {
    console.log('web.getDriversFileCount', id);
    return await httpGet("/api/v1/drivers/fileCount", { id });
}


export async function getDriversDirCount(id) {
    console.log('web.getDriversDirCount', id);
    return await httpGet("/api/v1/drivers/dirCount", { id });
}
