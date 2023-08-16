import Image from "./Image";
import { Box, Button, Checkbox, FormControl, FormControlLabel, FormGroup, Grid, Hidden, ImageList, ImageListItem, ImageListItemBar, InputLabel, MenuItem, Select, Stack } from "@mui/material";
import moment from 'moment';

export default function ({ exifMap, chosenHostComputer }) {
    let filterHashList = Object.keys(exifMap).sort((a, b) => exifMap[a].dateTime - exifMap[b].dateTime)
        .filter(hash => chosenHostComputer.includes(exifMap[hash].hostComputer));
    let dateMap = {};
    filterHashList.forEach(hash => {
        let date = moment(exifMap[hash].dateTime / 1000 / 1000).format("YYYY年");
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
