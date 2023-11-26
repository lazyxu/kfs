import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Grid, Tab, Tabs } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import moment from "moment/moment";
import { useState } from 'react';
import DriverAttributeNormal from './DriverAttributeNormal';
import DriverBaiduPhoto from './DriverBaiduPhoto';
import DriverLocalFile from './DriverLocalFile';
import DriverLocalFileFilter from './DriverLocalFileFilter';

function formatTime(t) {
    return moment(t / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
}

function Attr({ k, children }) {
    return <>
        <Grid xs={4} item sx={{ overflowWrap: "anywhere" }}><Box>{k}：</Box></Grid>
        <Grid xs={8} item sx={{ overflowWrap: "anywhere" }}>{children}</Grid>
    </>
}

function getDriverType(driver) {
    switch (driver.type) {
        case "baiduPhoto":
            return "一刻相册备份盘";
        case "localFile":
            return "本地文件备份盘";
        case "":
            return "普通云盘";
        default:
            break;
    }
}

export default ({ setOpen, driver }) => {
    const [attributeType, setAttributeType] = useState(0);
    return (
        <Dialog open={true} fullWidth={true} onClose={() => setOpen(false)}>
            <DialogTitle sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.secondary
            }}>
                云盘属性
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
            <DialogContent sx={{
                backgroundColor: theme => theme.background.primary,
                color: theme => theme.context.primary
            }}>
                <Tabs value={attributeType} onChange={(e, v) => setAttributeType(v)} sx={{
                    backgroundColor: theme => theme.background.primary,
                    color: theme => theme.context.secondary,
                    borderBottom: 1, borderColor: 'divider',
                    marginBottom: "1em"
                }}
                >
                    {driver.type === "" && <Tab key={0} value={0} label="常规" />}
                    {driver.type === "baiduPhoto" && [
                        <Tab key={0} value={0} label="常规" />,
                        <Tab key={1} value={1} label="同步" />,
                    ]}
                    {driver.type === "localFile" && [
                        <Tab key={0} value={0} label="常规" />,
                        <Tab key={1} value={1} label="同步" />,
                        <Tab key={2} value={2} label="过滤规则设置" />,
                    ]}
                </Tabs>
                {attributeType === 0 && <DriverAttributeNormal setOpen={setOpen} driver={driver} />}
                {attributeType === 1 && (driver.type === "baiduPhoto" ? <DriverBaiduPhoto driver={driver} /> : <DriverLocalFile driver={driver} />)}
                {attributeType === 2 && <DriverLocalFileFilter driver={driver} />}
            </DialogContent>
        </Dialog>
    );
};
