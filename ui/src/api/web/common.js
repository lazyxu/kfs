import axios from "axios";

export async function httpGet(path, params) {
    let resp = await axios.get(`${location.origin}${path}`, {params});
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
    let resp = await axios.post(`${location.origin}${path}`, null, {params});
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
    let resp = await axios.delete(`${location.origin}${path}`, {params});
    let json = resp.data;
    if (!json) {
        return;
    }
    if (json.code !== 0) {
        throw await json.msg;
    }
    return json.data;
}
