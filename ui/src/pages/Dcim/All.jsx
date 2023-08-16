import { Box, Button, Checkbox, FormControl, FormControlLabel, FormGroup, Grid, Hidden, ImageList, ImageListItem, ImageListItemBar, InputLabel, MenuItem, Select, Stack } from "@mui/material";
import Image from "./Image";

export default function ({ exifMap, chosenHostComputer }) {
    return (
        <Grid container spacing={1} style={{ overflowY: "scroll" }}>
            {Object.keys(exifMap).sort((a, b) => exifMap[a].dateTime - exifMap[b].dateTime)
                .filter(hash => chosenHostComputer.includes(exifMap[hash].hostComputer)).map(hash => {
                    return <Grid item style={{ width: "256px", height: "256px" }} key={hash}>
                        <Image hash={hash} exif={exifMap[hash]} />
                    </Grid>
                })}
        </Grid>
    );
}
