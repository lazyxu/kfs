import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Grid } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import moment from "moment/moment";
import DriverBaiduPhoto from './DriverBaiduPhoto';

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
            return "一刻相册";
        case "":
            return "普通云盘";
        default:
            break;
    }
}

export default ({ setOpen, driver }) => {
    // TODO: get more calculated attributes from server.
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
                <Grid container spacing={1.5} sx={{ alignItems: "center" }}>
                    <Attr k="云盘名称">{driver.name}</Attr>
                    <Attr k="描述">{driver.description}</Attr>
                    <Attr k="云盘类型">{getDriverType(driver)}</Attr>
                    {/* <Box variant="body"> */}
                    {/* <Box>文件总数：{driver.count}</Box> */}
                    {/* <Box>总大小：{humanize.filesize(driver.size)}</Box> */}
                    {/* <Typography>可修改该云盘的设备：any</Typography> */}
                    {/* <Typography>可读取该云盘的设备：any</Typography> */}
                    {/* </Box> */}
                    {driver.type === "baiduPhoto" && <>
                        <Attr k="同步"><DriverBaiduPhoto driver={driver} /></Attr>
                        <Attr k="accessToken">{driver.accessToken}</Attr>
                        <Attr k="refreshToken">{driver.refreshToken}</Attr>
                    </>}
                </Grid>
            </DialogContent>
        </Dialog>
    );
};
