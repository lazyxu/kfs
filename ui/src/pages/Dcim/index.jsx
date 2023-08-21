import { Box, Button, Checkbox, FormControl, FormControlLabel, FormGroup, Grid, Hidden, ImageList, ImageListItem, ImageListItemBar, InputLabel, MenuItem, Select, Stack } from "@mui/material";
// import FormControlLabel from '@mui/material/FormControlLabel';
import { useEffect, useState } from "react";
import { analysisExif, exifStatus, listExif } from 'api/web/exif';
import All from './All';
import Date from "./Date";
import Month from "./Month";
import Year from "./Year";
import Exif from "./Exif";

export default function ({ show }) {
    const [exifMap, setExifMap] = useState({});
    const [viewBy, setViewBy] = useState("所有照片");
    const [chosenModel, setChosenModel] = useState([]);
    const [ModelMap, setModelMap] = useState([]);
    return (
        <Stack style={{ width: "100%", height: "100%", display: show ? undefined : "none" }}>
            <Exif onNewExif={() => {
                listExif().then(exifMap => {
                    setExifMap(exifMap);
                    let ModelMap = {};
                    Object.values(exifMap).forEach(exif => {
                        if (ModelMap.hasOwnProperty(exif.Model)) {
                            ModelMap[exif.Model]++;
                        } else {
                            ModelMap[exif.Model] = 1;
                        }
                    })
                    setModelMap(ModelMap);
                });
            }} />
            <Box>
                <InputLabel sx={{display: "inline"}}>视图：</InputLabel>
                <Select
                    value={viewBy}
                    onChange={e => setViewBy(e.target.value)}
                    size="small"
                >
                    {["年", "月", "日", "所有照片"].map(value =>
                        <MenuItem key={value} value={value}>{value}</MenuItem>
                    )}
                </Select>
            </Box>
            <Box>
                <InputLabel sx={{display: "inline"}}>拍摄设备：</InputLabel>
                {Object.keys(ModelMap).map((Model, i) =>
                    <FormControlLabel key={i} control={
                        <Checkbox defaultChecked={chosenModel.includes(Model)} value={Model} onChange={e => {
                            setChosenModel(prev => {
                                let set = new Set(prev);
                                if (e.target.checked) {
                                    set.add(Model);
                                    return Array.from(set);
                                } else {
                                    set.delete(Model);
                                    return Array.from(set);
                                }
                            })
                        }} />
                    } label={(Model ? Model : "未知设备") + " (" + ModelMap[Model] + ")"} />
                )}
            </Box>
            {viewBy == "年" && <Year exifMap={exifMap} chosenModel={chosenModel} />}
            {viewBy == "月" && <Month exifMap={exifMap} chosenModel={chosenModel} />}
            {viewBy == "日" && <Date exifMap={exifMap} chosenModel={chosenModel} />}
            {viewBy == "所有照片" && <All exifMap={exifMap} chosenModel={chosenModel} />}
        </Stack>
    );
}
