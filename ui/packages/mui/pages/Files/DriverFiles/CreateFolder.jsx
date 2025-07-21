import { Close } from "@mui/icons-material";
import { Dialog, DialogTitle, IconButton, Tab, Tabs, Typography } from "@mui/material";
import { useState } from 'react';
import CreateBaiduPhoto from "./CreateBaiduPhoto";
import CreateNormalFolder from "./CreateNormalFolder";

export default function ({ onClose }) {
    let [folderType, setFolderType] = useState(0);
    console.log("window.kfsEnv", window.kfsEnv)
    return (
        <Dialog fullWidth={true} open={true} onClose={onClose}>
            <DialogTitle>
                <Typography>新建文件夹</Typography>
                <IconButton
                    aria-label="close"
                    onClick={() => onClose()}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                        color: (theme) => theme.palette.grey[500],
                    }}
                >
                    <Close />
                </IconButton>
            </DialogTitle>
            <Tabs
                sx={{borderBottom: 1, borderColor: 'divider' }}
                value={folderType} onChange={(e, v) => setFolderType(v)}
            >
                {(window.kfsEnv.VITE_APP_PLATFORM !== 'web' ?
                    ["普通文件夹", "来自 一刻相册", "同步本地文件夹"] :
                    ["普通文件夹", "来自 一刻相册"]).map((v, i) =>
                        <Tab key={i} value={i} label={v} id={i} />
                    )}
            </Tabs>
            {folderType === 0 && <CreateNormalFolder onClose={onClose} />}
            {folderType === 1 && <CreateBaiduPhoto onClose={onClose} />}
            {/* {folderType === 2 && <CreateLocalFileFolder onClose={onClose} />} */}
        </Dialog>
    );
};
