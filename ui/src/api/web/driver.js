import { httpDelete, httpGet, httpPost } from "./webServer";

export async function listDriver() {
    console.log('web.listDriver');
    return await httpGet("/api/v1/drivers");
}

export async function newDriver(name, description, type, code) {
    console.log('web.newDriver', name, description, type, code);
    return await httpPost("/api/v1/drivers", {name, description, type, code});
}

export async function getDriverSync(name) {
    console.log('web.getDriverSync', name);
    return await httpGet("/api/v1/getDriverSync", {name});
}

export async function updateDriverSync(name, sync, h, m, s) {
    console.log('web.updateDriverSync', name, sync, h, m, s);
    return await httpGet("/api/v1/updateDriverSync", {name, sync, h, m, s});
}

export async function deleteDriver(name) {
    console.log('web.deleteDriver', name);
    return await httpDelete("/api/v1/drivers", {name});
}

export async function getDriversFileSize(name) {
    console.log('web.getDriversFileSize', name);
    return await httpGet("/api/v1/drivers/fileSize", {name});
}


export async function getDriversFileCount(name) {
    console.log('web.getDriversFileCount', name);
    return await httpGet("/api/v1/drivers/fileCount", {name});
}


export async function getDriversDirCount(name) {
    console.log('web.getDriversDirCount', name);
    return await httpGet("/api/v1/drivers/dirCount", {name});
}
