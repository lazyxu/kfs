import { Box, Grid } from "@mui/material";
import humanize from 'humanize';
import { useEffect, useState } from 'react';

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

export default ({ driverFile, filePath }) => {
    const [attributes, setAttributes] = useState(null);

    useEffect(() => {
        // getDriversDirCalculatedInfo(driverFile.id, filePath).then(setAttributes);
    }, [driverFile.id, filePath]);

    if (!attributes) {
        return <Box>计算中...</Box>;
    }

    return (
        <Grid container spacing={1.5} sx={{ alignItems: "center" }}>
            <Attr k="目录下文件">共 {attributes.fileCount} 个，总大小 {humanize.filesize(attributes.fileSize)}</Attr>
            <Attr k="目录下重复文件">共 {attributes.fileCount - attributes.distinctFileCount} 个，总大小 {humanize.filesize(attributes.fileSize - attributes.distinctFileSize)}</Attr>
            <Attr k="子目录">共 {attributes.dirCount} 个</Attr>
        </Grid>
    );
};
