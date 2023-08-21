import Image from "./Image";
import { Box, Grid } from "@mui/material";
import { parseShotTime, timeSortFn } from "api/utils/api";

export default function ({ exifMap, chosenShotEquipment }) {
    let filterHashList = Object.keys(exifMap)
        .filter(hash => chosenShotEquipment.includes(exifMap[hash].shotEquipment))
        .sort((a, b) => timeSortFn(exifMap, a, b));
    let dateMap = {};
    filterHashList.forEach(hash => {
        let date = parseShotTime(exifMap[hash]);
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
                    {dateMap[date].sort((a, b) => timeSortFn(exifMap, a.hash, b.hash)).map(exif => {
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
