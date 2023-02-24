import {httpGet} from "./common";
import axios from "axios";
import {getSysConfig} from "hox/sysConfig";

export async function list(branchName, filePath) {
    console.log('web.list', branchName, filePath);
    return await httpGet("/api/v1/list", {
        branchName,
        filePath: filePath.join("/"),
    });
}

export async function openFile(branchName, filePath) {
    console.log('web.openFile', branchName, filePath);
    let resp = await axios.get(`http://127.0.0.1:1123/api/v1/openFile`, {
        params: {
            branchName,
            filePath: filePath.join("/"),
            maxContentSize: getSysConfig().sysConfig.maxContentSize,
        },
        transformResponse: x=>x,
    });
    let tooLarge = resp.headers.get("Kfs-Too-Large");
    return {tooLarge, content: resp.data}
}

export async function download(branchName, filePath) {
    console.log('web.download', branchName, filePath);
    let resp = await axios.get(`http://127.0.0.1:1123/api/v1/downloadFile`, {
        params: {
            branchName,
            filePath: filePath.join("/"),
        },
        responseType: "arraybuffer",
    });
    return resp.data;
}
