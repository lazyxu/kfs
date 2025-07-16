self.onmessage = async function (e) {
    const { dirPath, md5, sha256, chunkIndex, chunk, fileMeta } = e.data;
    try {
        const formData = new FormData();
        formData.append("dirPath", dirPath);
        formData.append("md5", md5);
        formData.append("sha256", sha256);
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
