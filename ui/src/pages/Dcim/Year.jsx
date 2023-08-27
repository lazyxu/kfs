import Image from "./Image";
import { Box, Grid } from "@mui/material";
import { parseShotTime, timeSortFn } from "api/utils/api";

export default function ({ metadataList, chosenShotEquipment, chosenFileType }) {
    let filterMetadataList = metadataList
        .filter(metadata => chosenShotEquipment.includes(metadata.shotEquipment) && chosenFileType.includes(metadata.fileType.extension))
        .sort(timeSortFn);
    let dateMap = {};
    filterMetadataList.forEach(metadata => {
        let date = parseShotTime(metadata);
        date = date ? date.format("YYYY年") : "未知时间";
        if (dateMap.hasOwnProperty(date)) {
            dateMap[date].push(metadata);
        } else {
            dateMap[date] = [metadata];
        }
    })
    return (
        <Grid container spacing={1} style={{ overflowY: "scroll" }}>
            {Object.keys(dateMap).map(date => <Grid item xs={12} style={{ width: "100%" }} key={date}>
                <Box>{date}</Box>
                <Grid container spacing={1}>
                    {dateMap[date].sort(timeSortFn).map(metadata => {
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
