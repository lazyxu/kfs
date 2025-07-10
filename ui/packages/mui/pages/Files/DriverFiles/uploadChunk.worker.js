self.importScripts("https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"); // 确保axios可用，或用fetch

self.onmessage = async function (e) {
    const { driverId, dirPath, hash, chunkIndex, chunk, fileMeta } = e.data;
    try {
        const formData = new FormData();
        formData.append("driverId", driverId);
        formData.append("dirPath", dirPath);
        formData.append("hash", hash);
        formData.append("chunkIndex", chunkIndex);
        formData.append("chunk", chunk);
        if (fileMeta) {
            formData.append("fileName", fileMeta.name);
            formData.append("fileSize", fileMeta.size);
            formData.append("fileLastModified", fileMeta.lastModified);
            formData.append("fileType", fileMeta.type);
        }
        await fetch("/api/v1/uploadFileChunk", {
            method: "POST",
            body: formData,
        });
        self.postMessage({ type: "success" });
    } catch (err) {
        self.postMessage({ type: "error", error: err.message });
    }
};
