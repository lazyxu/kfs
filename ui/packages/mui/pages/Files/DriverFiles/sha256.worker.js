import { md5 } from 'js-md5';
import { sha256 } from 'js-sha256';

self.onmessage = async function (e) {
    const { file } = e.data;
    try {
        const chunkSize = 10 * 1024 * 1024; // 10MB
        const chunks = Math.ceil(file.size / chunkSize);
        let sha256Hasher = sha256.create();
        let md5Hasher = md5.create();
        let offset = 0;
        for (let i = 0; i < chunks; i++) {
            const blob = file.slice(offset, offset + chunkSize);
            const buf = await blob.arrayBuffer();
            sha256Hasher.update(new Uint8Array(buf));
            md5Hasher.update(new Uint8Array(buf));
            offset += chunkSize;
            self.postMessage({ type: "progress", percent: ((i + 1) / chunks) * 100 });
        }
        self.postMessage({ type: "result", sha256: sha256Hasher.hex(), md5: md5Hasher.hex() });
    } catch (err) {
        self.postMessage({ type: "error", error: err.message });
    }
};
