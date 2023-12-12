const key = 'kfs-config';

function get() {
    const item = localStorage.getItem(key);
    try {
        return JSON.parse(item);
    } catch (_) {
        return undefined;
    }
}

function set(json) {
    localStorage.setItem(key, JSON.stringify(json, undefined, 2));
}

function remove() {
    localStorage.removeItem(key);
}

export default { get, set, remove };
