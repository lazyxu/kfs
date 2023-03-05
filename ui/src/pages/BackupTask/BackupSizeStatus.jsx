import {Alert, Box, Typography} from "@mui/material";
import humanize from "humanize";

export default function ({json}) {
    if (!json) {
        return (
            <Box/>
        );
    }
    if (json.errMsg) {
        return (
            <Alert variant="outlined" sx={{width: "max-content"}} severity="error">
                {json.errMsg}
            </Alert>
        );
    }
    return (
        <Alert variant="outlined" sx={{width: "max-content"}} severity={json.finished ? "success" : "info"}>
            <Typography>文件数量：{json.data.fileCount}</Typography>
            <Typography>目录数量：{json.data.dirCount}</Typography>
            <Typography>总大小：{humanize.filesize(json.data.fileSize)}</Typography>
        </Alert>
    );
}
