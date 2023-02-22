export async function httpGet(path) {
    let resp = await fetch(`http://127.0.0.1:1123${path}`);
    if (!resp.ok) {
        throw await resp.text();
    }
    let json = await resp.json();
    if (json.code !== 0) {
        throw await resp.msg;
    }
    return json.data;
}
