import Image from "./Image";
import { Box, Grid } from "@mui/material";
import { parseDateTime } from "api/utils/api";

export default function ({ exifMap, chosenModel }) {
    let filterHashList = Object.keys(exifMap).sort((a, b) => exifMap[a].DateTime - exifMap[b].DateTime)
        .filter(hash => chosenModel.includes(exifMap[hash].Model));
    let dateMap = {};
    filterHashList.forEach(hash => {
        let date = parseDateTime(exifMap[hash]);
        date = date ? date.format("YYYY年MM月DD日") : "未知时间";
        let elm = { hash, ...exifMap[hash] };
        if (dateMap.hasOwnProperty(date)) {
            dateMap[date].push(elm);
        } else {
            dateMap[date] = [elm];
        }
    })
    return (
        <Grid container spacing={1} style={{ overflowY: "scroll" }}>
            {Object.keys(dateMap).map(date => <Grid item xs={12} style={{ width: "100%" }} key={date}>
                <Box>{date}</Box>
                <Grid container spacing={1}>
                    {dateMap[date].map(exif => {
                        return <Grid item style={{ width: "256px", height: "256px" }} key={exif.hash}>
                            <Image hash={exif.hash} exif={exif} />
                        </Grid>
                    })}
                </Grid>
            </Grid>
            )}
        </Grid>
    );
}
