import { Box, Button, Checkbox, FormControl, FormControlLabel, FormGroup, Grid, Hidden, ImageList, ImageListItem, ImageListItemBar, InputLabel, MenuItem, Select, Stack } from "@mui/material";
// import FormControlLabel from '@mui/material/FormControlLabel';
import { useEffect, useState } from "react";
import { analysisExif, exifStatus, listExif } from 'api/web/exif';
import All from './All';
import Date from "./Date";
import Month from "./Month";
import Year from "./Year";
import Exif from "./Exif";
import { parseShotTime, parseShotEquipment } from "api/utils/api";

export default function ({ show }) {
    const [exifMap, setExifMap] = useState({});
    const [viewBy, setViewBy] = useState("所有照片");
    const [chosenShotEquipment, setChosenShotEquipment] = useState([]);
    const [shotEquipmentMap, setShotEquipmentMap] = useState([]);
    return (
        <Stack style={{ width: "100%", height: "100%", display: show ? undefined : "none" }}>
            <Exif onNewExif={() => {
                listExif().then(exifMap => {
                    setExifMap(exifMap);
                    let shotEquipmentMap = {};
                    Object.values(exifMap).forEach(exif => {
                        let shotEquipment = parseShotEquipment(exif);
                        let shotTime = parseShotTime(exif);
                        if (shotEquipmentMap.hasOwnProperty(shotEquipment)) {
                            shotEquipmentMap[shotEquipment]++;
                        } else {
                            shotEquipmentMap[shotEquipment] = 1;
                        }
                        exif.shotEquipment = shotEquipment;
                        exif.shotTime = shotTime;
                    })
                    setShotEquipmentMap(shotEquipmentMap);
                });
            }} />
            <Box>
                <InputLabel sx={{ display: "inline" }}>视图：</InputLabel>
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
                <InputLabel sx={{ display: "inline" }}>拍摄设备：</InputLabel>
                {Object.keys(shotEquipmentMap).map((shotEquipment, i) =>
                    <FormControlLabel key={i} control={
                        <Checkbox checked={chosenShotEquipment.includes(shotEquipment)} value={shotEquipment} onChange={e => {
                            setChosenShotEquipment(prev => {
                                let set = new Set(prev);
                                if (e.target.checked) {
                                    set.add(shotEquipment);
                                    return Array.from(set);
                                } else {
                                    set.delete(shotEquipment);
                                    return Array.from(set);
                                }
                            })
                        }} />
                    } label={(shotEquipment ? shotEquipment : "未知设备") + " (" + shotEquipmentMap[shotEquipment] + ")"} />
                )}
            </Box>
            {viewBy == "年" && <Year exifMap={exifMap} chosenShotEquipment={chosenShotEquipment} />}
            {viewBy == "月" && <Month exifMap={exifMap} chosenShotEquipment={chosenShotEquipment} />}
            {viewBy == "日" && <Date exifMap={exifMap} chosenShotEquipment={chosenShotEquipment} />}
            {viewBy == "所有照片" && <All exifMap={exifMap} chosenShotEquipment={chosenShotEquipment} />}
        </Stack>
    );
}
