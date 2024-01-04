import { httpGet } from "@kfs/common/api/webServer";

export async function listDriverFileByHash(hash) {
    try {
        console.log('api.listDriverFileByHash', hash);
        return await httpGet("/api/v1/listDriverFileByHash", {
            hash
        });
    } catch (e) {
        window.noteError("查找相同文件失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        throw e;
    }
}