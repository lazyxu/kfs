import { parseShotTime } from "@kfs/api";
import { Box, Grid } from "@mui/material";
import Image from "./Image";

export default function ({ metadataList }) {
    let dateMap = {};
    metadataList.forEach(metadata => {
        let date = parseShotTime(metadata);
        date = date ? date.format("YYYY年MM月") : "未知时间";
        if (dateMap.hasOwnProperty(date)) {
            dateMap[date].push(metadata);
        } else {
            dateMap[date] = [metadata];
        }
    })
    return (
        <Grid container padding={1} spacing={1}>
            {Object.keys(dateMap).map(date => <Grid item xs={12} style={{ width: "100%" }} key={date}>
                <Box>{date}</Box>
                <Grid container padding={1} spacing={1}>
                    {dateMap[date].map(metadata => {
                        return <Grid item style={{ width: "256px", height: "256px" }} key={metadata.hash}>
                            <Image metadata={metadata} />
                        </Grid>
                    })}
                </Grid>
            </Grid>
            )}
        </Grid>
    );
}
