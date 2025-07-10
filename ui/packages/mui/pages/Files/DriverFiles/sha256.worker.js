self.importScripts('https://cdn.jsdelivr.net/npm/js-sha256@0.9.0/build/sha256.min.js');

self.onmessage = async function (e) {
    const { file } = e.data;
    try {
        const chunkSize = 10 * 1024 * 1024; // 10MB
        const chunks = Math.ceil(file.size / chunkSize);
        let sha256 = self.sha256.create();
        let offset = 0;
        for (let i = 0; i < chunks; i++) {
            const blob = file.slice(offset, offset + chunkSize);
            const buf = await blob.arrayBuffer();
            sha256.update(new Uint8Array(buf));
            offset += chunkSize;
            self.postMessage({ type: "progress", percent: ((i + 1) / chunks) * 100 });
        }
        self.postMessage({ type: "result", hash: sha256.hex() });
    } catch (err) {
        self.postMessage({ type: "error", error: err.message });
    }
};
