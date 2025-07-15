import { httpDelete, httpGet, httpPostJson } from "@kfs/common/api/webServer";
import { noteError } from "@kfs/mui/components/Notification";

export async function listDevice(setDevices) {
    try {
        console.log('api.listDevice');
        let list = await httpGet("/api/v1/devices");
        setDevices(list);
    } catch (e) {
        noteError("获取设备列表失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function newDevice(id, name, os, user_agent, host_name) {
    try {
        console.log('api.newDevice', id, name, os, user_agent, host_name);
        await httpPostJson("/api/v1/devices", { id, name, os, user_agent, host_name });
        return id;
    } catch (e) {
        noteError("创建设备失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}

export async function deleteDevice(setDevices, deviceId) {
    try {
        console.log('api.deleteDevice', deviceId);
        await httpDelete("/api/v1/devices", { deviceId });
        await listDevice(setDevices);
    } catch (e) {
        noteError("删除设备失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
