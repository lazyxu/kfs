import { noteError } from "components/Notification/Notification";
import { getSysConfig } from "hox/sysConfig";
import { httpPostBody as httpPostBodyLocal } from "./web/localServer";
import { httpDelete, httpGet, httpPost } from "./web/webServer";

export async function listDriver(setResourceManager) {
    try {
        console.log('api.listDriver');
        let drivers = await httpGet("/api/v1/drivers");
        setResourceManager(prev => {
            return { ...prev, drivers, driverId: null, driverName: null, filePath: [], dirItems: null, file: null };
        });
    } catch (e) {
        noteError("获取云盘列表失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function newDriver(setResourceManager, name, description) {
    try {
        console.log('api.newDriver', name, description);
        let exist = await httpPost("/api/v1/drivers", { name, description });
        await listDriver(setResourceManager);
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

export async function startAllLocalFileSync(driverId, drivers) {
    try {
        let serverAddr = getSysConfig().sysConfig.socketServer;
        console.log('api.startAllLocalFileSync', driverId, serverAddr, drivers);
        return await httpPostBodyLocal("/api/v1/startAllLocalFileSync", { driverId, serverAddr, drivers });
    } catch (e) {
        noteError("启动本地文件备份盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
