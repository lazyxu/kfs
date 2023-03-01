import {httpDelete, httpGet, httpPost} from "./common";

export async function listBranch() {
    console.log('web.listBranch');
    return await httpGet("/api/v1/branches");
}

export async function newBranch(name) {
    console.log('web.newBranch', name);
    return await httpPost("/api/v1/branches", {name});
}

export async function deleteBranch(name) {
    console.log('web.deleteBranch', name);
    return await httpDelete("/api/v1/branches", {name});
}
