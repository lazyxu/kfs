import {httpGet} from "./common";
import axios from "axios";
import {getSysConfig} from "hox/sysConfig";

export async function list(driverName, filePath) {
    console.log('web.list', driverName, filePath);
    return await httpGet("/api/v1/list", {
        driverName,
        filePath,
    });
}

export async function openFile(driverName, filePath) {
    console.log('web.openFile', driverName, filePath);
    let resp = await axios.get(`${location.origin}/api/v1/openFile`, {
        params: {
            driverName,
            filePath,
            maxContentSize: getSysConfig().sysConfig.maxContentSize,
        },
        transformResponse: x=>x,
    });
    let tooLarge = resp.headers.get("Kfs-Too-Large");
    return {tooLarge, content: resp.data}
}

export async function download(driverName, filePath) {
    console.log('web.download', driverName, filePath);
    let resp = await axios.get(`${location.origin}/api/v1/downloadFile`, {
        params: {
            driverName,
            filePath,
        },
        responseType: "arraybuffer",
    });
    return resp.data;
}

export async function getImage(hash) {
    console.log('web.getImage', hash);
    let resp = await axios.get(`${location.origin}/api/v1/getImage`, {
        params: {
            hash
        },
        responseType: "arraybuffer",
    });
    return resp.data;
}
