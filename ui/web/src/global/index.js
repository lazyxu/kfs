Object.defineProperty(window, 'kfsEnv', {
    value: import.meta.env,
    enumerable: false,
    configurable: false,
    writable: false,
});


const KEY_KFS_CONFIG = 'kfs-config';
Object.defineProperty(window, "kfsConfig", {
    get() {
        const item = localStorage.getItem(KEY_KFS_CONFIG);
        try {
            return JSON.parse(item);
        } catch (_) {
            return undefined;
        }
    },
    set(json) {
        localStorage.setItem(KEY_KFS_CONFIG, JSON.stringify(json, undefined, 2));
    },
});
