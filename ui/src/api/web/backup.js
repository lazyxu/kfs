import { httpGet, httpPost, httpDelete } from "./localServer";

export async function listBackupTask() {
    console.log('web.listBackupTask');
    return await httpGet("/api/v1/backupTask");
}

export async function newBackupTask(name, description, srcPath, driverName, dstPath, encoder, concurrent) {
    console.log('web.newBackupTask', name, description, srcPath, driverName, dstPath, encoder, concurrent);
    return await httpPost("/api/v1/backupTask", { name, description, srcPath, driverName, dstPath, encoder, concurrent });
}

export async function deleteBackupTask(name) {
    console.log('web.deleteBackupTask', name);
    return await httpDelete("/api/v1/backupTask", { name });
}
