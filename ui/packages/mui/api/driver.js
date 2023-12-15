import { noteError } from "@kfs/mui/components/Notification/Notification";
import { getSysConfig } from "@kfs/mui/hox/sysConfig";
import { httpPostBody as httpPostBodyLocal } from "./localServer";
import { httpDelete, httpGet, httpPost } from "./webServer";

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

export async function newDriverBaiduPhoto(name, description, code) {
    try {
        console.log('api.newDriverBaiduPhoto', name, description, code);
        let exist = await httpPost("/api/v1/driverBaiduPhotos", { name, description, code });
        if (exist) {
            throw new Error("云盘名称重复: " + name);
        }
    } catch (e) {
        noteError("创建云盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function newLocalFileDriver(name, description, deviceId, srcPath, encoder, concurrent) {
    try {
        console.log('api.newLocalFileDriver', name, description, srcPath, encoder, concurrent);
        let exist = await httpPost("/api/v1/driverLocalFiles", { name, description, deviceId, srcPath, encoder, concurrent });
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

export async function deleteDriver(driverId) {
    try {
        console.log('api.deleteDriver', driverId);
        await httpDelete("/api/v1/drivers", { driverId });
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
        let serverAddr = getSysConfig().socketServer;
        console.log('api.startAllLocalFileSync', serverAddr, drivers);
        return await httpPostBodyLocal("/api/v1/startAllLocalFileSync", { serverAddr, drivers });
    } catch (e) {
        noteError("启动本地文件备份盘失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
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

export async function getDriversDirCalculatedInfo(driverId, filePath) {
    try {
        console.log('web.getDriversDirCalculatedInfo', driverId, filePath);
        return await httpGet("/api/v1/drivers/dirCalculatedInfo", { driverId, filePath });
    } catch (e) {
        noteError("获取云盘目录计算属性失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
