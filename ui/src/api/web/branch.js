import {httpGet} from "./common";

export async function listBranch() {
    console.log('web.listBranch');
    return await httpGet("/api/v1/branches");
}
