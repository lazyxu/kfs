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
    try {
        console.log('web.getDriverSync', id);
        return await httpGet("/api/v1/getDriverSync", { id });
    } catch (e) {
        noteError("获取同步配置失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
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

export async function updateDriverLocalFile(driverId, srcPath, ignores, encoder) {
    try {
        console.log('web.updateDriverLocalFile', driverId, srcPath, ignores, encoder);
        return await httpGet("/api/v1/updateDriverLocalFile", { driverId, srcPath, ignores, encoder });
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
    try {
        console.log('web.getDriversFileSize', id);
        return await httpGet("/api/v1/drivers/fileSize", { id });
    } catch (e) {
        noteError("获取云盘文件总大小失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}


export async function getDriversFileCount(id) {
    try {
        console.log('web.getDriversFileCount', id);
        return await httpGet("/api/v1/drivers/fileCount", { id });
    } catch (e) {
        noteError("获取云盘文件数量失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}


export async function getDriversDirCount(id) {
    try {
        console.log('web.getDriversDirCount', id);
        return await httpGet("/api/v1/drivers/dirCount", { id });
    } catch (e) {
        noteError("获取云盘目录数量失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function getDriversDirCalculatedInfo(driverId, filePath) {
    try {
        console.log('web.getDriversDirCalculatedInfo', driverId, filePath);
        return await httpGet("/api/v1/drivers/dirCalculatedInfo", { driverId, filePath });
    } catch (e) {
        noteError("获取云盘目录计算属性失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
