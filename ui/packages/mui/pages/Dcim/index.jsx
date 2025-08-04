import { Box, Checkbox, FormControlLabel, IconButton, InputLabel, Radio, RadioGroup, Stack, Typography } from "@mui/material";
// import FormControlLabel from '@mui/material/FormControlLabel';
import { parseShotEquipment } from "@kfs/common/api/utils";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { noteError } from "@kfs/mui/components/Notification";
import { CalendarMonth, FilterAlt, Refresh } from "@mui/icons-material";
import { EventStreamContentType, fetchEventSource } from "@microsoft/fetch-event-source";
import moment from 'moment';
import { useEffect, useRef, useState } from "react";
import ThumbnailList from '@kfs/common/components/ThumbnailList/ThumbnailList';
import ThumbnailGroupedList from "@kfs/common/components/ThumbnailList/ThumbnailGroupedList";
import { listDCIM } from "@kfs/common/api/webServer/eventsource";

export default function () {
    const [metadataList, setMetadataList] = useState([]);
    const [metadataListTotal, setMetadataListTotal] = useState(0);
    const [viewBy, setViewBy] = useState("所有照片");
    const [calendar, setCalendar] = useState(false);
    const [filter, setFilter] = useState(false);
    const [chosenShotEquipment, setChosenShotEquipment] = useState();
    const [shotEquipmentMap, setShotEquipmentMap] = useState({});
    const [chosenFileType, setChosenFileType] = useState();
    const [fileTypeMap, setFileTypeMap] = useState({});
    const abortControllerRef = useRef(null);
    const refersh = () => {
        console.log("===refresh DCIM===");
        if (abortControllerRef.current) {
            abortControllerRef.current.abort();
        }
        abortControllerRef.current = new AbortController();
        listDCIM(abortControllerRef.current.signal, () => {
            setDcimList([]);
            setMetadataList([]);
            setMetadataListTotal(0);
            setShotEquipmentMap({});
            setFileTypeMap({});
        }, json => {
            if (json.dcimList) {
                setMetadataList(prev => {
                    const newList = [...prev, ...json.dcimList];
                    let shotEquipmentMap = {};
                    let fileTypeMap = {};
                    newList.forEach(metadata => {
                        let shotEquipment = parseShotEquipment(metadata);
                        if (shotEquipmentMap.hasOwnProperty(shotEquipment)) {
                            shotEquipmentMap[shotEquipment]++;
                        } else {
                            shotEquipmentMap[shotEquipment] = 1;
                        }
                        if (fileTypeMap.hasOwnProperty(metadata.extension)) {
                            fileTypeMap[metadata.extension]++;
                        } else {
                            fileTypeMap[metadata.extension] = 1;
                        }
                        metadata.shotEquipment = shotEquipment;
                        metadata.parsedTime = moment.parseZone(metadata.time, "YYYY-MM-DDTHH:mm:ssZZ");
                    });
                    setShotEquipmentMap(shotEquipmentMap);
                    setFileTypeMap(fileTypeMap);
                    return newList;
                });
            } else {
                setMetadataListTotal(json.n);
            }
        }, (tag, err) => {
            noteError(tag+": " + err);
        });
    }
    useEffect(() => {
        refersh();
        return () => {
            console.log("===abort DCIM===");
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        }
    }, []);
    let filteredMetadataList = metadataList
        .filter(metadata =>
            (!chosenShotEquipment || chosenShotEquipment.includes(metadata.shotEquipment)) &&
            (!chosenFileType || chosenFileType.includes(metadata.extension)))
        .sort((a, b) => a.parsedTime.isAfter(b.parsedTime));
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
                                            return undefined; // undefined means no filter
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
                                            return undefined; // undefined means no filter
                                        }
                                        return Array.from(set);
                                    }
                                })
                            }} />
                        } label={fileType + " (" + fileTypeMap[fileType] + ")"} />
                    )}
                </Box>
            </Box>
            <Typography>共{metadataList.filter(m => m.type=== "image").length}张照片、{metadataList.filter(m => m.type=== "video").length}个视频 ({metadataList.length}/{metadataListTotal})</Typography>
            {(chosenShotEquipment || chosenFileType) && <Typography>筛选出{filteredMetadataList.filter(m => m.type=== "image").length}张照片、{filteredMetadataList.filter(m => m.type=== "video").length}个视频</Typography>}
            {viewBy === "年" && <ThumbnailGroupedList metadataList={filteredMetadataList} spacing={0} minCol={10} formatFn={date => date.format("YYYY年")} />}
            {viewBy === "月" && <ThumbnailGroupedList metadataList={filteredMetadataList} spacing={0.5} minCol={4} formatFn={date => date.format("YYYY年MM月")} />}
            {viewBy === "日" && <ThumbnailGroupedList metadataList={filteredMetadataList} spacing={1} minCol={3} formatFn={date => date.format("YYYY年MM月DD日")} />}
            {viewBy === "所有照片" && <ThumbnailList metadataList={filteredMetadataList} />}
        </Stack>
    );
}
