import { toPrecent } from '@kfs/api';
import { Box, InputLabel, Stack, Typography } from '@mui/material';
import { getDiskUsage } from 'api/disk';
import humanize from 'humanize';
import { useEffect, useState } from 'react';

export default () => {
    const [diskUsage, setDiskUsage] = useState();
    useEffect(() => {
        getDiskUsage().then(setDiskUsage);
    }, []);
    return (
        <Stack style={{ padding: "1em", overflowY: 'auto' }}>
            {!diskUsage ? <Typography>正在获取...</Typography> :
                <Box>
                    <Box sx={{ width: "100%", height: "1em", paddingBottom: "1em" }}>
                        <Box sx={{ borderRadius: "5px 0 0 5px", display: "inline-block", height: "1em", width: toPrecent((diskUsage.total - diskUsage.file - diskUsage.metadata - diskUsage.thumbnail - diskUsage.transCode - diskUsage.free) / diskUsage.total), background: "gray" }} />
                        <Box sx={{ display: "inline-block", height: "1em", width: toPrecent((diskUsage.file) / diskUsage.total), background: "green" }} />
                        <Box sx={{ display: "inline-block", height: "1em", width: toPrecent((diskUsage.metadata) / diskUsage.total), background: "blueviolet" }} />
                        <Box sx={{ display: "inline-block", height: "1em", width: toPrecent((diskUsage.thumbnail) / diskUsage.total), background: "orange" }} />
                        <Box sx={{ display: "inline-block", height: "1em", width: toPrecent((diskUsage.transCode) / diskUsage.total), background: "yellowgreen" }} />
                        <Box sx={{ borderRadius: "0 5px 5px 0", display: "inline-block", height: "1em", width: toPrecent((diskUsage.free) / diskUsage.total), background: "cornflowerblue" }} />
                    </Box>
                    <Box>
                        <InputLabel>硬盘总空间：{humanize.filesize(diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "gray", height: "1em", width: "1em", borderRadius: "1em" }} /> 其他：{humanize.filesize(diskUsage.total - diskUsage.file - diskUsage.metadata - diskUsage.thumbnail - diskUsage.transCode - diskUsage.free)} {toPrecent((diskUsage.total - diskUsage.file - diskUsage.metadata - diskUsage.thumbnail - diskUsage.free) / diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "green", height: "1em", width: "1em", borderRadius: "1em" }} /> 云盘文件：{humanize.filesize(diskUsage.file)} {toPrecent(diskUsage.file / diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "blueviolet", height: "1em", width: "1em", borderRadius: "1em" }} /> 元数据：{humanize.filesize(diskUsage.metadata)} {toPrecent(diskUsage.metadata / diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "orange", height: "1em", width: "1em", borderRadius: "1em" }} /> 缩略图：{humanize.filesize(diskUsage.thumbnail)} {toPrecent(diskUsage.thumbnail / diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "yellowgreen", height: "1em", width: "1em", borderRadius: "1em" }} /> 图片视频转码：{humanize.filesize(diskUsage.transCode)} {toPrecent(diskUsage.transCode / diskUsage.total)}</InputLabel>
                        <InputLabel><Box style={{ display: "inline-block", background: "cornflowerblue", height: "1em", width: "1em", borderRadius: "1em" }} /> 剩余：{humanize.filesize(diskUsage.free)} {toPrecent(diskUsage.free / diskUsage.total)}</InputLabel>
                    </Box>
                </Box>
            }
        </Stack>
    );
};
