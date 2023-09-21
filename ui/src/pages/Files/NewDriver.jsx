import { Dialog, DialogTitle, Tab, Tabs, Typography } from "@mui/material";
import { useState } from 'react';
import NewBaiduPhoto from "./NewBaiduPhoto";
import NewNormalDriver from "./NewNormalDriver";

export default function ({ setOpen }) {
    let [driverType, setDriverType] = useState(0);
    return (
        <Dialog open={true} onClose={() => setOpen(false)} fullScreen={true}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary,
            }}>
                <Typography>新建云盘</Typography>
            </DialogTitle>
            <Tabs value={driverType} onChange={(e, v) => setDriverType(v)} sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary,
                borderBottom: 1, borderColor: 'divider'
            }}
            >
                {["普通云盘", "一刻相册"].map((v, i) =>
                    <Tab key={i} value={i} label={v} id={`simple-tab-${i}`} />
                )}
            </Tabs>
            {driverType === 0 && <NewNormalDriver setOpen={setOpen} />}
            {driverType === 1 && <NewBaiduPhoto setOpen={setOpen} />}
        </Dialog>
    );
};
