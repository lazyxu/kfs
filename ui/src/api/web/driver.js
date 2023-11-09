import { httpDelete, httpGet, httpPost } from "./webServer";

export async function listDriver() {
    console.log('web.listDriver');
    return await httpGet("/api/v1/drivers");
}

export async function newDriver(name, description, type, code) {
    console.log('web.newDriver', name, description, type, code);
    return await httpPost("/api/v1/drivers", {name, description, type, code});
}

export async function getDriverSync(id) {
    console.log('web.getDriverSync', id);
    return await httpGet("/api/v1/getDriverSync", {id});
}

export async function updateDriverSync(id, sync, h, m, s) {
    console.log('web.updateDriverSync', id, sync, h, m, s);
    return await httpGet("/api/v1/updateDriverSync", {id, sync, h, m, s});
}

export async function deleteDriver(id) {
    console.log('web.deleteDriver', id);
    return await httpDelete("/api/v1/drivers", {id});
}

export async function getDriversFileSize(id) {
    console.log('web.getDriversFileSize', id);
    return await httpGet("/api/v1/drivers/fileSize", {id});
}


export async function getDriversFileCount(id) {
    console.log('web.getDriversFileCount', id);
    return await httpGet("/api/v1/drivers/fileCount", {id});
}


export async function getDriversDirCount(id) {
    console.log('web.getDriversDirCount', id);
    return await httpGet("/api/v1/drivers/dirCount", {id});
}
