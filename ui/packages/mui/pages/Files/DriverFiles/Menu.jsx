import { checkFileExists, createFolder, mergeFileChunks } from "@kfs/mui/api/fs";
import Menu from "@kfs/mui/components/Menu";
import { noteError, noteSuccess } from "@kfs/mui/components/Notification";
import { Add, Bookmark, ContentPaste, CreateNewFolder, DriveFolderUpload, History, IosShare, Settings, UploadFile } from "@mui/icons-material";
import { Box, LinearProgress, ListItemText, MenuItem, Typography } from "@mui/material";
import { useRef, useState } from "react";
import NewBaiduPhoto from "./CreateBaiduPhoto";
import Sha256Worker from "./sha256.worker.js?worker";
import UploadChunkWorker from "./uploadChunk.worker.js?worker";
import CreateFolder from "./CreateFolder";

const CHUNK_SIZE = 100 * 1024 * 1024; // 100MB

export default function ({ contextMenu, setContextMenu, onUploaded }) {
    const fileInputRef = useRef();
    const [progress, setProgress] = useState(null);
    const [createFolderDialog, setCreateFolderDialog] = useState(false);

    console.log("Menu contextMenu", contextMenu);
    const onClose = () => setContextMenu(c => ({ ...c, open: false }))
    const handleUploadClick = () => {
        onClose();
        fileInputRef.current?.click();
    };

    const handleFileChange = async (e) => {
        const file = e.target.files[0];
        if (!file) return;
        console.log("handleFileChange", file, contextMenu);
        setProgress({ percent: 0, status: "计算SHA256..." });

        // 1. 计算SHA256
        const {sha256, md5} = await calcFileSha256(file, setProgress);
        setProgress({ percent: 0, status: "检查服务器..." });

        // 2. 检查服务器是否已存在
        const exists = await checkFileExists(contextMenu.dirPath, file.name, sha256, md5);
        if (exists) {
            setProgress(null);
            noteSuccess("文件已存在，无需上传");
            // setContextMenu(null);
            onUploaded?.();
            return;
        }

        // 3. 分块上传
        setProgress({ percent: 0, status: "上传中..." });
        try {
            await uploadLargeFile({
                file,
                dirPath: contextMenu.dirPath,
                md5, sha256,
                setProgress,
                fileMeta: {
                    name: file.name,
                    size: file.size,
                    lastModified: file.lastModified,
                    type: file.type,
                },
            });
            setProgress(null);
            noteSuccess("上传成功");
            setContextMenu(null);
            onUploaded?.();
        } catch (err) {
            setProgress(null);
            noteError("上传失败：" + (err?.message || err));
        }
        e.target.value = "";
    };
    return (
        <>
            <input
                type="file"
                style={{ display: "none" }}
                ref={fileInputRef}
                onChange={handleFileChange}
            />
            <Menu
                contextMenu={contextMenu}
                open={contextMenu?.open}
                onClose={onClose}
            >
                <MenuItem onClick={handleUploadClick}>
                    <UploadFile />
                    <ListItemText>上传文件</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <DriveFolderUpload />
                    <ListItemText>上传文件夹</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <Add />
                    <ListItemText>新建文件</ListItemText>
                </MenuItem>
                <MenuItem onClick={async () => {
                    onClose();
                    setCreateFolderDialog(true);
                }}>
                    <CreateNewFolder />
                    <ListItemText>新建文件夹</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <ContentPaste />
                    <ListItemText>粘贴</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <Settings />
                    <ListItemText>属性</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <IosShare />
                    <ListItemText>分享</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <History />
                    <ListItemText>历史版本</ListItemText>
                </MenuItem>
                <MenuItem disabled>
                    <Bookmark />
                    <ListItemText>添加书签</ListItemText>
                </MenuItem>
            </Menu>
            {createFolderDialog && <CreateFolder onClose={()=> {
                setCreateFolderDialog(false);
            }} />}
            {progress && (
                <Box sx={{ width: 300, p: 2, position: "fixed", bottom: 40, left: "50%", transform: "translateX(-50%)", zIndex: 9999, bgcolor: "#fff", borderRadius: 2, boxShadow: 3 }}>
                    <Typography variant="body2">{progress.status}</Typography>
                    <LinearProgress variant="determinate" value={progress.percent} sx={{ mt: 1 }} />
                    <Typography variant="caption">{progress.percent.toFixed(1)}%</Typography>
                </Box>
            )}
        </>
    );
}

// 计算文件SHA256（worker）
function calcFileSha256(file, setProgress) {
    return new Promise((resolve, reject) => {
        const worker = new Sha256Worker();
        worker.postMessage({ file });
        worker.onmessage = (e) => {
            if (e.data.type === "progress") {
                setProgress({ percent: e.data.percent, status: "计算SHA256..." });
            } else if (e.data.type === "result") {
                resolve(e.data);
                worker.terminate();
            } else if (e.data.type === "error") {
                reject(e.data.error);
                worker.terminate();
            }
        };
    });
}

// 分块并发上传
async function uploadLargeFile({ file, dirPath, md5, sha256, setProgress, fileMeta }) {
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
    let uploaded = 0;
    let concurrency = navigator.hardwareConcurrency || 4;
    let uploading = 0;
    let nextChunk = 0;
    let errors = [];
    const chunkStatus = Array(totalChunks).fill(false);

    return new Promise((resolve, reject) => {
        function uploadNext() {
            if (errors.length) return;
            if (uploaded === totalChunks) {
                // 合并分块
                mergeFileChunks(dirPath, md5, sha256, totalChunks, fileMeta)
                    .then(resolve)
                    .catch(reject);
                return;
            }
            while (uploading < concurrency && nextChunk < totalChunks) {
                const chunkIndex = nextChunk++;
                uploading++;
                const start = chunkIndex * CHUNK_SIZE;
                const end = Math.min(file.size, start + CHUNK_SIZE);
                const chunk = file.slice(start, end);

                // 用worker上传分块
                const worker = new UploadChunkWorker();
                worker.postMessage({
                    dirPath,
                    md5, sha256,
                    chunkIndex,
                    chunk,
                    fileMeta, // 传递文件元信息
                });
                worker.onmessage = (e) => {
                    uploading--;
                    if (e.data.type === "success") {
                        chunkStatus[chunkIndex] = true;
                        uploaded++;
                        setProgress({
                            percent: (uploaded / totalChunks) * 100,
                            status: `上传中...（${uploaded}/${totalChunks}）`
                        });
                        uploadNext();
                    } else if (e.data.type === "error") {
                        errors.push(e.data.error);
                        reject(e.data.error);
                        worker.terminate();
                    }
                    worker.terminate();
                };
            }
        }
        uploadNext();
    });
}
