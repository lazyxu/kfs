import { Box, Button, Checkbox, FormControl, FormControlLabel, FormGroup, Grid, Hidden, ImageList, ImageListItem, ImageListItemBar, InputLabel, MenuItem, Select, Stack } from "@mui/material";
import moment from 'moment';

export default function ({ exifMap, chosenHostComputer }) {
    return (
        <Grid container spacing={1} style={{ overflowY: "scroll" }}>
            {Object.keys(exifMap).sort((a, b) => exifMap[a].dateTime - exifMap[b].dateTime)
                .filter(hash => chosenHostComputer.includes(exifMap[hash].hostComputer)).map(hash => {
                    let time = moment(exifMap[hash].dateTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
                    return <Grid item style={{ width: "256px", height: "256px" }} key={hash}>
                        <Box sx={{ width: "100%" }}>
                            <img style={{ width: "100%" }} src={"http://127.0.0.1:1123/thumbnail?size=256&cutSquare=true&hash=" + hash} loading="lazy" title={time + "\n" + hash} />
                        </Box>
                    </Grid>
                })}
        </Grid>
    );
}
