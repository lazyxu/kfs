import axios from "axios";
import { getSysConfig } from "hox/sysConfig";

export async function httpGet(path, params) {
    let resp = await axios.get(`${getSysConfig().sysConfig.webServer}${path}`, {params});
    let json = resp.data;
    if (!json) {
        return;
    }
    if (json.code !== 0) {
        throw await json.msg;
    }
    return json.data;
}

export async function httpPost(path, params) {
    let resp = await axios.post(`${getSysConfig().sysConfig.webServer}${path}`, null, {params});
    let json = resp.data;
    if (!json) {
        return;
    }
    if (json.code !== 0) {
        throw await json.msg;
    }
    return json.data;
}

export async function httpDelete(path, params) {
    let resp = await axios.delete(`${getSysConfig().sysConfig.webServer}${path}`, {params});
    let json = resp.data;
    if (!json) {
        return;
    }
    if (json.code !== 0) {
        throw await json.msg;
    }
    return json.data;
}
