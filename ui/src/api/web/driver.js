import {httpDelete, httpGet, httpPost} from "./webServer";

export async function listDriver() {
    console.log('web.listDriver');
    return await httpGet("/api/v1/drivers");
}

export async function newDriver(name) {
    console.log('web.newDriver', name);
    return await httpPost("/api/v1/drivers", {name});
}

export async function deleteDriver(name) {
    console.log('web.deleteDriver', name);
    return await httpDelete("/api/v1/drivers", {name});
}
