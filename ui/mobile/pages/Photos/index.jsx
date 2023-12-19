import { parseShotEquipment, parseShotTime, timeSortFn } from "@kfs/common/api/utils";
import { listExif } from '@kfs/common/api/webServer/exif';
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useEffect, useRef, useState } from "react";
import { Image, View } from 'react-native';
import { Appbar } from "react-native-paper";
function calImageWith(gridWith) {
    const n = gridWith / 100;
    return gridWith / Math.ceil(n);
}

export default function () {
    const [metadataList, setMetadataList] = useState([]);
    const [viewBy, setViewBy] = useState("所有照片");
    const [calendar, setCalendar] = useState(false);
    const [filter, setFilter] = useState(false);
    const [chosenShotEquipment, setChosenShotEquipment] = useState();
    const [shotEquipmentMap, setShotEquipmentMap] = useState({});
    const [chosenFileType, setChosenFileType] = useState();
    const [fileTypeMap, setFileTypeMap] = useState({});
    const ref = useRef(null);
    const [width, setWidth] = useState(0);
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
        console.log(ref);
        setWidth(calImageWith(ref.current.offsetWidth));
    }, []);
    console.log(width);
    let filteredMetadataList = metadataList
        .filter(metadata =>
            (!chosenShotEquipment || chosenShotEquipment.includes(metadata.shotEquipment)) &&
            (!chosenFileType || chosenFileType.includes(metadata.fileType.extension)))
        .sort(timeSortFn);
    return (
        <>
            <Appbar.Header>
                <Appbar.Content title="照片" />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <View style={{
                display: "flex",
                backgroundColor: '#fff',
                width: "100%",
                height: "100%",
                overflow: "scroll",
                flexDirection: 'row',
                flexWrap: "wrap",
            }} ref={ref}>
                {filteredMetadataList.map(metadata =>
                    <Image key={metadata.hash} style={{
                        width: width,
                        height: width,
                    }} source={{ uri: `${getSysConfig().webServer}/thumbnail?size=256&cutSquare=true&hash=${metadata.hash}` }} />
                )}
            </View>
        </>
    );
}
