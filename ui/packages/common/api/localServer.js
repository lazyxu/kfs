import { getSysConfig } from "@kfs/common/hox/sysConfig";
import axios from "axios";

export async function httpGet(path, params) {
    let resp = await axios.get(`http://127.0.0.1:${getSysConfig().port}${path}`, {params});
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
    let resp = await axios.post(`http://127.0.0.1:${getSysConfig().port}${path}`, null, {params});
    let json = resp.data;
    if (!json) {
        return;
    }
    if (json.code !== 0) {
        throw await json.msg;
    }
    return json.data;
}

export async function httpPostBody(path, params) {
    let resp = await axios.post(`http://127.0.0.1:${getSysConfig().port}${path}`, params);
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
    let resp = await axios.delete(`http://127.0.0.1:${getSysConfig().port}${path}`, {params});
    let json = resp.data;
    if (!json) {
        return;
    }
    if (json.code !== 0) {
        throw await json.msg;
    }
    return json.data;
}
