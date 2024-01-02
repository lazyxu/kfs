import AsyncStorage from '@react-native-async-storage/async-storage';
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

const KEY_KFS_CONFIG = 'kfs-config';

window.getKfsConfig = async () => {
    try {
        const item = await AsyncStorage.getItem(KEY_KFS_CONFIG);
        return JSON.parse(item);
    } catch (_) {
        return undefined;
    }
};

window.setKfsConfig = async (json) => {
    try {
        await AsyncStorage.setItem(KEY_KFS_CONFIG, JSON.stringify(json, undefined, 2));
    } catch (e) {
        noteError(e.message);
    }
};

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
