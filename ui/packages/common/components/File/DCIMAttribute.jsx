import { formatLinuxPermission, modeIsDir } from '@kfs/common/api/utils';
import { Close } from '@mui/icons-material';
import { Box, Dialog, DialogContent, DialogTitle, Grid } from "@mui/material";
import IconButton from "@mui/material/IconButton";
import humanize from 'humanize';

function Attr({ k, children }) {
    return <>
        <Grid size={{ xs: 4 }} sx={{ overflowWrap: "anywhere" }}>
            <Box fontWeight="bold">{k}：</Box>
        </Grid>
        <Grid size={{ xs: 8 }} sx={{ overflowWrap: "anywhere" }}>
            {children}
        </Grid>
    </>
}

export default ({ driverFile, onClose }) => {
    const isDir = modeIsDir(driverFile.mode);

    return (
        <Dialog open={true} fullWidth={true} onClose={onClose} maxWidth="sm">
            <DialogTitle>
                文件属性
                <IconButton
                    aria-label="close"
                    onClick={() => onClose()}
                    sx={{
                        position: 'absolute',
                        right: 8,
                        top: 8,
                    }}
                >
                    <Close />
                </IconButton>
            </DialogTitle>
            <DialogContent>
                <Box sx={{ flexGrow: 1, mt: 2 }}>
                    <Grid container spacing={1.5} sx={{ alignItems: "center" }}>
                        <Attr k="类型">{isDir ? "文件夹" : "文件"}</Attr>
                        {!isDir && <>
                            <Attr k="文件ID">{driverFile.file_id}</Attr>
                            <Attr k="文件大小">{humanize.filesize(driverFile.size)}</Attr>
                            <Attr k="md5">{driverFile.md5}</Attr>
                            <Attr k="sha256">{driverFile.sha256}</Attr>
                        </>}
                        <Attr k="文件权限">{formatLinuxPermission(driverFile.mode)}</Attr>
                        <Attr k="创建时间">{driverFile.create_time}</Attr>
                        <Attr k="属性修改时间">{driverFile.change_time}</Attr>
                        <Attr k="内容修改时间">{driverFile.modify_time}</Attr>
                        <Attr k="访问时间">{driverFile.access_time}</Attr>
                    </Grid>
                </Box>
            </DialogContent>
        </Dialog>
    );
};
