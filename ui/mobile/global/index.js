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
    webServer: "http://192.168.2.12:1123",
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

window.noteError = () => {
    Toast.show({
        type: 'success',
        text1: 'Hello',
        text2: 'This is some something 👋'
    });
}
