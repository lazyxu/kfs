export async function listBranch() {
    console.log('web.listBranch');
    let resp = await fetch("http://127.0.0.1:1123/api/v1/branches");
    console.log(resp);
    if (resp.ok) {
        return await resp.json();
    }
    return [];
}
