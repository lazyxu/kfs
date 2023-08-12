import {httpGet} from "./common";
import axios from "axios";
import {getSysConfig} from "hox/sysConfig";

export async function list(driverName, filePath) {
    console.log('web.list', driverName, filePath);
    return await httpGet("/api/v1/list", {
        driverName,
        filePath: filePath,
    });
}

export async function openFile(driverName, filePath) {
    console.log('web.openFile', driverName, filePath);
    let resp = await axios.get(`http://127.0.0.1:1123/api/v1/openFile`, {
        params: {
            driverName,
            filePath: filePath.join("/"),
            maxContentSize: getSysConfig().sysConfig.maxContentSize,
        },
        transformResponse: x=>x,
    });
    let tooLarge = resp.headers.get("Kfs-Too-Large");
    return {tooLarge, content: resp.data}
}

export async function download(driverName, filePath) {
    console.log('web.download', driverName, filePath);
    let resp = await axios.get(`http://127.0.0.1:1123/api/v1/downloadFile`, {
        params: {
            driverName,
            filePath: filePath.join("/"),
        },
        responseType: "arraybuffer",
    });
    return resp.data;
}
