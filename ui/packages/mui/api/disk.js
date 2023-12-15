import { noteError } from "@kfs/mui/components/Notification/Notification";
import { httpGet } from "./webServer";

export async function getDiskUsage() {
    try {
        console.log('web.diskUsage');
        return await httpGet("/api/v1/diskUsage");
    } catch (e) {
        noteError("获取存储空间失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}
