import { httpGet } from "./localServer";

export async function listBackupTask() {
    console.log('web.listBackupTask');
    return await httpGet("/api/v1/backupTask");
}
