import { noteError } from "components/Notification/Notification";
import { getSysConfig } from "hox/sysConfig";
import { httpPostBody as httpPostBodyLocal } from "./web/localServer";
import { httpDelete, httpGet, httpPost } from "./web/webServer";

export async function listDriver() {
    try {
        console.log('api.listDriver');
        let drivers = await httpGet("/api/v1/drivers");
        return drivers;
    } catch (e) {
        noteError("获取云盘列表失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function openDrivers(setResourceManager) {
    console.log('openDrivers');
    setResourceManager({});
}

export async function newDriver(name, description) {
    try {
        console.log('api.newDriver', name, description);
        let exist = await httpPost("/api/v1/drivers", { name, description });
        if (exist) {
            throw new Error("云盘名称重复: " + name);
        }
    } catch (e) {
        noteError("创建云盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function newDriverBaiduPhoto(setResourceManager, name, description, code) {
    try {
        console.log('api.newDriverBaiduPhoto', name, description, code);
        let exist = await httpPost("/api/v1/driverBaiduPhotos", { name, description, code });
        await listDriver(setResourceManager);
        if (exist) {
            throw new Error("云盘名称重复: " + name);
        }
    } catch (e) {
        noteError("创建云盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function newLocalFileDriver(setResourceManager, name, description, deviceId, srcPath, encoder, concurrent) {
    try {
        console.log('api.newLocalFileDriver', name, description, srcPath, encoder, concurrent);
        let exist = await httpPost("/api/v1/driverLocalFiles", { name, description, deviceId, srcPath, encoder, concurrent });
        await listDriver(setResourceManager);
        if (exist) {
            throw new Error("云盘名称重复: " + name);
        }
    } catch (e) {
        noteError("创建云盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function resetDriver(driverId) {
    try {
        console.log('api.deleteDriver', driverId);
        await httpGet("/api/v1/drivers/reset", { driverId });
    } catch (e) {
        noteError("重置云盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function deleteDriver(setResourceManager, driverId) {
    try {
        console.log('api.deleteDriver', driverId);
        await httpDelete("/api/v1/drivers", { driverId });
        await listDriver(setResourceManager);
    } catch (e) {
        noteError("删除云盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function listLocalFileDriver(deviceId) {
    try {
        console.log('api.listLocalFileDriver', deviceId);
        return await httpGet("/api/v1/listLocalFileDriver", { deviceId });
    } catch (e) {
        noteError("获取本地文件备份盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function startAllLocalFileSync(drivers) {
    try {
        let serverAddr = getSysConfig().sysConfig.socketServer;
        console.log('api.startAllLocalFileSync', serverAddr, drivers);
        return await httpPostBodyLocal("/api/v1/startAllLocalFileSync", { serverAddr, drivers });
    } catch (e) {
        noteError("启动本地文件备份盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
