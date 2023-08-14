import useResourceManager from 'hox/resourceManager';
import { Button, Stack } from "@mui/material";
import { useEffect } from "react";
import { analysisExif } from 'api/web/exif';

export default function ({ show }) {
    return (
        <Stack style={{ width: "100%", height: "100%", display: show ? undefined : "none" }}>
            <Button variant="outlined" sx={{ width: "10em" }}
                onClick={e => {
                    analysisExif(true);
                }}
            >
                开始解析exif
            </Button>
            <Button variant="outlined" sx={{ width: "10em" }}
                onClick={e => {
                    analysisExif(false);
                }}
            >
                结束解析exif
            </Button>
        </Stack>
    );
}
