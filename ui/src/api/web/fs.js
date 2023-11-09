import axios from "axios";
import { getSysConfig } from "hox/sysConfig";
import { httpGet } from "./webServer";

export async function list(id, filePath) {
    console.log('web.list', id, filePath);
    return await httpGet("/api/v1/list", {
        id,
        filePath,
    });
}

export async function listDriverFileByHash(hash) {
    console.log('web.listDriverFileByHash', hash);
    return await httpGet("/api/v1/listDriverFileByHash", {
        hash
    });
}

export async function openFile(id, filePath) {
    console.log('web.openFile', id, filePath);
    let resp = await axios.get(`${getSysConfig().sysConfig.webServer}/api/v1/openFile`, {
        params: {
            id,
            filePath,
            maxContentSize: getSysConfig().sysConfig.maxContentSize,
        },
        transformResponse: x => x,
    });
    let tooLarge = resp.headers.get("Kfs-Too-Large");
    return { tooLarge, content: resp.data }
}

export async function download(id, filePath) {
    console.log('web.download', id, filePath);
    let resp = await axios.get(`${getSysConfig().sysConfig.webServer}/api/v1/downloadFile`, {
        params: {
            id,
            filePath,
        },
        responseType: "arraybuffer",
    });
    return resp.data;
}

export async function getImage(hash) {
    console.log('web.getImage', hash);
    let resp = await axios.get(`${getSysConfig().sysConfig.webServer}/api/v1/getImage`, {
        params: {
            hash
        },
        responseType: "arraybuffer",
    });
    return resp.data;
}
