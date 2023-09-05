import {httpGet, httpPost} from "./webServer";

export async function getDiskUsage() {
    console.log('web.diskUsage');
    return await httpGet("/api/v1/diskUsage");
}
