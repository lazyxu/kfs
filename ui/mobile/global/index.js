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

let kfsConfig = {
    webServer: "http://127.0.0.1:1123",
};

const KEY_KFS_CONFIG = 'kfs-config';
Object.defineProperty(window, "kfsConfig", {
    get() {
        return kfsConfig;
    },
    set(json) {
        kfsConfig = json;
    },
});

window.noteError = (msg) => {
    Toast.show({
        type: 'error',
        text1: msg,
        text2: msg
    });
}