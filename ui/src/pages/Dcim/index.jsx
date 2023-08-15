import useResourceManager from 'hox/resourceManager';
import { Box, Button, Checkbox, FormControl, FormControlLabel, FormGroup, Grid, ImageList, ImageListItem, ImageListItemBar, InputLabel, MenuItem, Select, Stack } from "@mui/material";
// import FormControlLabel from '@mui/material/FormControlLabel';
import { useEffect, useState } from "react";
import { analysisExif, listExif } from 'api/web/exif';
import moment from 'moment';

export default function ({ show }) {
    const [exifMap, setExifMap] = useState({});
    const [viewBy, setViewBy] = useState("所有照片");
    const [chosenHostComputer, setChosenHostComputer] = useState([]);
    const [hostComputerMap, setHostComputerMap] = useState([]);
    useEffect(() => {
        listExif().then(exifMap => {
            setExifMap(exifMap);
            let hostComputerMap = {};
            Object.values(exifMap).forEach(exif => {
                if (hostComputerMap.hasOwnProperty(exif.hostComputer)) {
                    hostComputerMap[exif.hostComputer]++;
                } else {
                    hostComputerMap[exif.hostComputer] = 1;
                }
            })
            setHostComputerMap(hostComputerMap);
            setChosenHostComputer(Object.keys(hostComputerMap));
        });
    }, []);
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
            <FormControl sx={{ minWidth: "10em" }}>
                <InputLabel id="view-by">视图</InputLabel>
                <Select
                    labelId="view-by"
                    value={viewBy}
                    onChange={e => setViewBy(e.target.value)}
                    sx={{ width: "10em" }}
                >
                    {["年", "月", "日", "所有照片"].map(value =>
                        <MenuItem key={value} value={value}>{value}</MenuItem>
                    )}
                </Select>
            </FormControl>
            <FormGroup>
                <InputLabel>拍摄设备</InputLabel>
                {Object.keys(hostComputerMap).map((hostComputer, i) =>
                    <FormControlLabel key={i} control={
                        <Checkbox defaultChecked={chosenHostComputer.includes(hostComputer)} value={hostComputer} onChange={e => {
                            setChosenHostComputer(prev => {
                                let set = new Set(prev);
                                if (e.target.checked) {
                                    set.add(hostComputer);
                                    return Array.from(set);
                                } else {
                                    set.delete(hostComputer);
                                    return Array.from(set);
                                }
                            })
                        }} />
                    } label={(hostComputer ? hostComputer : "未知设备") + " (" + hostComputerMap[hostComputer] + ")"} />
                )}
            </FormGroup>
            <Grid container spacing={2} overflow="scroll">
                {Object.keys(exifMap).sort((a, b) => exifMap[a].dateTime - exifMap[b].dateTime)
                    .filter(hash => chosenHostComputer.includes(exifMap[hash].hostComputer)).map(hash => {
                        let time = moment(exifMap[hash].dateTime / 1000 / 1000).format("YYYY年MM月DD日 HH:mm:ss");
                        return <Grid xs={2.4} key={hash}>
                            <Box sx={{width: "100%"}}>
                                <img style={{width: "100%"}}src={"http://127.0.0.1:1123/thumbnail?size=256&hash=" + hash} loading="lazy" title={time + "\n" + hash} />
                            </Box>
                        </Grid>
                    })}
            </Grid>
        </Stack>
    );
}
