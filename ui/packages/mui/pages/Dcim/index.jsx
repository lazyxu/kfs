import { Box, Checkbox, FormControlLabel, IconButton, InputLabel, Radio, RadioGroup, Stack, Typography } from "@mui/material";
// import FormControlLabel from '@mui/material/FormControlLabel';
import { parseShotEquipment, parseShotTime, timeSortFn } from "@kfs/common/api/utils";
import { listExif } from '@kfs/mui/api/exif';
import { CalendarMonth, FilterAlt, Refresh } from "@mui/icons-material";
import { useEffect, useState } from "react";
import All from './All';
import Date from "./Date";
import Month from "./Month";
import Year from "./Year";

export default function () {
    const [metadataList, setMetadataList] = useState([]);
    const [viewBy, setViewBy] = useState("所有照片");
    const [calendar, setCalendar] = useState(false);
    const [filter, setFilter] = useState(false);
    const [chosenShotEquipment, setChosenShotEquipment] = useState();
    const [shotEquipmentMap, setShotEquipmentMap] = useState({});
    const [chosenFileType, setChosenFileType] = useState();
    const [fileTypeMap, setFileTypeMap] = useState({});
    const refersh = () => {
        listExif().then(metadataList => {
            let shotEquipmentMap = {};
            let fileTypeMap = {};
            metadataList.forEach(metadata => {
                let { fileType } = metadata;
                let shotEquipment = parseShotEquipment(metadata);
                let shotTime = parseShotTime(metadata);
                if (shotEquipmentMap.hasOwnProperty(shotEquipment)) {
                    shotEquipmentMap[shotEquipment]++;
                } else {
                    shotEquipmentMap[shotEquipment] = 1;
                }
                if (fileTypeMap.hasOwnProperty(fileType.extension)) {
                    fileTypeMap[fileType.extension]++;
                } else {
                    fileTypeMap[fileType.extension] = 1;
                }
                metadata.shotEquipment = shotEquipment;
                metadata.shotTime = shotTime;
            })
            setMetadataList(metadataList);
            setShotEquipmentMap(shotEquipmentMap);
            setFileTypeMap(fileTypeMap);
        });
    }
    useEffect(() => {
        refersh();
    }, []);
    let filteredMetadataList = metadataList
        .filter(metadata =>
            (!chosenShotEquipment || chosenShotEquipment.includes(metadata.shotEquipment)) &&
            (!chosenFileType || chosenFileType.includes(metadata.fileType.extension)))
        .sort(timeSortFn);
    return (
        <Stack style={{ flex: "1", overflowY: 'auto' }}>
            <Stack
                direction="row"
                justifyContent="flex-end"
                alignItems="flex-end"
                spacing={0.5}
            >
                <IconButton onClick={() => refersh()} ><Refresh /></IconButton>
                <IconButton onClick={() => setCalendar(f => !f)}><CalendarMonth /></IconButton>
                <IconButton onClick={() => setFilter(f => !f)}><FilterAlt /></IconButton>
            </Stack>
            <Box sx={{ ...(!calendar && { display: 'none' }) }}>
                <InputLabel sx={{ display: "inline" }}>视图：</InputLabel>
                <RadioGroup sx={{ display: "inline" }}
                    row
                    value={viewBy}
                    onChange={e => setViewBy(e.target.value)}
                    size="small"
                >
                    {["年", "月", "日", "所有照片"].map(value =>
                        <FormControlLabel key={value} value={value} control={<Radio />} label={value} />
                    )}
                </RadioGroup>
            </Box>
            <Box sx={{ ...(!filter && { display: 'none' }) }}>
                <Box>
                    <InputLabel sx={{ display: "inline" }}>拍摄设备：</InputLabel>
                    {Object.keys(shotEquipmentMap).map((shotEquipment, i) =>
                        <FormControlLabel key={i} control={
                            <Checkbox checked={chosenShotEquipment?.includes(shotEquipment)} value={shotEquipment} onChange={e => {
                                setChosenShotEquipment(prev => {
                                    let set = new Set(prev);
                                    if (e.target.checked) {
                                        set.add(shotEquipment);
                                        return Array.from(set);
                                    } else {
                                        set.delete(shotEquipment);
                                        if (set.size === 0) {
                                            return undefined;
                                        }
                                        return Array.from(set);
                                    }
                                })
                            }} />
                        } label={(shotEquipment ? shotEquipment : "未知设备") + " (" + shotEquipmentMap[shotEquipment] + ")"} />
                    )}
                </Box>
                <Box>
                    <InputLabel sx={{ display: "inline" }}>文件类型：</InputLabel>
                    {Object.keys(fileTypeMap).map((fileType, i) =>
                        <FormControlLabel key={i} control={
                            <Checkbox checked={chosenFileType?.includes(fileType)} value={fileType} onChange={e => {
                                setChosenFileType(prev => {
                                    let set = new Set(prev);
                                    if (e.target.checked) {
                                        set.add(fileType);
                                        return Array.from(set);
                                    } else {
                                        set.delete(fileType);
                                        if (set.size === 0) {
                                            return undefined;
                                        }
                                        return Array.from(set);
                                    }
                                })
                            }} />
                        } label={fileType + " (" + fileTypeMap[fileType] + ")"} />
                    )}
                </Box>
            </Box>
            <Typography>共{metadataList.filter(m => m.fileType.type === "image").length}张照片、{metadataList.filter(m => m.fileType.type === "video").length}个视频</Typography>
            {(chosenShotEquipment || chosenFileType) && <Typography>筛选出{filteredMetadataList.filter(m => m.fileType.type === "image").length}张照片、{filteredMetadataList.filter(m => m.fileType.type === "video").length}个视频</Typography>}
            {viewBy === "年" && <Year metadataList={filteredMetadataList} />}
            {viewBy === "月" && <Month metadataList={filteredMetadataList} />}
            {viewBy === "日" && <Date metadataList={filteredMetadataList} />}
            {viewBy === "所有照片" && <All metadataList={filteredMetadataList} />}
        </Stack>
    );
}
