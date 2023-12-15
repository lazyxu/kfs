import { exifStatus } from '@kfs/mui/api/exif';
import { Box, Button } from "@mui/material";
import { useEffect, useState } from "react";

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
            if (status.cnt !== newStatus.cnt) {
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
            <Button variant="outlined" sx={{ width: "10em" }}
                onClick={e => {
                    onNewExif?.();
                }}
            >
                刷新
            </Button>
        </Box>
    );
}