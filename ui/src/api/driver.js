import { noteError } from "components/Notification/Notification";
import { getSysConfig } from "../hox/sysConfig";
import * as mockApi from "./mock/driver";
import * as webApi from "./web/driver";
import { httpDelete, httpGet, httpPost } from "./web/webServer";

export function getDriverApi() {
    return getSysConfig().sysConfig.api === "web" ? webApi : mockApi;
}

export async function listDriver(setResourceManager) {
    try {
        console.log('api.listDriver');
        let drivers = await httpGet("/api/v1/drivers");
        setResourceManager(prev => {
            return { ...prev, drivers, driverId: null, driverName: null, filePath: [], dirItems: null, file: null };
        });
    } catch (e) {
        noteError("获取云盘列表失败：" + (e.response.data ? e.response.data : e.message));
        throw e;
    }
}

export async function newDriver(setResourceManager, name, description, type, code) {
    try {
        console.log('api.newDriver', name, description, type, code);
        let exist = await httpPost("/api/v1/drivers", { name, description, type, code });
        await listDriver(setResourceManager);
        if (exist) {
            throw new Error("云盘名称重复: " + name);
        }
    } catch (e) {
        noteError("创建云盘失败：" + (e.response?.data ? e.response?.data : e.message));
        throw e;
    }
}

export async function deleteDriver(setResourceManager, driverId) {
    try {
        console.log('api.deleteDriver', driverId);
        await httpDelete("/api/v1/drivers", { driverId });
        await listDriver(setResourceManager);
    } catch (e) {
        noteError("删除云盘失败：" + (e.response.data ? e.response.data : e.message));
        throw e;
    }
}
