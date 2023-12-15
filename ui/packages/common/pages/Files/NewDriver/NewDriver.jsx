import { Close } from "@mui/icons-material";
import { Dialog, DialogTitle, IconButton, Tab, Tabs, Typography } from "@mui/material";
import { useState } from 'react';
import NewBaiduPhoto from "./NewBaiduPhoto";
import NewLocalFileDriver from "./NewLocalFileDriver";
import NewNormalDriver from "./NewNormalDriver";

export default function ({ onClose, onSucc }) {
    let [driverType, setDriverType] = useState(0);
    console.log("window.kfsEnv", window.kfsEnv)
    return (
        <Dialog fullWidth={true} open={true} onClose={onClose}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary,
            }}>
                <Typography>新建云盘</Typography>
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
            <Tabs value={driverType} onChange={(e, v) => setDriverType(v)} sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary,
                borderBottom: 1, borderColor: 'divider'
            }}
            >
                {(window.kfsEnv.VITE_APP_PLATFORM !== 'web' ?
                    ["普通云盘", "一刻相册备份盘", "本地文件备份盘"] :
                    ["普通云盘", "一刻相册备份盘"]).map((v, i) =>
                        <Tab key={i} value={i} label={v} id={`simple-tab-${i}`} />
                    )}
            </Tabs>
            {driverType === 0 && <NewNormalDriver onSucc={onSucc} />}
            {driverType === 1 && <NewBaiduPhoto onSucc={onSucc} />}
            {driverType === 2 && <NewLocalFileDriver onSucc={onSucc} />}
        </Dialog>
    );
};
