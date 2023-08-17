import { Box, Button, Checkbox, FormControl, FormControlLabel, FormGroup, Grid, Hidden, ImageList, ImageListItem, ImageListItemBar, InputLabel, MenuItem, Select, Stack, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { analysisExif, exifStatus } from 'api/web/exif';
import LinearProgressWithLabel from "pages/BackupTask/LinearProgressWithLabel";

export default function ({ onNewExif }) {
    const [status, setStatus] = useState({ analyzing: false });
    let interval;
    let exifStatusCb;
    exifStatusCb = newStatus => {
        if (newStatus.analyzing) {
            if (!interval) {
                interval = setInterval(() => {
                    exifStatus().then(exifStatusCb);
                }, 500);
            }
            if (status.cnt != newStatus.cnt) {
                onNewExif?.();
            }
        } else {
            clearTimeout(interval);
            interval = undefined;
            onNewExif?.();
        }
        setStatus(newStatus);
    }
    useEffect(() => {
        exifStatus().then(exifStatusCb);
    }, []);
    return (
        <Box sx={{ width: "100%" }}>
            {status.analyzing ?
                <>
                    <LinearProgressWithLabel variant="determinate" value={status.cnt / status.total * 100} />
                    <Box sx={{ paddingLeft: "0.5em" }}>
                        <Typography variant="body2" color="text.secondary">{`${status.cnt}/${status.total}`}</Typography>
                    </Box>
                    <Button variant="outlined" sx={{ width: "10em" }}
                        onClick={e => {
                            analysisExif(false);
                        }}
                    >
                        取消解析exif
                    </Button>
                </> :
                <>
                    <Button variant="outlined" sx={{ width: "10em" }}
                        onClick={e => {
                            analysisExif(true);
                            setTimeout(() => {
                                exifStatus().then(exifStatusCb);
                            }, 500);
                        }}
                    >
                        开始解析exif
                    </Button>
                    {status.finished && <Typography variant="body2" color="text.secondary">已完成 {status.cnt}</Typography>}
                </>}
        </Box>
    );
}