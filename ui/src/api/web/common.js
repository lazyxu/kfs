import axios from "axios";

export async function httpGet(path, params) {
    let resp = await axios.get(`http://127.0.0.1:1123${path}`, { params });
    let json = resp.data;
    if (json.code !== 0) {
        throw await json.msg;
    }
    return json.data;
}

