import { Close } from "@mui/icons-material";
import { Dialog, DialogTitle, IconButton, Tab, Tabs, Typography } from "@mui/material";
import { useState } from 'react';
import NewBaiduPhoto from "./NewBaiduPhoto";
import NewLocalFile from "./NewLocalFile";
import NewNormalDriver from "./NewNormalDriver";

export default function ({ setOpen }) {
    let [driverType, setDriverType] = useState(0);
    return (
        <Dialog fullWidth={true} open={true} onClose={() => setOpen(false)}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary,
            }}>
                <Typography>新建云盘</Typography>
                <IconButton
                    aria-label="close"
                    onClick={() => setOpen(false)}
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
                {["普通云盘", "一刻相册备份盘", "本地文件备份盘"].map((v, i) =>
                    <Tab key={i} value={i} label={v} id={`simple-tab-${i}`} />
                )}
            </Tabs>
            {driverType === 0 && <NewNormalDriver setOpen={setOpen} />}
            {driverType === 1 && <NewBaiduPhoto setOpen={setOpen} />}
            {driverType === 2 && <NewLocalFile setOpen={setOpen} />}
        </Dialog>
    );
};
