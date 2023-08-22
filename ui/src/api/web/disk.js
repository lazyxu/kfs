import {httpGet, httpPost} from "./common";

export async function getDiskUsage() {
    console.log('web.diskUsage');
    return await httpGet("/api/v1/diskUsage");
}
