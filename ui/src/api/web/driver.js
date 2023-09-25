import { httpDelete, httpGet, httpPost } from "./webServer";

export async function listDriver() {
    console.log('web.listDriver');
    return await httpGet("/api/v1/drivers");
}

export async function newDriver(name, description, type, code) {
    console.log('web.newDriver', name, description, type, code);
    return await httpPost("/api/v1/drivers", {name, description, type, code});
}

export async function deleteDriver(name) {
    console.log('web.deleteDriver', name);
    return await httpDelete("/api/v1/drivers", {name});
}
