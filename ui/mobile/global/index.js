import { MMKV } from 'react-native-mmkv';
import Toast from "react-native-toast-message";

let env = {
    VITE_APP_PLATFORM: "mobile",
    MODE: "development",
}

Object.defineProperty(window, 'kfsEnv', {
    value: env,
    enumerable: false,
    configurable: false,
    writable: false,
});

export const storage = new MMKV();

const KEY_KFS_CONFIG = 'kfs-config';
Object.defineProperty(window, "kfsConfig", {
    get() {
        const item = storage.getString(KEY_KFS_CONFIG);
        try {
            return JSON.parse(item);
        } catch (_) {
            return undefined;
        }
    },
    set(json) {
        storage.set(KEY_KFS_CONFIG, JSON.stringify(json, undefined, 2));
    },
});

window.noteError = (msg) => {
    Toast.show({
        type: 'error',
        text1: msg,
        text2: msg
    });
}

window.noteInfo = (msg) => {
    Toast.show({
        type: 'info',
        text1: msg
    });
}
